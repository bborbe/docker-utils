package factory

import (
	"net/http"

	"github.com/bborbe/docker-utils/repositories"
	"github.com/bborbe/docker-utils/tags"
	http_client_builder "github.com/bborbe/http/client_builder"
)

type dockerUtilsFactory struct{}

func New() *dockerUtilsFactory {
	return new(dockerUtilsFactory)
}

func (d *dockerUtilsFactory) Repositories() repositories.Repositories {
	return repositories.New(d.httpClient())
}

func (d *dockerUtilsFactory) Tags() tags.Tags {
	return tags.New(d.httpClient())
}

func (d *dockerUtilsFactory) httpClient() *http.Client {
	return http_client_builder.New().WithoutProxy().Build()
}
