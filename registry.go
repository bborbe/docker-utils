// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
