package docker

import (
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type Registry struct {
	Url      string
	Username string
	Password string
}

func (r *Registry) RegistryPasswordFromFile(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "read file failed")
	}
	r.Password = strings.TrimSpace(string(content))
	return nil
}
