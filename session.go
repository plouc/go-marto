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

func (s *Session) ConsumeRequest() *Request {
	req := s.requests[s.current]
	req.Session = s

	if s.current < (len(s.requests) - 1) {
		s.current++
	} else {
		s.finished = true
	}

	return req
}

func (s *Session) HasFinished() bool {
	return s.finished
}