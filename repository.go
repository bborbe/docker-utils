// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

type RepositoryName string

func (r RepositoryName) String() string {
	return string(r)
}

type RepositoryNamesByName []RepositoryName

func (t RepositoryNamesByName) Len() int {
	return len(t)
}

func (t RepositoryNamesByName) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t RepositoryNamesByName) Less(i, j int) bool {
	return t[i] < t[j]
}

type Repository struct {
	Name RepositoryName
	Tag  TagName
}
