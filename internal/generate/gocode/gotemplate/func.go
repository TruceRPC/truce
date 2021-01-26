package gotemplate

import (
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/TruceRPC/truce"
)

var tmplFuncs = template.FuncMap{
	"args":      args,
	"backtick":  backtick,
	"errorFmt":  errorFmt,
	"method":    method,
	"name":      name,
	"pathJoin":  pathJoin,
	"signature": signature,
	"tags":      tags,
}

func name(f truce.Field) (v string) {
	parts := strings.Split(f.Name, "_")
	for _, p := range parts {
		v += strings.Title(p)
	}
	return
}

func tags(f truce.Field) string {
	return fmt.Sprintf("`json:%q`", f.Name)
}

func signature(f truce.Function) string {
	builder := &strings.Builder{}
	fmt.Fprintf(builder, "%s(ctxt context.Context, ", f.Name)
	for i, arg := range f.Arguments {
		if i > 0 {
			builder.WriteString(", ")
		}
		fmt.Fprintf(builder, "%s %s", arg.Name, arg.Type)
	}

	builder.WriteString(") (")
	if rtn := f.Return; rtn.Name != "" {
		fmt.Fprintf(builder, "%s, ", rtn.Type)
	}
	builder.WriteString("error)")

	return builder.String()
}

func pathJoin(path Path) string {
	if len(path) == 0 {
		return fmt.Sprintf("%q", path)
	}

	var hasVariables bool
	for _, elem := range path {
		if elem.Type == "variable" {
			hasVariables = true
		}
	}

	if hasVariables {
		return fmt.Sprintf(`fmt.Sprintf(%q, %s)`, path.FmtString(), path.ArgString())
	}
	return fmt.Sprintf(`%q`, path.FmtString())
}

func args(f truce.Function) (v string) {
	for i, arg := range f.Arguments {
		if i > 0 {
			v += ", "
		}
		v += arg.Name
	}
	return
}

func method(b *Function) string {
	return strings.Title(strings.ToLower(b.Method))
}

func errorFmt(t truce.Type) string {
	var (
		i              int
		fmtStr, argStr string
	)
	for _, field := range sortedFields(t.Fields) {
		if i > 0 {
			fmtStr += " "
			argStr += ", "
		}

		fmtStr += fmt.Sprintf("%s=%%q", field.Name)
		argStr += fmt.Sprintf("e.%s", name(field))
		i++
	}

	return fmt.Sprintf(`"error: %s", %s`, fmtStr, argStr)
}

func backtick(v string) string {
	return "`" + v + "`"
}

func sortedFields(m map[string]truce.Field) []truce.Field {
	var l []truce.Field
	for _, f := range m {
		l = append(l, f)
	}
	sort.Slice(l, func(i, j int) bool {
		return l[i].Name < l[j].Name
	})
	return l
}
