package marto

import (
	"io"
)

type ScenarioConfig struct {
	Id       string             `json:"id"`
	Repeat   int                `json:"repeat"`
	Every    int                `json:"every"`
	Requests []*RequestTemplate `json:"requests"`
}

type Scenario struct {
	Id               string
	requestTemplates []*RequestTemplate
	repeat           int
	every            int
	sessions         []*Session
}

func NewScenario(id string) *Scenario {
	return &Scenario {
		Id:               id,
		requestTemplates: make([]*RequestTemplate, 0),
		repeat:           1,
		every:            0,
		sessions:         make([]*Session, 0),
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

func (s *Scenario) Request(id int) *RequestTemplate {
	return s.requestTemplates[id]
}

func (s *Scenario) Requests() []*RequestTemplate {
	return s.requestTemplates
}

func (s *Scenario) Repeat(count int) *Scenario {
	if count < 1 {
		panic("You must repeat at least one time")
	}
	
	s.repeat = count

	return s
}

func (s *Scenario) ShouldRepeat() bool {
	return s.repeat > 0
}

func (s *Scenario) RepeatCount() int {
	return s.repeat
}

func (s *Scenario) Every(seconds int) *Scenario {
	s.every = seconds

	return s
}

func (s *Scenario) IsDelayed() bool {
	return s.every > 0
}

func (s *Scenario) GetDelay() int {
	return s.every
}

func (s *Scenario) RequestCount() int {
	return len(s.requestTemplates)
}

func (s *Scenario) HasRequest() bool {
	return s.RequestCount() > 0
}

func (s *Scenario) Append(method string, url string, body io.Reader) *RequestTemplate {

	reqTpl := NewRequestTemplate(method, url, body)

    s.requestTemplates = append(s.requestTemplates, reqTpl)

    return reqTpl
}

