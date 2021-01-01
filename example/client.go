package types

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	client *http.Client
	host   *url.URL
}

func (c *Client) GetResource(ctxt context.Context, v0 string) (rtn Resource, err error) {
	u, err := c.host.Parse(fmt.Sprintf("/api/v0/resources/%v", v0))
	if err != nil {
		return
	}

	var (
		body io.Reader
		resp *http.Response
	)

	req, err := http.NewRequest("GET", u.String(), body)
	if err != nil {
		return
	}

	resp, err = c.client.Do(req.WithContext(ctxt))
	if err != nil {
		return
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	err = json.NewDecoder(resp.Body).Decode(&rtn)

	return
}

func (c *Client) GetResources(ctxt context.Context) (rtn []Resource, err error) {
	u, err := c.host.Parse(fmt.Sprintf("/api/v0/resources"))
	if err != nil {
		return
	}

	var (
		body io.Reader
		resp *http.Response
	)

	req, err := http.NewRequest("GET", u.String(), body)
	if err != nil {
		return
	}

	resp, err = c.client.Do(req.WithContext(ctxt))
	if err != nil {
		return
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	err = json.NewDecoder(resp.Body).Decode(&rtn)

	return
}

func (c *Client) PutResource(ctxt context.Context, v0 string, v1 PutResourceRequest) (rtn Resource, err error) {
	u, err := c.host.Parse(fmt.Sprintf("/api/v0/group/%v/resources", v0))
	if err != nil {
		return
	}

	var (
		body io.Reader
		resp *http.Response
	)

	buf := &bytes.Buffer{}
	body = buf
	if err = json.NewEncoder(buf).Encode(v1); err != nil {
		return
	}

	req, err := http.NewRequest("PUT", u.String(), body)
	if err != nil {
		return
	}

	resp, err = c.client.Do(req.WithContext(ctxt))
	if err != nil {
		return
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	err = json.NewDecoder(resp.Body).Decode(&rtn)

	return
}
