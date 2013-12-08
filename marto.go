package marto

import (
	"fmt"
	"net/http"
	//"io/ioutil"
	"log"
	"time"
)

type Marto struct {
	client                 *http.Client
	RequestStats           []*RequestStat
	AggregatedRequestStats map[string]*AggregatedRequestStat
	scenarios              map[string]*Scenario
	runningScenarios       []string
}

func NewMarto() *Marto {

	client := &http.Client{
	}

	return &Marto{
		client:                 client,
		RequestStats:           make([]*RequestStat, 0),
		AggregatedRequestStats: map[string]*AggregatedRequestStat{},
		scenarios:              map[string]*Scenario{},               
		runningScenarios:       make([]string, 0),
	}
}

func (m *Marto) AddScenario(id string, s *Scenario) {
	m.scenarios[id] = s
}

func (m *Marto) IsRunningScenario(id string) bool {
    for _, runningId := range m.runningScenarios {
        if id == runningId {
            return true
        }
    }
    return false
}

func (m *Marto) AggregateRequestStats() {
	for _, reqStat := range m.RequestStats {
		
		statKey := reqStat.Method + reqStat.Url

		if _, ok := m.AggregatedRequestStats[statKey]; !ok {
			m.AggregatedRequestStats[statKey] = &AggregatedRequestStat{
				Method:          reqStat.Method,
				Url:             reqStat.Url,
				Count:           0,
				Total:           0,
				AverageDuration: 0,
			}
		}

		m.AggregatedRequestStats[statKey].Count++
		m.AggregatedRequestStats[statKey].Total += reqStat.Duration.Nanoseconds()
		m.AggregatedRequestStats[statKey].AverageDuration = m.AggregatedRequestStats[statKey].Total / m.AggregatedRequestStats[statKey].Count
	}
}

func (m *Marto) Start(id string) {

	fmt.Println("starting scenario", id)

	if _, ok := m.scenarios[id]; !ok {
		panic(fmt.Sprintf("No scenario defined for id: %s", id))
	}

	scenario := m.scenarios[id]

	if !scenario.HasRequest() {
		panic(fmt.Sprintf("No request found on scenario: %s", id))	
	}

	session := scenario.CreateSession()
	m.processSession(session)

	if scenario.ShouldRepeat() {
		for i := 1; i < scenario.RepeatCount(); i++ {
			session = scenario.CreateSession()
			m.processSession(session)
		}
	}

	fmt.Println("finished scenario", id)
}

func (m *Marto) processSession(s *Session) {
	if !s.HasFinished() {
		req := s.CurrentRequest()

		if req.HasDelay() {
			fmt.Printf("...delaying execution of %s %s for %dms...\n", req.Method, req.URL.String(), req.Delay())
			select {
        	case <-time.After(time.Duration(req.Delay() * uint64(time.Millisecond))):
        		m.doRequest(req)
				s.Next()
				m.processSession(s)
        	}
		} else {
			m.doRequest(req)
			s.Next()
			m.processSession(s)
		}
	}
}

func (m *Marto) doRequest(req *Request) *http.Response {

	reqStat := &RequestStat{
		Url:       req.URL.String(),
		Method:    req.Method,
		StartedAt: time.Now(),
	}

	m.RequestStats = append(m.RequestStats, reqStat)

	fmt.Printf("-> request.start - %s [%d]\n", req.URL.String(), len(m.RequestStats))

	defer reqStat.Finished()

	res, err := m.client.Do(req.Request)
	if err != nil {
		log.Fatal(err)
	}

	reqStat.StatusCode = res.StatusCode

	fmt.Printf("   request.end [%d]\n", res.StatusCode)

	return res
}