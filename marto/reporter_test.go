package marto

import (
	"testing"
	"bytes"

	"github.com/stretchr/testify/assert"
)


func TestHandlers(t *testing.T) {

	reporter := NewBaseReporter()
	scenario := NewScenario("test_00")

	assert.Equal(t, 0, len(reporter.ScenariosStats))
	assert.Equal(t, 0, len(reporter.Writers))
	
	writer := bytes.NewBuffer(make([]byte, 0))
	reporter.AddWriter(writer)

	assert.Equal(t, 1, len(reporter.Writers))

	reporter.OnScenarioStarted(scenario)

	assert.Equal(t, 1, len(reporter.ScenariosStats))
	//assert.Equal(t, 0, len(reporter.ScenariosStats["test_0"].Histogram))
	//assert.Equal(t, 0, len(reporter.ScenariosStats["test_0"].Stats))

	assert.Contains(t, writer.String(), "scenario.started | test_00")
}
