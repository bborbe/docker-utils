package model

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"github.com/bborbe/io/util"
	"os"
	"encoding/json"
	"encoding/base64"
	"io"
)

type RegistryUsername string

func (r RegistryUsername) String() string {
	return string(r)
}

func (r RegistryUsername) Validate() error {
	if len(r) == 0 {
		return errors.New("username empty")
	}
	return nil
}

type RegistryPassword string

func (r RegistryPassword) String() string {
	return string(r)
}

func (r RegistryPassword) Validate() error {
	if len(r) == 0 {
		return errors.New("password empty")
	}
	return nil
}

func RegistryPasswordFromFile(path string) (RegistryPassword, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return RegistryPassword(strings.TrimSpace(string(content))), nil
}

type RegistryName string

func (r RegistryName) String() string {
	return string(r)
}

func (r RegistryName) Url() string {
	return fmt.Sprintf("https://%s", r.String())
}

func (r RegistryName) Validate() error {
	if len(r) == 0 {
		return errors.New("registry empty")
	}
	return nil
}

type Registry struct {
	Name     RegistryName
	Username RegistryUsername
	Password RegistryPassword
}

func (r *Registry) ReadCredentialsFromDockerConfig() error {
	dockerConfig := "~/.docker/config.json"
	path, err := util.NormalizePath(dockerConfig)
	if err != nil {
		return fmt.Errorf("normalize path %s failed: %v", dockerConfig, err)
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file %s failed: %v", path, err)
	}
	return r.CredentialsFromDockerConfig(file)
}

func (r *Registry) CredentialsFromDockerConfig(reader io.Reader) error {
	var data struct {
		Domain map[string]struct {
			Auth string `json:"auth"`
		} `json:"auths"`
	}
	if err := json.NewDecoder(reader).Decode(&data); err != nil {
		return fmt.Errorf("decode json failed: %v", err)
	}
	auth, ok := data.Domain[nameToDomain(r.Name)];
	if !ok {
		return fmt.Errorf("domain %s not found in docker config", r.Name)
	}
	value, err := base64.StdEncoding.DecodeString(auth.Auth)
	if err != nil {
		return fmt.Errorf("base64 decode auth failed: %v", err)
	}
	parts := strings.SplitN(string(value), ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("split auth failed")
	}
	r.Username = RegistryUsername(parts[0])
	r.Password = RegistryPassword(parts[1])
	return nil
}

func nameToDomain(name RegistryName) string {
	if "docker.io" == name.String()  {
		return "https://index.docker.io/v1/"
	}
	return name.String()
}


func (r Registry) Validate() error {
	if err := r.Name.Validate(); err != nil {
		return err
	}
	if err := r.Username.Validate(); err != nil {
		return err
	}
	if err := r.Password.Validate(); err != nil {
		return err
	}
	return nil
}
