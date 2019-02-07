package tags

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/pkg/errors"

	"github.com/bborbe/docker-utils/model"
	"github.com/bborbe/io/reader_shadow_copy"
	"github.com/golang/glog"
)

type Tags interface {
	Delete(registry model.Registry, repositoryName model.RepositoryName, tag model.Tag) error
	Exists(registry model.Registry, repositoryName model.RepositoryName, tag model.Tag) (bool, error)
	List(registry model.Registry, repositoryName model.RepositoryName) ([]model.Tag, error)
	Sha(registry model.Registry, repositoryName model.RepositoryName, tag model.Tag) (string, error)
	Manifest(registry model.Registry, repositoryName model.RepositoryName, tag model.Tag) (*Manifest, error)
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
	dockerContentDigest, err := r.Sha(registry, repositoryName, tag)
	if err != nil {
		return errors.Wrap(err, "get content digest failed")
	}
	url := fmt.Sprintf("%s/v2/%v/manifests/%v", registry.Name.Url(), repositoryName.String(), dockerContentDigest)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return errors.Wrap(err, "build request failed")
	}
	if err := registry.SetAuth(req); err != nil {
		return errors.Wrap(err, "set auth failed")
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "perform http request failed")
	}
	if resp.StatusCode/100 != 2 {
		return errors.Errorf("http status code %v != 2xx", resp.StatusCode)
	}
	glog.V(2).Infof("tag deleted")
	return nil
}

func (r *tagsConnector) Sha(registry model.Registry, repositoryName model.RepositoryName, tag model.Tag) (string, error) {
	url := fmt.Sprintf("%s/v2/%v/manifests/%v", registry.Name.Url(), repositoryName.String(), tag.String())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", errors.Wrap(err, "build request failed")
	}
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	if err := registry.SetAuth(req); err != nil {
		return "", errors.Wrap(err, "set auth failed")
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "perform http request failed")
	}
	if resp.StatusCode/100 != 2 {
		return "", errors.Errorf("http status code %v != 2xx", resp.StatusCode)
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%v", string(bytes))

	return resp.Header.Get("Docker-Content-Digest"), nil
}

func (r *tagsConnector) Exists(registry model.Registry, repositoryName model.RepositoryName, tag model.Tag) (bool, error) {
	tags, err := r.List(registry, repositoryName)
	if err != nil {
		return false, errors.Wrap(err, "list tags failed")
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
		return nil, errors.Wrap(err, "create http request failed")
	}
	if err := registry.SetAuth(req); err != nil {
		return nil, errors.Wrap(err, "set auth failed")
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "perform http request failed")
	}
	glog.V(4).Infof("response %d", resp.StatusCode)
	if resp.StatusCode/100 != 2 {
		return nil, errors.Errorf("request failed with status: %d", resp.StatusCode)
	}
	var response struct {
		Name string      `json:"name"`
		Tags []model.Tag `json:"tags"`
	}
	reader := reader_shadow_copy.New(resp.Body)
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode http response to json failed")
	}
	if glog.V(4) {
		glog.Infof(string(reader.Bytes()))
	}
	tags := response.Tags
	sort.Sort(model.TagsByName(tags))
	return tags, nil
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

func (r *tagsConnector) Manifest(registry model.Registry, repositoryName model.RepositoryName, tag model.Tag) (*Manifest, error) {
	url := fmt.Sprintf("%s/v2/%v/manifests/%v", registry.Name.Url(), repositoryName.String(), tag.String())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "build request failed")
	}
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	if err := registry.SetAuth(req); err != nil {
		return nil, errors.Wrap(err, "set auth failed")
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "perform http request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return nil, errors.Errorf("http status code %v != 2xx", resp.StatusCode)
	}
	var manifest Manifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, errors.Wrap(err, "decode json failed")
	}
	glog.V(2).Infof("manifest %v", manifest)
	return &manifest, nil
}
