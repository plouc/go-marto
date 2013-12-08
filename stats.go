package marto

import (
	"time"
)

type RequestStat struct {
	Url        string
	Method     string
	StartedAt  time.Time
	Duration   time.Duration
	StatusCode int
}

type AggregatedRequestStat struct {
	Url             string
	Method          string
	Count           int64
	Total           int64
	AverageDuration int64
	StatusCodes     map[int]uint64
}

func NewRequestStat(method string, url string, startedAt time.Time) *RequestStat {
	return &RequestStat{url, method, startedAt, 0, 0}
}

func (rs *RequestStat) Finished() {
	rs.Duration = time.Since(rs.StartedAt)
}
