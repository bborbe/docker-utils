package requestbuilder

import (
	"testing"

	"io/ioutil"
	"strings"

	. "github.com/bborbe/assert"
)

func TestImplementsHttpRequestBuilder(t *testing.T) {
	r := NewHTTPRequestBuilder("http://www.example.com")
	var i *HttpRequestBuilder
	err := AssertThat(r, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetRequestWithHeader(t *testing.T) {
	r := NewHTTPRequestBuilder("http://www.benjamin-borbe.de")
	r.AddHeader("a", "b")
	request, err := r.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(request, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(request.Header), Is(1))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(request.Header["a"]), Is(1))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(request.Header["a"][0], Is("b"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetRequest(t *testing.T) {
	r := NewHTTPRequestBuilder("http://www.benjamin-borbe.de")
	request, err := r.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(request, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}

func TestDefaultMethodIsGet(t *testing.T) {
	r := NewHTTPRequestBuilder("http://www.benjamin-borbe.de")
	request, err := r.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(request.Method, Is("GET"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetMethod(t *testing.T) {
	r := NewHTTPRequestBuilder("http://www.benjamin-borbe.de")
	r.SetMethod("POST")
	request, err := r.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(request.Method, Is("POST"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetBody(t *testing.T) {
	r := NewHTTPRequestBuilder("http://www.benjamin-borbe.de")
	r.SetBody(strings.NewReader("hello world"))
	request, err := r.Build()
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(string(content), Is("hello world"))
	if err != nil {
		t.Fatal(err)
	}
}
