package marto

type Session struct {
	id       int // private because it must be set according to scenario.Sessions size
	current  int
	requests []*Request
	finished bool
	Scenario *Scenario
}

func NewSession(scenario *Scenario) *Session {
	return &Session{
		id:       len(scenario.Sessions()),
		current:  0,
		Scenario: scenario,
		requests: scenario.Requests(),
		finished: false,
	}
}

func (s *Session) Id() int {
	return s.id
}

func (s *Session) CurrentRequest() *Request {
	return s.requests[s.current]
}

func (s *Session) HasFinished() bool {
	return s.finished
}

func (s *Session) Next() {
	if s.current < (len(s.requests) - 1) {
		s.current++
	} else {
		s.finished = true
	}
}