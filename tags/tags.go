package tags

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/bborbe/docker-utils/model"
	"github.com/bborbe/io/reader_shadow_copy"
	"github.com/golang/glog"
)

type Tags interface {
	Delete(registry model.Registry, repositoryName model.RepositoryName, tag model.Tag) error
	Exists(registry model.Registry, repositoryName model.RepositoryName, tag model.Tag) (bool, error)
	List(registry model.Registry, repositoryName model.RepositoryName) ([]model.Tag, error)
}

type tagsConnector struct {
	httpClient *http.Client
}

func New(httpClient *http.Client) *tagsConnector {
	c := new(tagsConnector)
	c.httpClient = httpClient
	return c
}

func (r *tagsConnector) Delete(registry model.Registry, repositoryName model.RepositoryName, tag model.Tag) error {
	dockerContentDigest, err := r.dockerContentDigest(registry, repositoryName, tag)
	if err != nil {
		return fmt.Errorf("get content digest failed: %v", err)
	}
	url := fmt.Sprintf("%s/v2/%v/manifests/%v", registry.Name.Url(), repositoryName.String(), dockerContentDigest)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("build request failed: %v", err)
	}
	if err := registry.SetAuth(req); err != nil {
		glog.V(0).Info("set auth failed: %v", err)
		return err
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform http request failed: %v", err)
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("http status code %v != 2xx", resp.StatusCode)
	}
	glog.V(2).Infof("tag deleted")
	return nil
}

func (r *tagsConnector) dockerContentDigest(registry model.Registry, repositoryName model.RepositoryName, tag model.Tag) (string, error) {
	url := fmt.Sprintf("%s/v2/%v/manifests/%v", registry.Name.Url(), repositoryName.String(), tag.String())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("build request failed: %v", err)
	}
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	if err := registry.SetAuth(req); err != nil {
		glog.V(0).Info("set auth failed: %v", err)
		return "", err
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("perform http request failed: %v", err)
	}
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("http status code %v != 2xx", resp.StatusCode)
	}
	return resp.Header.Get("Docker-Content-Digest"), nil

}

func (r *tagsConnector) Exists(registry model.Registry, repositoryName model.RepositoryName, tag model.Tag) (bool, error) {
	tags, err := r.List(registry, repositoryName)
	if err != nil {
		glog.V(2).Infof("list tags failed: %v", err)
		return false, err
	}
	for _, t := range tags {
		if t == tag {
			glog.V(2).Infof("found tag")
			return true, nil
		}
	}
	glog.V(2).Infof("tag not found")
	return false, nil
}

func (r *tagsConnector) List(registry model.Registry, repositoryName model.RepositoryName) ([]model.Tag, error) {
	url := fmt.Sprintf("%s/v2/%v/tags/list", registry.Name.Url(), repositoryName.String())
	glog.V(2).Infof("request url: %v", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		glog.V(0).Infof("create http request failed: %v", err)
		return nil, err
	}
	if err := registry.SetAuth(req); err != nil {
		glog.V(0).Infof("set auth failed: %v", err)
		return nil, err
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		glog.V(0).Infof("perform http request failed: %v", err)
		return nil, err
	}
	glog.V(4).Infof("response %d", resp.StatusCode)
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}
	var response struct {
		Name string      `json:"name"`
		Tags []model.Tag `json:"tags"`
	}
	reader := reader_shadow_copy.New(resp.Body)
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		glog.V(0).Infof("decode http response to json failed: %v", err)
		return nil, err
	}
	if glog.V(4) {
		glog.Infof(string(reader.Bytes()))
	}
	tags := response.Tags
	sort.Sort(model.TagsByName(tags))
	return tags, nil
}
