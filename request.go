package marto

import (
	"net/http"
	"io"
)

type Request struct {
	*http.Request
	delay uint64
}

func NewRequest(method string, strUrl string, body io.Reader, delay uint64) (*Request, error) {

	req, err := http.NewRequest(method, strUrl, body)
	if err != nil {
		return nil, err
	}

	return &Request{req, delay}, nil
}

func (r *Request) HasDelay() bool {
	return r.delay > 0
}

func (r *Request) Delay() uint64 {
	return r.delay
}