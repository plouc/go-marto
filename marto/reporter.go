package marto

import (
	"net/http"
	"time"
	"strconv"
	"io"
	"fmt"
	"encoding/json"
)

// Interface for a basic reporter
type Reporter interface {
	// scenario related callbacks
	OnScenarioStarted(scenario *Scenario)
	OnScenarioFinished(scenario *Scenario)

	// session related callbacks
	OnSessionStarted(session *Session)
	OnSessionFinished(session *Session)

	// request/response related callbacks
	OnRequest(session *Session, request *http.Request)
	OnResponse(session *Session, request *http.Request, response *http.Response)
}

// Stats for a single scenario
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

// Stats for a single request
type RequestStat struct {
	StartedAt time.Time
	EndedAt   time.Time
	Duration  time.Duration
}

// Stats for a group of similar requests
type RequestStats struct {
	Id      string `json:"id"`
	Method  string `json:"method"`
	Url     string `json:"url"`
	Count   int    `json:"count"`
	Total   int64  `json:"total"`
	Average int64  `json:"average"`
}

func computeRequestId(request *http.Request) string {
	return request.Method + "." + request.URL.String()
}

// Base reporter
type BaseReporter struct {

	StartedAt time.Time

	Writers        []io.Writer
	ScenariosStats map[string]*ScenarioStats
	RequestStats   map[*http.Request]*RequestStat

	// Used to register callbacks
	OnScenarioStartedFn  func(scenario *Scenario)
	OnScenarioFinishedFn func(scenario *Scenario)
	OnSessionStartedFn   func(session *Session)
	OnSessionFinishedFn  func(session *Session)
	OnRequestFn          func(session *Session, request *http.Request)
	OnResponseFn         func(session *Session, request *http.Request, response *http.Response, stat *RequestStats)
}

func NewBaseReporter() *BaseReporter {
	return &BaseReporter{
		Writers:        make([]io.Writer, 0),
		ScenariosStats: map[string]*ScenarioStats{},
		RequestStats:   map[*http.Request]*RequestStat{},
	}
}

func (r *BaseReporter) AddWriter(writer io.Writer) {
	r.Writers = append(r.Writers, writer)
}

// Called when a scenario started
func (r *BaseReporter) OnScenarioStarted(scenario *Scenario) {
	if len(r.Writers) > 0 {
		for _, writer := range r.Writers {
			fmt.Fprintf(writer, "scenario.started | %s [%s]\n", scenario.Id, time.Since(r.StartedAt))
		}
	}

	r.ScenariosStats[scenario.Id] = &ScenarioStats{
		Histogram: make(map[string]int, 0),
		Stats:     map[string]*RequestStats{},
	}

	if r.OnScenarioStartedFn != nil {
		r.OnScenarioStartedFn(scenario)
	}
}

// Called when a scenario finished
func (r *BaseReporter) OnScenarioFinished(scenario *Scenario) {
	if len(r.Writers) > 0 {
		for _, writer := range r.Writers {
			fmt.Fprintf(writer, "scenario.finished | %s [%s]\n", scenario.Id, time.Since(r.StartedAt))
		}
	}

	if r.OnScenarioFinishedFn != nil { r.OnScenarioFinishedFn(scenario)	}
}

// Called when a session started
func (r *BaseReporter) OnSessionStarted(session *Session) {
	if len(r.Writers) > 0 {
		for _, writer := range r.Writers {
			fmt.Fprintf(writer, "session.started | scenario \"%s\" [%d/%d] [%s]\n", session.Scenario.Id, session.Id() + 1, session.Scenario.RepeatCount(), time.Since(r.StartedAt))
		}
	}

	if r.OnSessionStartedFn != nil { r.OnSessionStartedFn(session) }
}

// Called when a session finished
func (r *BaseReporter) OnSessionFinished(session *Session) {
	if len(r.Writers) > 0 {
		for _, writer := range r.Writers {
			fmt.Fprintf(writer, "%s.session.finished [%d of %d] [%s]\n", session.Scenario.Id, session.Id(), session.Scenario.RepeatCount(), time.Since(r.StartedAt))
		}
	}

	if r.OnSessionFinishedFn != nil { r.OnSessionFinishedFn(session) }
}

// Called when a request started:
// * check if there are available writers, and write if so
// * create a new requestStats for the given request and mark start time
// * check if there is an available callback function and call it if so
func (r *BaseReporter) OnRequest(session *Session, request *http.Request) {
	if len(r.Writers) > 0 {
		for _, writer := range r.Writers {
			fmt.Fprintf(writer, "-> %s %s [%d of %d] [%s]\n", request.Method, request.URL.String(), session.Current + 1, session.Scenario.RequestCount(), "n/a")
		}
	}

	r.RequestStats[request] = &RequestStat{
		StartedAt: time.Now(),
	}

	scenId := session.Scenario.Id
	reqId  := strconv.Itoa(session.Current)

	r.ScenariosStats[scenId].UpdateReqCount(1)
	if _, ok := r.ScenariosStats[scenId].Stats[reqId]; !ok {
		r.ScenariosStats[scenId].Stats[reqId] = &RequestStats{
			Id:     reqId,
			Method: request.Method,
			Url:    request.URL.String(),
		}
	}

	if r.OnRequestFn != nil { r.OnRequestFn(session, request) }
}

// Called when a request finished
func (r *BaseReporter) OnResponse(session *Session, request *http.Request, response *http.Response) {

	r.RequestStats[request].EndedAt  = time.Now()
	r.RequestStats[request].Duration = r.RequestStats[request].EndedAt.Sub(r.RequestStats[request].StartedAt)

	if len(r.Writers) > 0 {
		for _, writer := range r.Writers {
			fmt.Fprintf(writer, "<- %d | %s | %s in %s\n", response.StatusCode, request.Method, request.URL.String(), r.RequestStats[request].Duration)
		}
	}

	scenId := session.Scenario.Id
	reqId  := strconv.Itoa(session.Current)

	r.ScenariosStats[scenId].UpdateReqCount(-1)
	r.ScenariosStats[scenId].Stats[reqId].Count++
	r.ScenariosStats[scenId].Stats[reqId].Total += r.RequestStats[request].Duration.Nanoseconds()
	r.ScenariosStats[scenId].Stats[reqId].Average = r.ScenariosStats[scenId].Stats[reqId].Total / int64(r.ScenariosStats[scenId].Stats[reqId].Count)

	// if callback not nil, call it
	if r.OnResponseFn != nil {
		r.OnResponseFn(session, request, response, r.ScenariosStats[scenId].Stats[reqId])
	}
}

// Dump currently collected metrics
func (r *BaseReporter) Dump(w io.Writer) {
	for scenarioId, scenarioStats := range r.ScenariosStats {
		fmt.Fprintf(w, "+ scenario \"%s\":\n", scenarioId)
		for _, reqStats := range scenarioStats.Stats {
			fmt.Fprintf(w, "  - %s %s\n", reqStats.Method, reqStats.Url)
			fmt.Fprintf(w, "    %d request(s) | average: %d\n", reqStats.Count, reqStats.Average)
		}
		/*
		for ts, count := range scenarioStats.Histogram {
			fmt.Fprintf(w, "    %s | %d\n", ts, count)
		}
		*/
	}
}

func (r *BaseReporter) DumpJson(w io.Writer) {
	serialized, _ := json.MarshalIndent(r.ScenariosStats, "", "   ")

	w.Write(serialized)
}

