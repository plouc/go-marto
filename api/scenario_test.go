package api

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"labix.org/v2/mgo"
	"encoding/json"

	"github.com/plouc/go-marto/marto"
)

func getMongo() *mgo.Session {
	mongo, err := mgo.Dial("localhost")
	if err != nil { panic(err) }

	return mongo
}

func clearScenarios(mongo *mgo.Session) {
	_, err := mongo.DB("marto_test").C("scenarios").RemoveAll(nil)
	if err != nil { panic(err) }
}

func loadScenarios(mongo *mgo.Session) {
	scenario := marto.ScenarioConfig{
		Id:       "test0",
		Repeat:   1,
		Every:    0,
		Requests: make([]*marto.RequestTemplate, 0),
	}

	err := mongo.DB("marto_test").C("scenarios").Insert(scenario)
	if err != nil { panic(err) }
}


func TestGetScenariosWithNoScenario(t *testing.T) {
	mongo := getMongo()
	defer mongo.Close()

	clearScenarios(mongo)

	req, err := http.NewRequest("GET", "http://test.dev", nil)
	if err != nil { panic(err) }
	w := httptest.NewRecorder()

	scenarioService := NewScenarioService(mongo, "marto_test")
	scenarioService.GetScenarios(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "[]", w.Body.String())
}


func TestGetScenariosWithScenarios(t *testing.T) {
	mongo := getMongo()
	defer mongo.Close()

	loadScenarios(mongo)

	req, err := http.NewRequest("GET", "http://test.dev", nil)
	if err != nil { panic(err) }
	w := httptest.NewRecorder()

	scenarioService := NewScenarioService(mongo, "marto_test")
	scenarioService.GetScenarios(w, req)

	assert.Equal(t, 200, w.Code)

	var scenarios []*marto.ScenarioConfig
	err = json.Unmarshal(w.Body.Bytes(), &scenarios)
	if err != nil { panic(err) }

	assert.Equal(t, 1, len(scenarios))
	assert.Equal(t, "test0", scenarios[0].Id)
	assert.Equal(t, 1,       scenarios[0].Repeat)
	assert.Equal(t, 0,       scenarios[0].Every)
	assert.Equal(t, 0,       len(scenarios[0].Requests))
}
