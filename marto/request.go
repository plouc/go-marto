package marto

import (
	"net/http"
	"io"
)

type BasicAuth struct {
	username string `json:"username"schema:"username"`
	password string `json:"password"schema:"password"`
}

type RequestTemplate struct {
	Method    string            `json:"method"schema:"method"`
	Url       string            `json:"url"schema:"url"`
	Body      io.Reader         `json:"body"schema:"body"`  
	delay     uint64            `json:"delay"schema:"delay"`
	Headers   map[string]string `json:"headers"schema:"headers"`
	BasicAuth *BasicAuth        `json:"basic_auth"schema:"basic_auth"`
}

func NewRequestTemplate(method string, url string, body io.Reader) *RequestTemplate {
	return &RequestTemplate{
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   make(map[string]string),
		BasicAuth: nil,
	}
}

func (rt *RequestTemplate) SetDelay(delay uint64) *RequestTemplate {
	rt.delay = delay

	return rt
}

func (rt *RequestTemplate) HasDelay() bool {
	return rt.delay > 0
}

func (rt *RequestTemplate) Delay() uint64 {
	return rt.delay
}

func (rt *RequestTemplate) AddHeader(key string, value string) {
	rt.Headers[key] = value
}

func (rt *RequestTemplate) SetBasicAuth(username string, password string) {
	rt.BasicAuth = &BasicAuth{username, password}
}

func BuildRequest(tpl *RequestTemplate) *http.Request {
	req, err := http.NewRequest(tpl.Method, tpl.Url, tpl.Body)
    if err != nil {
    	panic(err)
    }

    if tpl.BasicAuth != nil {
    	req.SetBasicAuth(tpl.BasicAuth.username, tpl.BasicAuth.password)
    }

    for hKey, hValue := range tpl.Headers {
		req.Header.Set(hKey, hValue)
	}

    return req
}