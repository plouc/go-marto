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
}

func (rs *RequestStat) Finished() {
	rs.Duration = time.Since(rs.StartedAt)
}
