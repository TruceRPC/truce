package generate

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"strings"
	"text/template"

	"github.com/georgemac/truce"
	"github.com/georgemac/truce/pkg/http"
)

func GenerateTypes(w io.Writer, api truce.API) error {
	context := struct {
		Package string
		Imports string
		Types   []truce.Type
	}{
		Package: "types",
		Imports: "",
		Types:   api.Types,
	}

	return writeGo(w, typeTmpl, context)
}

func GenerateClient(w io.Writer, api truce.API) error {
	context := struct {
		Package    string
		Imports    string
		ClientName string
		Bindings   http.Bindings
	}{
		Package:    "types",
		Imports:    "",
		ClientName: "Client",
		Bindings:   http.BindingsFrom(api),
	}

	return writeGo(w, clientTmpl, context)
}

func GenerateServer(w io.Writer, api truce.API) error {
	context := struct {
		Package    string
		Imports    string
		ServerName string
		Bindings   http.Bindings
	}{
		Package:    "types",
		Imports:    "",
		ServerName: "Server",
		Bindings:   http.BindingsFrom(api),
	}

	return writeGo(w, serverTmpl, context)
}

func writeGo(w io.Writer, tmpl *template.Template, context interface{}) error {
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, context); err != nil {
		return err
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write(src)
	return err
}

var tmplFuncs = template.FuncMap{
	"name": func(f truce.Field) (v string) {
		parts := strings.Split(f.Name, "_")
		for _, p := range parts {
			v += strings.Title(p)
		}
		return
	},
	"signature": func(f truce.Function) string {
		builder := &strings.Builder{}
		fmt.Fprintf(builder, "%s(ctxt context.Context, ", f.Name)
		for i, arg := range f.Arguments {
			if i > 0 {
				builder.Write([]byte(", "))
			}

			fmt.Fprintf(builder, "v%d %s", i, arg.Type)
		}

		builder.Write([]byte(") ("))
		if rtn := f.Return; rtn.Name != "" {
			fmt.Fprintf(builder, "rtn %s, ", rtn.Type)
		}

		builder.Write([]byte("err error)"))

		return builder.String()
	},
	"path": func(b http.Binding) string {
		if len(b.PathArgs) == 0 {
			return `"` + b.PathFmt + `"`
		}

		return fmt.Sprintf(`fmt.Sprintf(%q, %s)`, b.PathFmt, b.PathArgs)
	},
	"args": func(f truce.Function) (v string) {
		for i := range f.Arguments {
			if i > 0 {
				v += ", "
			}

			v += fmt.Sprintf("v%d", i)
		}
		return
	},
	"method": func(b http.Binding) string {
		return strings.Title(strings.ToLower(b.Method))
	},
}

var typeTmpl = template.Must(
	template.
		New("types").
		Funcs(tmplFuncs).
		Parse(`package {{ .Package }}

import (
{{ .Imports }}
)

{{ range $type := .Types }}
type {{ $type.Name }} struct {
{{ range $field := $type.Fields }}  {{name .}}   {{.Type}}
{{ end }}
}
{{ end }}`),
)

var clientTmpl = template.Must(
	template.
		New("client").
		Funcs(tmplFuncs).
		Parse(`package {{ .Package }}

import (
"fmt"
"io"
"io/ioutil"
"net/http"
"net/url"
"encoding/json"
"bytes"
"context"
{{ .Imports }}
)

type {{.ClientName}} struct {
    client *http.Client
    host *url.URL
}

{{ $ctxt := . }}{{ range .Bindings }}
func (c *{{ $ctxt.ClientName }}) {{signature .Function}} {
    u, err := c.host.Parse({{path .}})
    if err != nil {
        return
    }

    var (
        body io.Reader
        resp *http.Response
    )

    {{if ne .BodyArg ""}}buf := &bytes.Buffer{}
    body = buf
    if err = json.NewEncoder(buf).Encode({{.BodyArg}}); err != nil {
        return
    }{{end}}

    req, err := http.NewRequest("{{.Method}}", u.String(), body)
    if err != nil {
        return
    }

    resp, err = c.client.Do(req.WithContext(ctxt))
    if err != nil {
        return
    }

    defer func(){
        _, _ = io.Copy(ioutil.Discard, resp.Body)
        _ = resp.Body.Close()
    }()

    {{if .HasReturn}}err = json.NewDecoder(resp.Body).Decode(&rtn){{end}}

    return
}
{{ end }}`),
)

var serverTmpl = template.Must(
	template.
		New("server").
		Funcs(tmplFuncs).
		Parse(`package {{ .Package }}

import (
"context"
"net/http"

"github.com/go-chi/chi"
{{ .Imports }}
)

type Service interface {
    {{range .Bindings}}{{signature .Function}}
    {{end}}
}

type {{.ServerName}} struct {
    chi.Router
}
{{ $ctxt := . }}
func NewServer() *{{.ServerName}} {
    s := &{{.ServerName}}{
      Router: chi.NewRouter(),
    }

    {{ range .Bindings }}s.Router.{{method .}}("{{.Path}}", s.handle{{.Function.Name}})
    {{end}}

    return s
}

{{ range .Bindings }}
func (c *{{ $ctxt.ServerName }}) handle{{.Function.Name}}(w http.ResponseWriter, r *http.Request) {
    p0 := chi.URLParam(r, "id")
	r0, err := c.srv.{{.Function.Name}}(r.Context(), {{args .Function}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}
{{ end }}`),
)
