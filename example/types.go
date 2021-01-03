package types

import ()

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PutUserRequest struct {
	Name string `json:"name"`
}

type PatchUserRequest struct {
	Name string `json:"name"`
}

type Post struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PutPostRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PatchPostRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
