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

		delay := scenario.GetDelay() * int(time.Millisecond) * i

		select {
        case <-time.After(time.Duration(delay)):
        	sess := scenario.CreateSession()
		    for _, reporter := range m.reporters {
		    	reporter.OnSessionStarted(sess)
		    }
		    go func() { m.processSession(sess) }()
        }
	}
}

// send current session request
func (m *Marto) processSession(session *Session) {
	if session.HasFinished() == false {
		req, tpl := session.Request()
		
		delay := int(tpl.Delay() * uint64(time.Millisecond))
		if session.Current == 0 {
			delay += session.Scenario.GetDelay() * int(time.Millisecond) * session.Id()
		}

		select {
        case <-time.After(time.Duration(delay)):
        	m.doSessionRequest(session, req)
        }
	} else {
		for _, reporter := range m.reporters {
			reporter.OnSessionFinished(session)
		}
		ch <- session
	}
}

// send a request
func (m *Marto) doSessionRequest(session *Session, req *http.Request) {

	for _, reporter := range m.reporters {
		reporter.OnRequest(session, req)
	}

	res, err := m.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	for _, reporter := range m.reporters {
		reporter.OnResponse(session, req, res)
	}

	m.processSession(session)
}