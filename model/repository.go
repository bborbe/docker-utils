package model

type Tag string

type RepositoryName string

func (r RepositoryName) String() string {
	return string(r)
}

type Repository struct {
	Name RepositoryName
	Tag  Tag
}
