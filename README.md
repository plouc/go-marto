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

	s := marto.NewScenario("search")
	s.Append("GET", "http://google.com", nil)
	req := s.Append("GET", "http://google.com/search?q=test", nil)
	req.SetDelay(2000)
	s.Repeat(2)
	m.AddScenario(s)
	
	m.AddReporter(&marto.SimpleReporter{})

	m.Run()
}
````
