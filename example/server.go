package example

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type Service interface {
	GetPost(ctxt context.Context, v0 string) (rtn Post, err error)
	GetPosts(ctxt context.Context) (rtn []Post, err error)
	GetUser(ctxt context.Context, v0 string) (rtn User, err error)
	GetUsers(ctxt context.Context) (rtn []User, err error)
	PatchPost(ctxt context.Context, v0 string, v1 PatchPostRequest) (rtn Post, err error)
	PatchUser(ctxt context.Context, v0 string, v1 PatchUserRequest) (rtn User, err error)
	PutPost(ctxt context.Context, v0 PutPostRequest) (rtn Post, err error)
	PutUser(ctxt context.Context, v0 PutUserRequest) (rtn User, err error)
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

	s.Router.Get("/api/v1/posts/{id}", s.handleGetPost)
	s.Router.Get("/api/v1/posts", s.handleGetPosts)
	s.Router.Get("/api/v1/users/{id}", s.handleGetUser)
	s.Router.Get("/api/v1/users", s.handleGetUsers)
	s.Router.Patch("/api/v1/posts/{id}", s.handlePatchPost)
	s.Router.Patch("/api/v1/users/{id}", s.handlePatchUser)
	s.Router.Put("/api/v1/posts", s.handlePutPost)
	s.Router.Put("/api/v1/users", s.handlePutUser)

	return s
}

func (c *Server) handleGetPost(w http.ResponseWriter, r *http.Request) {
	v0 := chi.URLParam(r, "id")

	r0, err := c.srv.GetPost(r.Context(), v0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (c *Server) handleGetPosts(w http.ResponseWriter, r *http.Request) {

	r0, err := c.srv.GetPosts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (c *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	v0 := chi.URLParam(r, "id")

	r0, err := c.srv.GetUser(r.Context(), v0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (c *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {

	r0, err := c.srv.GetUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (c *Server) handlePatchPost(w http.ResponseWriter, r *http.Request) {
	v0 := chi.URLParam(r, "id")

	var v1 PatchPostRequest
	if err := json.NewDecoder(r.Body).Decode(&v1); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r0, err := c.srv.PatchPost(r.Context(), v0, v1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (c *Server) handlePatchUser(w http.ResponseWriter, r *http.Request) {
	v0 := chi.URLParam(r, "id")

	var v1 PatchUserRequest
	if err := json.NewDecoder(r.Body).Decode(&v1); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r0, err := c.srv.PatchUser(r.Context(), v0, v1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (c *Server) handlePutPost(w http.ResponseWriter, r *http.Request) {

	var v0 PutPostRequest
	if err := json.NewDecoder(r.Body).Decode(&v0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r0, err := c.srv.PutPost(r.Context(), v0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (c *Server) handlePutUser(w http.ResponseWriter, r *http.Request) {

	var v0 PutUserRequest
	if err := json.NewDecoder(r.Body).Decode(&v0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r0, err := c.srv.PutUser(r.Context(), v0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}
