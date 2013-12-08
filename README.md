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
	"os"
	"github.com/plouc/go-marto"
)
```

##Usage

````go
package main

import (
	"os"
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
	
	m.AddReporter(marto.NewWriterReporter(os.Stdout))

	m.Run()
}
````

##Reporting

The package provide a simple reporter accepting an io.Writer, it can be used to send reporting to stdout, file…

Note that you can add several reporters on a single Marto.

###using stdout

````go
m.AddReporter(marto.NewWriterReporter(os.Stdout))
````

###using a file 

````go
fo, err := os.Create("marto.log")
if err != nil { panic(err) }
defer func() {
    if err := fo.Close(); err != nil {
        panic(err)
    }
}()
m.AddReporter(marto.NewWriterReporter(fo))
````