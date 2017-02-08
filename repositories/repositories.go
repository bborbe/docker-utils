package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/bborbe/docker_utils/model"
	"github.com/bborbe/io/reader_shadow_copy"
	"github.com/golang/glog"
	"net/http"
)

type Repositories interface {
	List(registry model.Registry) ([]model.RepositoryName, error)
}

type repositories struct {
	httpClient *http.Client
}

func New(httpClient *http.Client) *repositories {
	c := new(repositories)
	c.httpClient = httpClient
	return c
}

func (r *repositories) List(registry model.Registry) ([]model.RepositoryName, error) {
	url := fmt.Sprintf("%s/v2/_catalog", registry.Name.Url())
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
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}
	var response response
	reader := reader_shadow_copy.New(resp.Body)
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		glog.V(0).Infof("decode http response to json failed: %v", err)
		if glog.V(4) {
			glog.Infof(string(reader.Bytes()))
		}
		return nil, err
	}
	return response.Repositories, nil
}

type response struct {
	Repositories []model.RepositoryName `json:"repositories"`
}
