package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	types "github.com/georgemac/truce/example"
)

var errNotFound = errors.New("not found")

type service struct {
	posts map[string]types.Post
	users map[string]types.User
}

func (s *service) GetPost(ctxt context.Context, v0 string) (rtn types.Post, err error) {
	var ok bool
	rtn, ok = s.posts[v0]
	if !ok {
		err = fmt.Errorf("%s %q: %w", "post", v0, errNotFound)
	}

	return
}

func (s *service) GetPosts(ctxt context.Context) (rtn []types.Post, err error) {
	for _, post := range s.posts {
		rtn = append(rtn, post)
	}
	return
}

func (s *service) GetUser(ctxt context.Context, v0 string) (rtn types.User, err error) {
	var ok bool
	rtn, ok = s.users[v0]
	if !ok {
		err = fmt.Errorf("%s %q: %w", "user", v0, errNotFound)
	}

	return
}

func (s *service) GetUsers(ctxt context.Context) (rtn []types.User, err error) {
	for _, user := range s.users {
		rtn = append(rtn, user)
	}
	return
}

func (s *service) PatchPost(ctxt context.Context, v0 string, v1 types.PatchPostRequest) (rtn types.Post, err error) {
	rtn, err = s.GetPost(ctxt, v0)
	if err != nil {
		return
	}

	rtn.Title = v1.Title
	rtn.Body = v1.Body
	s.posts[rtn.Id] = rtn

	return
}

func (s *service) PatchUser(ctxt context.Context, v0 string, v1 types.PatchUserRequest) (rtn types.User, err error) {
	rtn, err = s.GetUser(ctxt, v0)
	if err != nil {
		return
	}

	rtn.Name = v1.Name
	s.users[rtn.Id] = rtn

	return
}

func (s *service) PutPost(ctxt context.Context, v0 types.PutPostRequest) (rtn types.Post, err error) {
	rtn = types.Post{
		Id:    fmt.Sprintf("%x", rand.Int63()),
		Title: v0.Title,
		Body:  v0.Body,
	}

	s.posts[rtn.Id] = rtn

	return
}

func (s *service) PutUser(ctxt context.Context, v0 types.PutUserRequest) (rtn types.User, err error) {
	rtn = types.User{
		Id:   fmt.Sprintf("%x", rand.Int63()),
		Name: v0.Name,
	}

	s.users[rtn.Id] = rtn

	return
}

func main() {
	srv := service{
		posts: map[string]types.Post{},
		users: map[string]types.User{},
	}
	server := types.NewServer(&srv)
	if err := http.ListenAndServe(":8080", server); err != nil {
		panic(err)
	}
}
