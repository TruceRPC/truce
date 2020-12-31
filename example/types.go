package types

import ()

type Resource struct {
	Id      string
	GroupId string
	Name    string
}

type PutResourceRequest struct {
	Name string
}
