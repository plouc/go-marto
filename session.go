package marto

type Session struct {
	current  int
	requests []*Request
	finished bool
}

func NewSession(requests []*Request) *Session {
	return &Session{
		current:  0,
		requests: requests,
		finished: false,
	}
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