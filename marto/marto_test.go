package marto

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"fmt"
	"github.com/stretchr/testify/assert"
)


func TestRunOneScenarioWithOneRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
  	}))
	defer ts.Close()

	m := NewMarto()

	scenario := NewScenario("test_0")
	scenario.Append("GET", ts.URL, nil)
	m.AddScenario(scenario)

	reporter := NewBaseReporter()
	m.AddReporter(reporter)

	m.Run()

	scenariosStats := reporter.ScenariosStats

	assert.Equal(t, len(scenariosStats), 1)
	report, ok := scenariosStats["test_0"]
	assert.True(t, ok, "reports contain the given scenario")
	assert.Equal(t, len(report.Stats), 1, "reports stats contain one request")

	stats, ok := report.Stats["0"]
	assert.True(t, ok, "report contain a request with the given key")
	assert.IsType(t, new(RequestStats), stats)
	assert.Equal(t, stats.Url, ts.URL)
	assert.Equal(t, stats.Method, "GET")
	assert.Equal(t, stats.Count, 1)
	assert.Equal(t, len(report.Histogram), 1, "reports histogram contain one request")
}


func TestRunOneScenarioWithThreeRequests(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
  	}))
	defer ts.Close()

	m := NewMarto()

	scenario := NewScenario("test_0")
	scenario.Append("GET", ts.URL, nil)
	scenario.Append("GET", ts.URL, nil)
	scenario.Append("GET", ts.URL, nil)
	m.AddScenario(scenario)

	reporter := NewBaseReporter()
	m.AddReporter(reporter)

	m.Run()

	scenariosStats := reporter.ScenariosStats

	assert.Equal(t, len(scenariosStats), 1)
	report, ok := scenariosStats["test_0"]
	assert.True(t, ok, "reports contain the given scenario")
	assert.Equal(t, len(report.Stats), 3, "reports stats contain one request")

	stats, ok := report.Stats["0"]
	assert.True(t, ok)
	assert.IsType(t, new(RequestStats), stats)
	assert.Equal(t, stats.Method, "GET")
	assert.Equal(t, stats.Url, ts.URL)
	assert.Equal(t, stats.Count, 1)
	assert.Equal(t, len(report.Histogram), 1)

	stats, ok = report.Stats["1"]
	assert.True(t, ok)
	assert.IsType(t, new(RequestStats), stats)
	assert.Equal(t, stats.Method, "GET")
	assert.Equal(t, stats.Url, ts.URL)
	assert.Equal(t, stats.Count, 1)
	assert.Equal(t, len(report.Histogram), 1)

	stats, ok = report.Stats["2"]
	assert.True(t, ok)
	assert.IsType(t, new(RequestStats), stats)
	assert.Equal(t, stats.Method, "GET")
	assert.Equal(t, stats.Url, ts.URL)
	assert.Equal(t, stats.Count, 1)
	assert.Equal(t, len(report.Histogram), 1)
}