package example

import (
	"encoding/json"
	"errors"
	"fmt"
)

type NotAuthorized struct {
	Message string `json:"message"`
}

func (e NotAuthorized) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Truce_ErrorType string `json:"error_type"`
		Message         string `json:"message"`
	}{
		Truce_ErrorType: "NotAuthorized",
		Message:         e.Message,
	})
}

func (e NotAuthorized) Error() string {
	return fmt.Sprintf("error: message=%q", e.Message)
}

type NotFound struct {
	Message string `json:"message"`
}

func (e NotFound) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Truce_ErrorType string `json:"error_type"`
		Message         string `json:"message"`
	}{
		Truce_ErrorType: "NotFound",
		Message:         e.Message,
	})
}

func (e NotFound) Error() string {
	return fmt.Sprintf("error: message=%q", e.Message)
}

type PatchPostRequest struct {
	Body  string `json:"body"`
	Title string `json:"title"`
}

type PatchUserRequest struct {
	Name string `json:"name"`
}

type Post struct {
	Body  string `json:"body"`
	Id    string `json:"id"`
	Title string `json:"title"`
}

type PutPostRequest struct {
	Body  string `json:"body"`
	Title string `json:"title"`
}

type PutUserRequest struct {
	Name string `json:"name"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func UnmarshalJSONError(data []byte, dst *error) error {
	v := struct {
		Truce_ErrorType string `json:"error_type"`
	}{}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Truce_ErrorType {
	case "NotAuthorized":

		v := NotAuthorized{}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}

		*dst = v
	case "NotFound":

		v := NotFound{}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}

		*dst = v

	default:
		return errors.New("internal service error")
	}

	return nil
}
