package marto

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"time"
	"strconv"
)

type AggregatorReporter struct {
	Scenarios map[string]*ScenarioStats `json:"scenarios"`
}

type ScenarioStats struct {
	requestCount int
	Histogram    map[string]int           `json:"histogram"`
	Stats        map[string]*RequestStats `json:"stats"`
}

func (s *ScenarioStats) UpdateReqCount(count int) {
	s.requestCount += count

	ts := strconv.Itoa(int(time.Now().Unix())) // 60

	if _, ok := s.Histogram[ts]; ok {
		if count > 0 {
			s.Histogram[ts] += count
		}
	} else {
		s.Histogram[ts] = s.requestCount
	}
}

type RequestStats struct {
	Method  string `json:"method"`
	Url     string `json:"url"`
	Count   int    `json:"count"`
	Total   int64  `json:"total"`
	Average int64  `json:"average"`
}

func NewAggregatorReporter() *AggregatorReporter {
	return &AggregatorReporter{
		Scenarios: map[string]*ScenarioStats{},
	}
}

func (r *AggregatorReporter) OnScenarioStarted(scenario *Scenario) {
	r.Scenarios[scenario.Id] = &ScenarioStats{
		Histogram: make(map[string]int),
		Stats:     map[string]*RequestStats{},
	}
}

func (r *AggregatorReporter) OnScenarioFinished(scenario *Scenario) {}
func (r *AggregatorReporter) OnSessionStarted(session *Session) {}
func (r *AggregatorReporter) OnSessionFinished(session *Session) {}

// Called when a request starts
func (r *AggregatorReporter) OnRequest(request *Request) {
	r.Scenarios[request.Scenario.Id].UpdateReqCount(1)
	if _, ok := r.Scenarios[request.Scenario.Id].Stats[strconv.Itoa(request.Id())]; !ok {
		r.Scenarios[request.Scenario.Id].Stats[strconv.Itoa(request.Id())] = &RequestStats{
			Method: request.Method,
			Url:    request.URL.String(),
		}
	}
}

// Dump currently collected metrics
func (r *AggregatorReporter) Dump(w io.Writer) {
	for scenarioId, scenarioStats := range r.Scenarios {
		fmt.Fprintf(w, "+ scenario \"%s\":\n", scenarioId)
		for _, reqStats := range scenarioStats.Stats {
			fmt.Fprintf(w, "  - %s %s\n", reqStats.Method, reqStats.Url)
			fmt.Fprintf(w, "    %d request(s) | average: %d\n", reqStats.Count, reqStats.Average)
		}
		for ts, count := range scenarioStats.Histogram {
			fmt.Fprintf(w, "    %s | %d\n", ts, count)
		}
	}
}

func (r *AggregatorReporter) DumpJson(w io.Writer) {
	serialized, _ := json.MarshalIndent(r.Scenarios, "", "   ")

	w.Write(serialized)
}


// Called on request response
func (r *AggregatorReporter) OnResponse(request *Request, response *http.Response) {
	r.Scenarios[request.Scenario.Id].UpdateReqCount(-1)

	reqId  := strconv.Itoa(request.Id())
	scenId := request.Scenario.Id

	r.Scenarios[scenId].Stats[reqId].Count++
	r.Scenarios[scenId].Stats[reqId].Total += request.Elapsed.Nanoseconds()
	r.Scenarios[scenId].Stats[reqId].Average = r.Scenarios[scenId].Stats[reqId].Total / int64(r.Scenarios[scenId].Stats[reqId].Count)
}