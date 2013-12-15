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
		Current:  0,
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

	tpl := s.Scenario.Request(s.Current)
	req, err := http.NewRequest(tpl.Method, tpl.Url, tpl.Body)
    if err != nil {
    	panic(err)
    }

	if s.Current < s.Scenario.RequestCount()-1 {
		s.Current++
	} else {
		s.finished = true
	}

	return req, tpl
}

func (s *Session) HasFinished() bool {
	return s.finished
}