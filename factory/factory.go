package factory

import (
	"github.com/bborbe/docker_utils/repositories"
	http_client_builder "github.com/bborbe/http/client_builder"
	"net/http"
)

type dockerUtilsFactory struct{}

func New() *dockerUtilsFactory {
	return new(dockerUtilsFactory)
}

func (d *dockerUtilsFactory) Repositories() repositories.Repositories {
	return repositories.New(d.httpClient())
}

func (d *dockerUtilsFactory) httpClient() *http.Client {
	return http_client_builder.New().WithoutProxy().Build()
}
