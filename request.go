package marto

import (
	"net/http"
	"io"
	"time"
	//"fmt"
)

type Request struct {
	*http.Request

	Scenario *Scenario

	id        int

	delay     uint64
	duration  uint64

	feeders   []Feeder

	StartedAt time.Time
	EndedAt   time.Time
	Elapsed   time.Duration
}

func NewRequest(id int, scenario *Scenario, method string, strUrl string, body io.Reader) (*Request, error) {

	req, err := http.NewRequest(method, strUrl, body)
	if err != nil {
		return nil, err
	}

	return &Request{
		Request:   req,
		Scenario:  scenario,
		id:        id,
		delay:     0,
		duration:  0,
		feeders:   make([]Feeder, 0),
		StartedAt: time.Time{},
		EndedAt:   time.Time{},
		Elapsed:   time.Duration(0),
	}, nil
}

func (r *Request) Start() {
	r.StartedAt = time.Now()
}

func (r *Request) End() {
	r.EndedAt = time.Now()
	r.Elapsed = time.Since(r.StartedAt)
}

func (r *Request) Id() int {
	return r.id
}

func (r *Request) HasDelay() bool {
	return r.delay > 0
}

func (r *Request) SetDelay(delay uint64) *Request {
	r.delay = delay

	return r
}

func (r *Request) Delay() uint64 {
	return r.delay
}

func (r *Request) SetIn(duration uint64) *Request {
	r.duration = duration

	return r
}

func (r *Request) In() uint64 {
	return r.duration
}

func (r *Request) BindFeeder(f Feeder) {
	r.feeders = append(r.feeders, f)
}

func (r *Request) Resolve() {
	//fmt.Printf("%#v\n", r.URL)
}