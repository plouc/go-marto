package marto

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
  	}))
	defer ts.Close()

	m := NewMarto()

	scenario := NewScenario("test_0")
	scenario.Append("GET", ts.URL, nil)
	m.AddScenario(scenario)

	reporter := NewAggregatorReporter()
	m.AddReporter(reporter)

	m.Run()

	reports := reporter.Scenarios

	assert.Equal(t, len(reports), 1, "reports contain one scenario")
	_, ok := reports["test_0"]
	assert.True(t, ok, "reports contain the given scenario")
}