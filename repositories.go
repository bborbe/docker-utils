package docker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/bborbe/io/reader_shadow_copy"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type Repositories interface {
	List(registry Registry) ([]RepositoryName, error)
}

type repositories struct {
	httpClient *http.Client
}

func NewRepositories(httpClient *http.Client) Repositories {
	return &repositories{
		httpClient: httpClient,
	}
}

func (r *repositories) List(registry Registry) ([]RepositoryName, error) {
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
	sort.Sort(RepositoryNamesByName(repositories))
	return repositories, nil
}

type response struct {
	Repositories []RepositoryName `json:"repositories"`
}
