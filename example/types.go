package types

import ()

type Resource struct {
	Id      string `json:"id"`
	GroupId string `json:"group_id"`
	Name    string `json:"name"`
}

type PutResourceRequest struct {
	Name string `json:"name"`
}
