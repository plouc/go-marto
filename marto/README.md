go-marto
========

go-marto is an http stress tool written in golang.

[![Build Status](https://travis-ci.org/plouc/go-marto.png?branch=master)](https://travis-ci.org/plouc/go-marto)

##Installation

To install go-marto, use `go get`:

    go get github.com/plouc/go-marto

Import the `go-marto` package into your code:

```go
import "github.com/plouc/go-marto/marto"
```

##Basic usage

````go
package main

import (
    "os"
    _marto "github.com/plouc/go-marto/marto"
)

func main() {
    marto := _marto.NewMarto()

    scenario := _marto.NewScenario("search").Repeat(2)
    scenario.Append("GET", "http://google.com", nil)
    scenario.Append("GET", "http://google.com/search?q=test", nil).SetDelay(2000)
    marto.AddScenario(scenario)

    reporter := _marto.NewBaseReporter()
    reporter.AddWriter(os.Stdout)
    marto.AddReporter(reporter)

    marto.Run()
}
````

##Building scenarios

A **Scenario** contains a suite of **Requests**, it is used to create a **Session**, finally **Requests** from the **Scenario** aren't directly used to make http calls, instead, thoses **Requests** are copied into a new **Session**. You can think of the **Scenario** as a Template and the **Session** as an instance of this template.

The decoupling is usefull to distribute requests over time.

###Creating a scenario

````go
scenario := marto.NewScenario("search")
scenario.Append("GET", "http://google.com", nil)
````

We've just created a basic Scenario composed of one GET request, but wait, if I want to make thousand of requests simultaneously, I have to create each Scenario manually ? No, you just have to make your scenario repeatable:

````go
scenario := marto.NewScenario("search").Repeat(100)
scenario.Append("GET", "http://google.com", nil)
````

You will now spawn 100 requests simultaneously, but having those all spawned at the same time can be weird, you should distribute them among time, to make the tool behave in a more natural way, launching the first session, waiting 100ms, launching the second session and so on…:

````go
scenario := marto.NewScenario("search").Repeat(100).Every(100)
scenario.Append("GET", "http://google.com", nil)
````

Now you will spawn a request every 100ms until the desired count is reached.

Consider this other example:

````go
scenario := marto.NewScenario("search").Repeat(100).Every(100)
scenario.Append("GET", "http://google.com", nil)
scenario.Append("GET", "http://google.com/search?q=test", nil)
````

The first session will start at 0ms, the second one at 100ms, the third at 200ms… but the two requests of the scenario will again spawn simultaneously, to change this you can alter the **Request** sent back when you call "s.Append()" to make it wait for a given time:

````go
scenario := marto.NewScenario("search").Repeat(100).Every(100)
scenario.Append("GET", "http://google.com", nil)
scenario.Append("GET", "http://google.com/search?q=test", nil).SetDelay(2000)
````

Now the first session will start at 0ms, the first request of this session will also start at 0ms, the second request will start at 2000ms, the second session will start at 100ms, the first request of the second session will start at 100ms and the second one at 2100ms (the value passed to Every() * the position of the session + the delay set on the request)…

##Reporting

###WriterReporter

The package provide a simple reporter accepting an **io.Writer**, it can be used to send reporting to stdout, file…

Note that you can add several reporters on a single Marto.

####WriterReporter using stdout

````go
m.AddReporter(marto.NewWriterReporter(os.Stdout))
````

####WriterReporter using a file 

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

###AggregatorReporter

The **AggregatorReporter** automatically store the number of request iterations and compute the average duration plus an histogram of requests per second.

````go
m.AddReporter(marto.NewAggregatorReporter())
reporter.Dump(os.Stdout)
reporter.DumpJson(os.Stdout)
````

###Customize

You can easily add **custom reporters**, you just have to conform to the **Reporter interface**:

````go
type Reporter interface {
	OnScenarioStarted(scenario *Scenario)
	OnScenarioFinished(scenario *Scenario)

	OnSessionStarted(session *Session)
	OnSessionFinished(session *Session)

	OnRequest(request *Request)
	OnResponse(request *Request, response *http.Response)
}
````

Then, you should build a custom alerting system which send an e-mail when slow requests are encountred.
Add your logic in the **OnResponse** function and you're done.
