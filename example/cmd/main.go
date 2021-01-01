package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	types "github.com/georgemac/truce/example"
)

type service map[string]types.Resource

func (s *service) GetResource(ctxt context.Context, v0 string) (rtn types.Resource, err error) {
	var ok bool
	rtn, ok = (*s)[v0]
	if !ok {
		err = fmt.Errorf("not found")
	}

	return
}

func (s *service) GetResources(ctxt context.Context) (rtn []types.Resource, err error) {
	for _, r := range *s {
		rtn = append(rtn, r)
	}
	return
}

func (s *service) PutResource(ctxt context.Context, v0 string, v1 types.PutResourceRequest) (rtn types.Resource, err error) {
	resource := types.Resource{
		Id:      fmt.Sprintf("%x", rand.Int63()),
		Name:    v1.Name,
		GroupId: v0,
	}
	(*s)[resource.Id] = resource
	return resource, nil
}

func main() {
	srv := service{}
	server := types.NewServer(&srv)
	if err := http.ListenAndServe(":8080", server); err != nil {
		panic(err)
	}
}
