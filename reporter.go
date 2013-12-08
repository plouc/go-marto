package marto

import (
	"net/http"
)

type Reporter interface {
	OnScenarioStarted(scenario *Scenario)
	OnScenarioFinished(scenario *Scenario)

	OnSessionStarted(session *Session)
	OnSessionFinished(session *Session)

	OnRequest(request *Request)
	OnResponse(request *Request, response *http.Response)
}