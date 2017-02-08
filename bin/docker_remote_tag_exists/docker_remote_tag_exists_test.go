package main

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestResumeFail(t *testing.T) {
	if err := AssertThat(do(), NilValue()); err != nil {
		t.Fatal(err)
	}
}
