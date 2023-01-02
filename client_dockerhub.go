// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type DockerHubClient interface {
	ListRepositories(ctx context.Context, repositoryName RepositoryName, ch chan<- DockerHubTagRepository) error
	ListTags(ctx context.Context, repositoryName RepositoryName, ch chan<- DockerHubTag) error
	DeleteTag(ctx context.Context, repositoryName RepositoryName, tag TagName) error
}

type dockerHubClient struct {
	httpClient HttpClient
	registry   Registry

	dockerhubMux   sync.Mutex
	dockerhubToken string

	dockerIoMux   sync.Mutex
	dockerIoToken string
}

func NewDockerHubClient(
	httpClient HttpClient,
	registry Registry,
) DockerHubClient {
	return &dockerHubClient{
		httpClient: httpClient,
		registry:   registry,
	}
}

type DockerHubTagRepository struct {
	User string `json:"user"`
	Name string `json:"name"`
}

func (d DockerHubTagRepository) RepositoryName() RepositoryName {
	return RepositoryName(d.User + "/" + d.Name)
}

func (c *dockerHubClient) ListRepositories(ctx context.Context, repositoryName RepositoryName, ch chan<- DockerHubTagRepository) error {
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/", repositoryName)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			glog.V(2).Infof("request url: %v", url)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				return errors.Wrap(err, "create http request failed")
			}
			var response struct {
				Next    string                   `json:"next"`
				Results []DockerHubTagRepository `json:"results"`
			}
			if err := c.doJSON(ctx, req, &response); err != nil {
				return errors.Wrap(err, "perform http request failed")
			}
			for _, result := range response.Results {
				ch <- result
			}
			if len(response.Next) == 0 {
				return nil
			}
			url = response.Next
		}
	}
}

type DockerHubTag struct {
	Tag         TagName `json:"name"`
	LastUpdated string  `json:"last_updated"`
}

func (c *dockerHubClient) ListTags(ctx context.Context, repositoryName RepositoryName, ch chan<- DockerHubTag) error {
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags/", repositoryName.String())
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			glog.V(2).Infof("request url: %v", url)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				return errors.Wrap(err, "create http request failed")
			}
			var response struct {
				Next    string         `json:"next"`
				Results []DockerHubTag `json:"results"`
			}
			if err := c.doJSON(ctx, req, &response); err != nil {
				return errors.Wrap(err, "perform http request failed")
			}
			for _, result := range response.Results {
				ch <- result
			}
			if len(response.Next) == 0 {
				return nil
			}
			url = response.Next
		}
	}
}
func (c *dockerHubClient) DeleteTag(ctx context.Context, repositoryName RepositoryName, tag TagName) error {
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags/%s/", repositoryName.String(), tag.String())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return errors.Wrap(err, "create http request failed")
	}
	if _, err := c.doSuccess(ctx, req); err != nil {
		return err
	}
	return nil
}

func (c *dockerHubClient) addAuth(ctx context.Context, req *http.Request) error {
	glog.V(2).Infof("auth with hub.docker.com")
	token, err := c.getDockerHubToken(ctx)
	if err != nil {
		return errors.Wrap(err, "get token failed")
	}
	req.Header.Add("Authorization", fmt.Sprintf("JWT %s", token))
	glog.V(2).Infof("set Authorization header")
	return nil
}

func (c *dockerHubClient) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	err := c.addAuth(ctx, req)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(ctx, req)
}

func (c *dockerHubClient) doSuccess(ctx context.Context, req *http.Request) (*http.Response, error) {
	err := c.addAuth(ctx, req)
	if err != nil {
		return nil, err
	}
	return c.httpClient.DoSuccess(ctx, req)
}

func (c *dockerHubClient) doJSON(ctx context.Context, req *http.Request, data interface{}) error {
	err := c.addAuth(ctx, req)
	if err != nil {
		return err
	}
	return c.httpClient.DoJSON(ctx, req, data)
}

func (c *dockerHubClient) getDockerHubToken(ctx context.Context) (string, error) {
	defer c.dockerhubMux.Unlock()
	c.dockerhubMux.Lock()
	if c.dockerhubToken == "" {
		b := bytes.NewBufferString(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, c.registry.Username, c.registry.Password))
		req, err := http.NewRequest("POST", "https://hub.docker.com/v2/users/login/", b)
		if err != nil {
			return "", errors.Wrap(err, "create request failed")
		}
		req.Header.Add("Content-Type", "application/json")
		var data struct {
			Token       string `json:"token"`
			AccessToken string `json:"access_token"`
		}
		if err := c.httpClient.DoJSON(ctx, req, &data); err != nil {
			return "", errors.Wrap(err, "request failed")
		}
		glog.V(2).Infof("got token: %s", data.Token)
		c.dockerhubToken = data.Token
	}
	return c.dockerhubToken, nil
}
