package main

import (
	"testing"

	"bytes"
	. "github.com/bborbe/assert"
)

func TestResumeFail(t *testing.T) {
	writer := &bytes.Buffer{}
	if err := AssertThat(do(writer), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
