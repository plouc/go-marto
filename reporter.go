package marto

import (
	"fmt"
	"net/http"
	"io"
	//"time"
)

type Reporter interface {
	OnScenarioStarted(scenario *Scenario)
	OnScenarioFinished(scenario *Scenario)

	OnSessionStarted(session *Session)
	OnSessionFinished(session *Session)

	OnRequest(request *Request)
	OnResponse(request *Request, response *http.Response)
}

type WriterReporter struct {
	Writer io.Writer
}

func NewWriterReporter(writer io.Writer) *WriterReporter {
	return &WriterReporter{writer}
}

func (wr *WriterReporter) OnScenarioStarted(scenario *Scenario) {
	fmt.Printf("scenario.started | %s\n", scenario.Id)
}
func (wr *WriterReporter) OnScenarioFinished(scenario *Scenario) {
	fmt.Printf("scenario.finished | %s\n", scenario.Id)
}

func (wr *WriterReporter) OnSessionStarted(session *Session) {
	fmt.Printf("session.started | scenario \"%s\" [%d/%d]\n", session.Scenario.Id, session.Id() + 1, session.Scenario.RepeatCount())
}
func (wr *WriterReporter) OnSessionFinished(session *Session) {
	fmt.Println("session.finished")
}

func (wr *WriterReporter) OnRequest(request *Request) {
	fmt.Printf("request | %s %s [step %d/%d]\n", request.Method, request.URL.String(), request.Id() + 1, request.Scenario.RequestCount())
}
func (wr *WriterReporter) OnResponse(request *Request, response *http.Response) {
	fmt.Printf("response | %d | %s |%s %s\n", response.StatusCode, request.Elapsed, request.Method, request.URL.String())
}