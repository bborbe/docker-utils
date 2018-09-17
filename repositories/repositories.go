package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"sort"

	"github.com/bborbe/docker-utils/model"
	"github.com/bborbe/io/reader_shadow_copy"
	"github.com/golang/glog"
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
		return nil, errors.Wrap(err, "create http request failed")
	}
	if err := registry.SetAuth(req); err != nil {
		return nil, errors.Wrap(err, "set auth failed")
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "perform http request failed")
	}
	if resp.StatusCode/100 != 2 {
		return nil, errors.Wrapf(err, "request failed with status: %d", resp.StatusCode)
	}
	var response response
	reader := reader_shadow_copy.New(resp.Body)
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode http response to json failed")
	}
	if glog.V(4) {
		glog.Infof(string(reader.Bytes()))
	}
	repositories := response.Repositories
	sort.Sort(model.RepositoryNamesByName(repositories))
	return repositories, nil
}

type response struct {
	Repositories []model.RepositoryName `json:"repositories"`
}
