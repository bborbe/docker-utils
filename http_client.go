// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bborbe/io/reader_shadow_copy"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type HttpClient interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
	DoSuccess(ctx context.Context, req *http.Request) (*http.Response, error)
	DoJSON(ctx context.Context, req *http.Request, data interface{}) error
}

func NewHttpClient(client *http.Client) HttpClient {
	return &httpClient{
		client: client,
	}
}

type httpClient struct {
	client *http.Client
}

func (h *httpClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := h.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrapf(err, "%s request to %s failed", req.Method, req.URL.String())
	}
	glog.V(2).Infof("%s request to %s completed with status %d", req.Method, req.URL.String(), resp.StatusCode)
	return resp, err
}

func (h *httpClient) DoSuccess(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := h.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		if glog.V(4) {
			defer resp.Body.Close()
			bytes, _ := ioutil.ReadAll(resp.Body)
			glog.Infof(string(bytes))
		}
		return nil, errors.Errorf("%s request to %s failed with statusCode %d", req.Method, req.URL.String(), resp.StatusCode)
	}
	return resp, err
}

func (h *httpClient) DoJSON(ctx context.Context, req *http.Request, data interface{}) error {
	resp, err := h.DoSuccess(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	reader := reader_shadow_copy.New(resp.Body)
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		if glog.V(4) {
			glog.Infof(string(reader.Bytes()))
		}
		return errors.Wrap(err, "decode http response to json failed")
	}
	return nil
}
