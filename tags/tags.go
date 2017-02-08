package tags

import (
	"encoding/json"
	"fmt"
	"github.com/bborbe/docker_utils/model"
	"github.com/bborbe/io/reader_shadow_copy"
	"github.com/golang/glog"
	"net/http"
	"sort"
)

type Tags interface {
	List(registry model.Registry, repositoryName model.RepositoryName) ([]model.Tag, error)
}

type repositories struct {
	httpClient *http.Client
}

func New(httpClient *http.Client) *repositories {
	c := new(repositories)
	c.httpClient = httpClient
	return c
}

func (r *repositories) List(registry model.Registry, repositoryName model.RepositoryName) ([]model.Tag, error) {
	url := fmt.Sprintf("%s/v2/%v/tags/list", registry.Name.Url(), repositoryName.String())
	glog.V(2).Infof("request url: %v", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		glog.V(0).Infof("create http request failed: %v", err)
		return nil, err
	}
	req.SetBasicAuth(registry.Username.String(), registry.Password.String())
	resp, err := r.httpClient.Do(req)
	if err != nil {
		glog.V(0).Infof("perform http request failed: %v", err)
		return nil, err
	}
	glog.V(4).Infof("response %d", resp.StatusCode)
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}
	var response response
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

type response struct {
	Tags []model.Tag `json:"tags"`
}
