// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"runtime"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type V2Client interface {
	ListRepositories(ctx context.Context, ch chan<- RepositoryName) error
	DeleteTag(ctx context.Context, repositoryName RepositoryName, tag TagName) error
	ExistsTag(ctx context.Context, repositoryName RepositoryName, tag TagName) (bool, error)
	ListTags(ctx context.Context, repositoryName RepositoryName, ch chan<- TagName) error
	Sha(ctx context.Context, repositoryName RepositoryName, tag TagName) (string, error)
	Manifest(ctx context.Context, repositoryName RepositoryName, tag TagName) (*Manifest, error)
}

type v2Client struct {
	httpClient HttpClient
	registry   Registry
}

func NewV2Client(
	httpClient HttpClient,
	registry Registry,
) V2Client {
	return &v2Client{
		httpClient: httpClient,
		registry:   registry,
	}
}

func (c *v2Client) ListRepositories(ctx context.Context, ch chan<- RepositoryName) error {
	url := fmt.Sprintf("%s/v2/_catalog", c.registry.Url)
	glog.V(2).Infof("request url: %v", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "create http request failed")
	}
	var response struct {
		Repositories []RepositoryName `json:"repositories"`
	}
	if err := c.doJSON(ctx, req, &response); err != nil {
		return errors.Wrap(err, "perform http request failed")
	}
	for _, repositoryName := range response.Repositories {
		ch <- repositoryName
	}
	return nil
}

func (c *v2Client) DeleteTag(ctx context.Context, repositoryName RepositoryName, tag TagName) error {
	dockerContentDigest, err := c.Sha(ctx, repositoryName, tag)
	if err != nil {
		return errors.Wrap(err, "get content digest failed")
	}
	url := fmt.Sprintf("%s/v2/%v/manifests/%v", c.registry.Url, repositoryName.String(), dockerContentDigest)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return errors.Wrap(err, "build request failed")
	}
	_, err = c.doSuccess(ctx, req)
	if err != nil {
		return errors.Wrap(err, "perform http request failed")
	}
	glog.V(2).Infof("tag deleted")
	return nil
}

func (c *v2Client) Sha(ctx context.Context, repositoryName RepositoryName, tag TagName) (string, error) {
	url := fmt.Sprintf("%s/v2/%v/manifests/%v", c.registry.Url, repositoryName.String(), tag.String())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", errors.Wrap(err, "build request failed")
	}
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	resp, err := c.doSuccess(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "perform http request failed")
	}
	return resp.Header.Get("Docker-Content-Digest"), nil
}

func (c *v2Client) ExistsTag(ctx context.Context, repositoryName RepositoryName, tag TagName) (bool, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	tags := make(chan TagName, runtime.NumCPU())
	go func() {
		defer close(tags)
		if err := c.ListTags(ctx, repositoryName, tags); err != nil {
			glog.Warningf("list tags failed :%v", err)
		}
	}()
	for t := range tags {
		if t == tag {
			glog.V(2).Infof("found tag")
			return true, nil
		}
	}
	glog.V(2).Infof("tag not found")
	return false, nil
}

func (c *v2Client) ListTags(ctx context.Context, repositoryName RepositoryName, ch chan<- TagName) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/%s/tags/list", c.registry.Url, repositoryName.String()), nil)
	if err != nil {
		return errors.Wrap(err, "create http request failed")
	}
	var response struct {
		Tags []TagName `json:"tags"`
	}
	if err := c.doJSON(ctx, req, &response); err != nil {
		return errors.Wrap(err, "perform http request failed")
	}
	for _, result := range response.Tags {
		ch <- result
	}
	return nil
}

type ManifestConfig struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}

type Manifest struct {
	SchemaVersion int              `json:"schemaVersion"`
	MediaType     string           `json:"mediaType"`
	Config        ManifestConfig   `json:"config"`
	Layers        []ManifestConfig `json:"layers"`
}

func (c *v2Client) Manifest(ctx context.Context, repositoryName RepositoryName, tag TagName) (*Manifest, error) {
	url := fmt.Sprintf("%s/v2/%v/manifests/%v", c.registry.Url, repositoryName.String(), tag.String())
	method := http.MethodGet
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "build request failed")
	}
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	var manifest Manifest
	if err := c.doJSON(ctx, req, &manifest); err != nil {
		return nil, errors.Wrap(err, "perform http request failed")
	}
	glog.V(2).Infof("manifest %v", manifest)
	return &manifest, nil
}

func (c *v2Client) addAuth(ctx context.Context, req *http.Request) error {
	if req.URL.Host == "registry-1.docker.io" {
		glog.V(2).Infof("auth with registry.docker.io")
		token, err := c.getDockerIoToken(ctx, req)
		if err != nil {
			return errors.Wrap(err, "get token failed")
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		glog.V(2).Infof("set Authorization header")
		return nil
	}
	if c.registry.Username != "" && c.registry.Password != "" {
		glog.V(2).Infof("basic auth")
		req.SetBasicAuth(c.registry.Username, c.registry.Password)
		glog.V(2).Infof("set basic auth")
		return nil
	}
	return nil
}

func (c *v2Client) getDockerIoToken(ctx context.Context, req *http.Request) (string, error) {
	re, err := regexp.Compile(`(?is)^/v2/(.*?/.*?)/`)
	if err != nil {
		return "", errors.Wrap(err, "invalid regex")
	}
	matches := re.FindStringSubmatch(req.URL.Path)
	if len(matches) < 2 {
		return "", errors.New("regex does not match")
	}
	var data struct {
		Token string `json:"token"`
	}
	repository := matches[1]
	glog.V(2).Infof("get token for repository: %s", repository)
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull,push,delete", repository), nil)
	if err != nil {
		return "", errors.Wrap(err, "create request failed")
	}
	if err := c.httpClient.DoJSON(ctx, req, &data); err != nil {
		return "", err
	}
	return data.Token, nil
}

func (c *v2Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	err := c.addAuth(ctx, req)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(ctx, req)
}

func (c *v2Client) doSuccess(ctx context.Context, req *http.Request) (*http.Response, error) {
	err := c.addAuth(ctx, req)
	if err != nil {
		return nil, err
	}
	return c.httpClient.DoSuccess(ctx, req)
}

func (c *v2Client) doJSON(ctx context.Context, req *http.Request, data interface{}) error {
	err := c.addAuth(ctx, req)
	if err != nil {
		return err
	}
	return c.httpClient.DoJSON(ctx, req, data)
}
