package marto

import (
	"fmt"
	"net/http"
	"log"
	"time"
)

type Marto struct {
	client              *http.Client
	scenarios           map[string]*Scenario
	pendingSessionCount int
	reporters           []Reporter
}

func NewMarto() *Marto {
	return &Marto{
		client:              &http.Client{},
		scenarios:           map[string]*Scenario{},               
		pendingSessionCount: 0,
	}
}

// AddReporter add a new reporter to the list of reporters.
// You should use the WriterReporter for example.
func (m *Marto) AddReporter(reporter Reporter) {
	m.reporters = append(m.reporters, reporter)
}

// AddScenario add a new scenario to the list of scenarios.
func (m *Marto) AddScenario(scenario *Scenario) {
	m.scenarios[scenario.Id] = scenario
}

var ch = make(chan *Session)

// Run all available scenarios.
func (m *Marto) Run() {
	for _, scenario := range m.scenarios {
		m.pendingSessionCount += scenario.RepeatCount()
		m.runScenario(scenario)
	}

	for {
		select {
		case _ = <-ch:
			m.pendingSessionCount--
			if m.pendingSessionCount == 0 {
				return
			}
		}
	}

}

// run the scenario
func (m *Marto) runScenario(scenario *Scenario) {
	if !scenario.HasRequest() {
		panic(fmt.Sprintf("No request found on scenario: %s", scenario.Id))	
	}

	for _, reporter := range m.reporters {
		reporter.OnScenarioStarted(scenario)
	}

	for i := 0; i < scenario.RepeatCount(); i++ {
		session := scenario.CreateSession()

		for _, reporter := range m.reporters {
			reporter.OnSessionStarted(session)
		}

		go func() {
			m.processSession(session)
		}()
	}
}

// send current session request
func (m *Marto) processSession(sess *Session) {
	if !sess.HasFinished() {
		req := sess.ConsumeRequest()
		
		delay := int(req.Delay() * uint64(time.Millisecond))
		if req.IsFirst() {
			delay += sess.Scenario.GetDelay() * int(time.Millisecond) * sess.Id()
		}

		select {
        case <-time.After(time.Duration(delay)):
        	m.doRequest(req)
        }
	} else {
		ch <- sess
	}
}

// send a request
func (m *Marto) doRequest(req *Request) {

	for _, reporter := range m.reporters { reporter.OnRequest(req) }

	req.Start()

	res, err := m.client.Do(req.Request)
	if err != nil {
		log.Fatal(err)
	}

	req.End()

	for _, reporter := range m.reporters { reporter.OnResponse(req, res) }

	m.processSession(req.Session)
}