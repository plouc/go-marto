package marto

import "net/http"

type Session struct {
	id       int // private because it must be set according to scenario.Sessions size
	Current  int
	finished bool
	Scenario *Scenario
}

func NewSession(scenario *Scenario) *Session {
	return &Session{
		id:       len(scenario.Sessions()),
		Current:  -1,
		Scenario: scenario,
		finished: false,
	}
}

func (s *Session) Id() int {
	return s.id
}

func (s *Session) Request() (*http.Request, *RequestTemplate) {
	if s.finished {
		panic("Cannot get request from finished session")
	}

	if s.Current < s.Scenario.RequestCount()-1 {
		s.Current++
	}

	if s.Current == s.Scenario.RequestCount()-1 {
		s.finished = true
	}

	tpl := s.Scenario.Request(s.Current)
	req := BuildRequest(tpl)

	return req, tpl
}

func (s *Session) HasFinished() bool {
	return s.finished
}