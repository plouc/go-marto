package marto

import (
	"net/http"
	"io"
	"time"
	//"fmt"
)

type RequestTemplate struct {
	Method string
	Url    string
	Body   io.Reader
	delay  uint64
}

func NewRequestTemplate(method string, url string, body io.Reader) *RequestTemplate {
	return &RequestTemplate{
		Method: method,
		Url:    url,
		Body:   body,
	}
}

func (rt *RequestTemplate) SetDelay(delay uint64) *RequestTemplate {
	rt.delay = delay

	return rt
}

func (rt *RequestTemplate) HasDelay() bool {
	return rt.delay > 0
}

func (rt *RequestTemplate) Delay() uint64 {
	return rt.delay
}



type Request struct {
	*http.Request

	Scenario *Scenario
	Session  *Session

	id        int

	delay     uint64

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
		Session:   new(Session),
		id:        id,
		delay:     0,
		feeders:   make([]Feeder, 0),
		StartedAt: time.Time{},
		EndedAt:   time.Time{},
		Elapsed:   time.Duration(0),
	}, nil
}

func (r *Request) IsFirst() bool {
	return r.id == 0
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

func (r *Request) BindFeeder(f Feeder) {
	r.feeders = append(r.feeders, f)
}

func (r *Request) Resolve() {
	//fmt.Printf("%#v\n", r.URL)
}