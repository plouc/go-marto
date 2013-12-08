package marto

import (
	"io"
)

type Scenario struct {
	Id       string
	requests []*Request
	repeat   int
	sessions []*Session
}

func NewScenario(id string) *Scenario {
	return &Scenario {
		Id:       id,
		requests: make([]*Request, 0),
		repeat:   0,
	}
}

func (s *Scenario) CreateSession() *Session {
	session := NewSession(s)

	s.sessions = append(s.sessions, session)

	return session
}

func (s *Scenario) Sessions() []*Session {
	return s.sessions
}

func (s *Scenario) Requests() []*Request {
	return s.requests
}

func (s *Scenario) Repeat(count int) {
	s.repeat = count
}

func (s *Scenario) ShouldRepeat() bool {
	return s.repeat > 0
}

func (s *Scenario) RepeatCount() int {
	return s.repeat
}

func (s *Scenario) RequestCount() int {
	return len(s.requests)
}

func (s *Scenario) Append(method string, strUrl string, body io.Reader) *Request {
	req, err := NewRequest(s.RequestCount(), s, method, strUrl, body)
    if err != nil {
    	panic(err)
    }

    s.requests = append(s.requests, req)

    return req
}

func (s *Scenario) HasRequest() bool {
	return len(s.requests) > 0
}

