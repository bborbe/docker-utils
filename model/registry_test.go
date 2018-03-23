package model

import (
	"testing"
	"bytes"
	. "github.com/bborbe/assert"
)

func TestCredentialsFromDockerConfigQuay(t *testing.T) {
	b := bytes.NewBufferString(`{"auths": {"quay.io": {"auth": "aGVsbG86d29ybGQ="}}}`)
	r := &Registry{
		Name: "quay.io",
	}
	err := r.CredentialsFromDockerConfig(b)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(r.Username.String(), Is("hello")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(r.Password.String(), Is("world")); err != nil {
		t.Fatal(err)
	}
}

func TestCredentialsFromDockerConfigDockerHub(t *testing.T) {
	b := bytes.NewBufferString(`{"auths": {"https://index.docker.io/v1/": {"auth": "aGVsbG86d29ybGQ="}}}`)
	r := &Registry{
		Name: "docker.io",
	}
	err := r.CredentialsFromDockerConfig(b)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(r.Username.String(), Is("hello")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(r.Password.String(), Is("world")); err != nil {
		t.Fatal(err)
	}
}


