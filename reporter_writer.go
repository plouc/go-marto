package marto

import (
	"fmt"
	"net/http"
	"io"
	"time"
)

type WriterReporter struct {
	Writer    io.Writer
	StartedAt time.Time
}

func NewWriterReporter(writer io.Writer) *WriterReporter {
	return &WriterReporter{
		Writer:    writer,
		StartedAt: time.Now(),
	}
}

func (r *WriterReporter) OnScenarioStarted(scenario *Scenario) {
	fmt.Fprintf(r.Writer, "scenario.started | %s [%s]\n", scenario.Id, time.Since(r.StartedAt))
}
func (r *WriterReporter) OnScenarioFinished(scenario *Scenario) {
	fmt.Fprintf(r.Writer, "scenario.finished | %s [%s]\n", scenario.Id, time.Since(r.StartedAt))
}

func (r *WriterReporter) OnSessionStarted(session *Session) {
	fmt.Fprintf(r.Writer, "session.started | scenario \"%s\" [%d/%d] [%s]\n", session.Scenario.Id, session.Id() + 1, session.Scenario.RepeatCount(), time.Since(r.StartedAt))
}
func (r *WriterReporter) OnSessionFinished(session *Session) {
	fmt.Fprintf(r.Writer, "session.finished [%s]\n", time.Since(r.StartedAt))
}

func (r *WriterReporter) OnRequest(request *Request) {
	fmt.Fprintf(r.Writer, "request | %s %s [step %d/%d] [%s]\n", request.Method, request.URL.String(), request.Id() + 1, request.Scenario.RequestCount(), time.Since(r.StartedAt))
}
func (r *WriterReporter) OnResponse(request *Request, response *http.Response) {
	fmt.Fprintf(r.Writer, "response | %d | %s | %s %s\n", response.StatusCode, request.Elapsed, request.Method, request.URL.String())
}