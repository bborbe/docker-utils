package model

import (
	"errors"
	"fmt"
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
