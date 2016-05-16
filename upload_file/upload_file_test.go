package upload_file

import (
	"net/http"
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsHandler(t *testing.T) {
	r := New("/tmp")
	var i *http.Handler
	err := AssertThat(r, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}

func TestTarget(t *testing.T) {
	d, err := target("/tmp", "/foo")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(d, Is("/tmp/foo")); err != nil {
		t.Fatal(err)
	}
}

func TestTargetError(t *testing.T) {
	_, err := target("/tmp", "../foo")
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
