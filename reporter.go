package marto

import (
	"fmt"
	"net/http"
)

type Reporter interface {
	OnScenarioStarted(scenario *Scenario)
	OnScenarioFinished(scenario *Scenario)

	OnSessionStarted(session *Session)
	OnSessionFinished(session *Session)

	OnRequest(request *Request)
	OnResponse(response *http.Response)
}

type SimpleReporter struct {
}

func (sr *SimpleReporter) OnScenarioStarted(scenario *Scenario) {
	fmt.Printf("scenario.started | %s\n", scenario.Id)
}
func (sr *SimpleReporter) OnScenarioFinished(scenario *Scenario) {
	fmt.Printf("scenario.finished | %s\n", scenario.Id)
}

func (sr *SimpleReporter) OnSessionStarted(session *Session) {
	fmt.Printf("session.started | scenario \"%s\" [%d/%d]\n", session.Scenario.Id, session.Id() + 1, session.Scenario.RepeatCount())
}
func (sr *SimpleReporter) OnSessionFinished(session *Session) {
	fmt.Println("session.finished")
}

func (sr *SimpleReporter) OnRequest(request *Request) {
	fmt.Printf("request | %s %s [step %d/%d]\n", request.Method, request.URL.String(), request.Id() + 1, request.Scenario.RequestCount())
}
func (sr *SimpleReporter) OnResponse(response *http.Response) {
	fmt.Printf("response | %d | %s %s\n", response.StatusCode, response.Request.Method, response.Request.URL.String())
}