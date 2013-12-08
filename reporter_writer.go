package marto

import (
	"fmt"
	"net/http"
	"io"
)

type WriterReporter struct {
	Writer io.Writer
}

func NewWriterReporter(writer io.Writer) *WriterReporter {
	return &WriterReporter{writer}
}

func (r *WriterReporter) OnScenarioStarted(scenario *Scenario) {
	fmt.Fprintf(r.Writer, "scenario.started | %s\n", scenario.Id)
}
func (r *WriterReporter) OnScenarioFinished(scenario *Scenario) {
	fmt.Fprintf(r.Writer, "scenario.finished | %s\n", scenario.Id)
}

func (r *WriterReporter) OnSessionStarted(session *Session) {
	fmt.Fprintf(r.Writer, "session.started | scenario \"%s\" [%d/%d]\n", session.Scenario.Id, session.Id() + 1, session.Scenario.RepeatCount())
}
func (r *WriterReporter) OnSessionFinished(session *Session) {
	fmt.Fprintf(r.Writer, "session.finished\n")
}

func (r *WriterReporter) OnRequest(request *Request) {
	fmt.Fprintf(r.Writer, "request | %s %s [step %d/%d]\n", request.Method, request.URL.String(), request.Id() + 1, request.Scenario.RequestCount())
}
func (r *WriterReporter) OnResponse(request *Request, response *http.Response) {
	fmt.Fprintf(r.Writer, "response | %d | %s | %s %s\n", response.StatusCode, request.Elapsed, request.Method, request.URL.String())
}