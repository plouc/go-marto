package marto

import (
	"io"
)

type Scenario struct {
	requests []*Request
	repeat   int
	sessions []*Session
}

func NewScenario() *Scenario {
	return &Scenario {
		requests: make([]*Request, 0),
		repeat:   0,
	}
}

func (s *Scenario) CreateSession() *Session {
	session := NewSession(s.requests)

	s.sessions = append(s.sessions, session)

	return session
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

func (s *Scenario) AppendRequest(req *Request) {
	s.requests = append(s.requests, req)
}

func (s *Scenario) AppendRequestFromConfig(method string, strUrl string, body io.Reader, delay uint64) {
	req, err := NewRequest(method, strUrl, body, delay)
    if err != nil {
    	panic(err)
    }
    s.AppendRequest(req)
}

func (s *Scenario) HasRequest() bool {
	return len(s.requests) > 0
}

