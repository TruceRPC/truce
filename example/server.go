// DO NOT EDIT.
// This code was generated by Truce.
package example

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type Service interface {
	GetPost(ctxt context.Context, id string) (Post, error)
	GetPosts(ctxt context.Context) ([]Post, error)
	GetUser(ctxt context.Context, id string) (User, error)
	GetUsers(ctxt context.Context) ([]User, error)
	PatchPost(ctxt context.Context, id string, post PatchPostRequest) (Post, error)
	PatchUser(ctxt context.Context, id string, user PatchUserRequest) (User, error)
	PutPost(ctxt context.Context, post PutPostRequest) (Post, error)
	PutUser(ctxt context.Context, user PutUserRequest) (User, error)
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
	id := chi.URLParam(r, "id")
	r0, err := c.srv.GetPost(r.Context(), id)
	if err != nil {
		handleError(w, err)
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
		handleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (c *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	r0, err := c.srv.GetUser(r.Context(), id)
	if err != nil {
		handleError(w, err)
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
		handleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (c *Server) handlePatchPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var post PatchPostRequest
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r0, err := c.srv.PatchPost(r.Context(), id, post)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (c *Server) handlePatchUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var user PatchUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r0, err := c.srv.PatchUser(r.Context(), id, user)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (c *Server) handlePutPost(w http.ResponseWriter, r *http.Request) {
	var post PutPostRequest
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r0, err := c.srv.PutPost(r.Context(), post)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (c *Server) handlePutUser(w http.ResponseWriter, r *http.Request) {
	var user PutUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r0, err := c.srv.PutUser(r.Context(), user)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(r0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func handleError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case NotAuthorized:
		w.WriteHeader(401)
		if merr := json.NewEncoder(w).Encode(err); merr != nil {
			http.Error(w, merr.Error(), http.StatusInternalServerError)
		}

		return
	case NotFound:
		w.WriteHeader(404)
		if merr := json.NewEncoder(w).Encode(err); merr != nil {
			http.Error(w, merr.Error(), http.StatusInternalServerError)
		}

		return

	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
