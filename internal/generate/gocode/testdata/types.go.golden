package example

import (
	"fmt"
	"time"
)

var _ = time.After

type NotAuthorized struct {
	Message string `json:"message"`
}

func (e NotAuthorized) Error() string {
	return fmt.Sprintf("error: message=%q", e.Message)
}

type NotFound struct {
	Message string `json:"message"`
}

func (e NotFound) Error() string {
	return fmt.Sprintf("error: message=%q", e.Message)
}

type PatchPostRequest struct {
	Body    []byte     `json:"body"`
	Created *time.Time `json:"created"`
	Draft   bool       `json:"draft"`
	Title   string     `json:"title"`
}

type PatchUserRequest struct {
	Age    int                    `json:"age"`
	Height float64                `json:"height"`
	Labels map[string]interface{} `json:"labels"`
	Name   string                 `json:"name"`
}

type Post struct {
	Body    []byte     `json:"body"`
	Created *time.Time `json:"created"`
	Draft   bool       `json:"draft"`
	Id      string     `json:"id"`
	Title   string     `json:"title"`
}

type PutPostRequest struct {
	Body    []byte     `json:"body"`
	Created *time.Time `json:"created"`
	Draft   bool       `json:"draft"`
	Title   string     `json:"title"`
}

type PutUserRequest struct {
	Age    int                    `json:"age"`
	Height float64                `json:"height"`
	Labels map[string]interface{} `json:"labels"`
	Name   string                 `json:"name"`
}

type User struct {
	Age    int                    `json:"age"`
	Height float64                `json:"height"`
	Id     string                 `json:"id"`
	Labels map[string]interface{} `json:"labels"`
	Name   string                 `json:"name"`
}
