go-marto
========

go-marto is an http stress tool written in golang.

##Installation

To install go-marto, use `go get`:

    go get github.com/plouc/go-marto

Import the `go-marto` package into your code:

```go
package whatever

import (
	"github.com/plouc/go-marto"
)
```

##Usage

````go
package main

import (
	"fmt"
	"github.com/plouc/go-marto"
)

func main() {
	m := marto.NewMarto()

	s := marto.NewScenario()
	s.AppendRequestFromConfig("GET", "http://rbenitte.com", nil, 0)
	s.AppendRequestFromConfig("GET", "http://rbenitte.com/14/skull", nil, 1000)
	s.AppendRequestFromConfig("GET", "http://rbenitte.com/13/studio-ilelle", nil, 1000)
	s.Repeat(2)
	m.AddScenario("rbenitte", s)

	m.Start("rbenitte")
	m.AggregateRequestStats()

	for _, aggReqStat := range m.AggregatedRequestStats {
		fmt.Printf("%s - %d request(s) - average: %dms (total %dms)\n", aggReqStat.Url, aggReqStat.Count, aggReqStat.AverageDuration / 1000000, aggReqStat.Total / 1000000)

		for statusCode, count := range aggReqStat.StatusCodes {
			fmt.Printf("  %d -> %d request(s)\n", statusCode, count)
		}
	}
}
````
