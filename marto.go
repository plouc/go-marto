package marto

import (
	"fmt"
	"net/http"
	"log"
	"time"
)

var ch = make(chan *Session)

type Marto struct {
	client                 *http.Client
	RequestStats           []*RequestStat
	AggregatedRequestStats map[string]*AggregatedRequestStat
	scenarios              map[string]*Scenario
	runningScenarios       []string
	reporters              []Reporter
}

func NewMarto() *Marto {
	return &Marto{
		client:                 &http.Client{},
		RequestStats:           make([]*RequestStat, 0),
		AggregatedRequestStats: map[string]*AggregatedRequestStat{},
		scenarios:              map[string]*Scenario{},               
		runningScenarios:       make([]string, 0),
	}
}

func (m *Marto) AddReporter(reporter Reporter) {
	m.reporters = append(m.reporters, reporter)
}


func (m *Marto) AddScenario(scenario *Scenario) {
	m.scenarios[scenario.Id] = scenario
}

// run all scenarios
func (m *Marto) Run() {
	for _, scenario := range m.scenarios {
		m.RunScenario(scenario)
	}

	for {
		select {
		case _ = <-ch:
		}
	}
}


// run the scenario
func (m *Marto) RunScenario(scenario *Scenario) {
	if !scenario.HasRequest() {
		panic(fmt.Sprintf("No request found on scenario: %s", scenario.Id))	
	}

	for _, reporter := range m.reporters {
		reporter.OnScenarioStarted(scenario)
	}

	m.StartScenarioSession(scenario)
	if scenario.ShouldRepeat() {
		for i := 1; i < scenario.RepeatCount(); i++ {
			m.StartScenarioSession(scenario)
		}
	}
}


func (m *Marto) StartScenarioSession(scenario *Scenario) {
	session := scenario.CreateSession()

	for _, reporter := range m.reporters {
		reporter.OnSessionStarted(session)
	}

	go func() {
		m.processSession(session)
	}()
}


// send current session request
func (m *Marto) processSession(sess *Session) {
	if !sess.HasFinished() {
		req := sess.CurrentRequest()
		if req.HasDelay() {
			//fmt.Printf("...delaying execution of %s %s for %dms...\n", req.Method, req.URL.String(), req.Delay())
			select {
        	case <-time.After(time.Duration(req.Delay() * uint64(time.Millisecond))):
        		ch <- m.processSessionRequest(sess)
        	}
		} else {
			ch <- m.processSessionRequest(sess)
		}
	}
}

// send current session request and try to process next request
func (m *Marto) processSessionRequest(s *Session) *Session {
	req := s.CurrentRequest()
	req.Resolve()
	m.doRequest(req)
	s.Next()
	m.processSession(s)

	return s
}

// send a request
func (m *Marto) doRequest(req *Request) {

	for _, reporter := range m.reporters { reporter.OnRequest(req) }

	req.Start()
	//defer req.End()

	res, err := m.client.Do(req.Request)
	if err != nil {
		log.Fatal(err)
	}

	req.End()

	for _, reporter := range m.reporters { reporter.OnResponse(req, res) }
}