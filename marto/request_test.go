package marto

import (
	"testing"
	"net/http"
	"github.com/stretchr/testify/assert"
	//"fmt"
)


func TestBuildRequest(t *testing.T) {
	tpl := NewRequestTemplate("GET", "localhost", nil)
	tpl.AddHeader("Content-Type", "text/plain")
	tpl.SetBasicAuth("user", "pass")

	req := BuildRequest(tpl)

	assert.IsType(t, new(http.Request), req)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "text/plain", req.Header.Get("Content-Type"))
	assert.Equal(t, "Basic dXNlcjpwYXNz", req.Header.Get("Authorization"))
}