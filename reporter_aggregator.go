package marto

import (
	"fmt"
	"net/http"
	"io"
	"os"
)

type AggregatorReporter struct {
	Scenarios map[string]map[int]*RequestStats
}

type RequestStats struct {
	Method  string
	Url     string
	Count   int
	Total   int64
	Average int64
}

func NewAggregatorReporter() *AggregatorReporter {
	return &AggregatorReporter{
		Scenarios: make(map[string]map[int]*RequestStats),
	}
}

func (r *AggregatorReporter) OnScenarioStarted(scenario *Scenario) {}
func (r *AggregatorReporter) OnScenarioFinished(scenario *Scenario) {}
func (r *AggregatorReporter) OnSessionStarted(session *Session) {}
func (r *AggregatorReporter) OnSessionFinished(session *Session) {}
func (r *AggregatorReporter) OnRequest(request *Request) {}


func (r *AggregatorReporter) Dump(w io.Writer) {
	for scenarioId, stats := range r.Scenarios {
		fmt.Fprintf(w, "+ scenario \"%s\":\n", scenarioId)
		for _, reqStats := range stats {
			fmt.Fprintf(w, "  - %s %s\n", reqStats.Method, reqStats.Url)
			fmt.Fprintf(w, "    %d request(s) | average: %d\n", reqStats.Count, reqStats.Average)
		}
	}
}


func (r *AggregatorReporter) OnResponse(request *Request, response *http.Response) {

	if _, ok := r.Scenarios[request.Scenario.Id]; !ok {
    	r.Scenarios[request.Scenario.Id] = make(map[int]*RequestStats)
	}

	if _, ok := r.Scenarios[request.Scenario.Id][request.Id()]; !ok {
		r.Scenarios[request.Scenario.Id][request.Id()] = &RequestStats{
			Method:  request.Method,
			Url:     request.URL.String(),
			Count:   0,
			Total:   0,
			Average: 0,
		}
	}

	r.Scenarios[request.Scenario.Id][request.Id()].Count++
	r.Scenarios[request.Scenario.Id][request.Id()].Total += request.Elapsed.Nanoseconds()
	r.Scenarios[request.Scenario.Id][request.Id()].Average = r.Scenarios[request.Scenario.Id][request.Id()].Total / int64(r.Scenarios[request.Scenario.Id][request.Id()].Count)

	r.Dump(os.Stdout)
}