package types

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type Service interface {
	GetResource(ctxt context.Context, v0 string) (rtn Resource, err error)
	GetResources(ctxt context.Context) (rtn []Resource, err error)
	PutResource(ctxt context.Context, v0 string, v1 PutResourceRequest) (rtn Resource, err error)
}

type Server struct {
	chi.Router
	srv Service
}

func NewServer(srv Service) *Server {
	s := &Server{
		Router: chi.NewRouter(),
		srv:    srv,
	}

	s.Router.Get("/api/v0/resources/{id}", s.handleGetResource)
	s.Router.Get("/api/v0/resources", s.handleGetResources)
	s.Router.Put("/api/v0/group/{group_id}/resources", s.handlePutResource)

	return s
}

func (c *Server) handleGetResource(w http.ResponseWriter, r *http.Request) {
	v0 := chi.URLParam(r, "id")

	r0, err := c.srv.GetResource(r.Context(), v0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (c *Server) handleGetResources(w http.ResponseWriter, r *http.Request) {

	r0, err := c.srv.GetResources(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (c *Server) handlePutResource(w http.ResponseWriter, r *http.Request) {
	v0 := chi.URLParam(r, "group_id")

	var v1 PutResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&v1); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r0, err := c.srv.PutResource(r.Context(), v0, v1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}
