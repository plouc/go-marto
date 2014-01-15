package api

import (
	"net/http"
	"fmt"

	"github.com/gorilla/mux"
	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"github.com/plouc/go-marto/marto"
)

type ScenarioService struct {
	mongo      *mgo.Session
	collection *mgo.Collection
}

func NewScenarioService(mgoSession *mgo.Session) *ScenarioService {

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	collection := mgoSession.DB("marto").C("scenarios")
	err := collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return &ScenarioService{
		mongo:      mgoSession,
		collection: collection,
	}
}


func (s *ScenarioService) GetScenarios(w http.ResponseWriter, r *http.Request) {

	results := []marto.ScenarioConfig{}
	err := s.collection.Find(nil).All(&results)
	if err != nil {
		panic(err)
	}

	js, _ := json.Marshal(results)
	fmt.Fprint(w, string(js))
}


func (s *ScenarioService) GetScenario(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id   := vars["id"]

	result := marto.ScenarioConfig{}
	err := s.collection.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		panic(err)
	}

	js, _ := json.Marshal(result)
	fmt.Fprint(w, string(js))
}


func (s *ScenarioService) SaveScenario(w http.ResponseWriter, r *http.Request) {

	fmt.Println("saveScenario")

	decoder := json.NewDecoder(r.Body)
	var scenario marto.ScenarioConfig
	err := decoder.Decode(&scenario)
	if err != nil {
		panic(err)
	}

	err = s.collection.Insert(scenario)
	if err != nil {
		panic(err)
	}

	js, _ := json.Marshal(scenario)
	fmt.Fprint(w, string(js))
}
