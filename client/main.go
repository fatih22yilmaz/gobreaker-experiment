package main

import (
	"fmt"
	"github.com/sony/gobreaker"
	"gobreaker-experiment/infrastructure/http"
	"log"
	netHttp "net/http"
)

var downstreamErrCount int
var circuitOpenErrCount int
var timeoutErrCount int

var httpClient = http.NewHTTPClient()

type clientStruct struct {
	circuitBreaker *gobreaker.CircuitBreaker
}

func (c *clientStruct) main() {
	downstreamErrCount = 0
	circuitOpenErrCount = 0
	timeoutErrCount = 0

	c.circuitBreaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "GoBreaker",
		MaxRequests: 3,
		Interval:    0,
		Timeout:     5,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Gobreaker's state changed from: %v, to: %v", from, to)
		},
	})

	netHttp.HandleFunc("/", handlerFunc)

	log.Fatal(netHttp.ListenAndServe(":8080", nil))

}

func (c *clientStruct) handlerFunc(w netHttp.ResponseWriter, r *netHttp.Request) {
	var response []byte
	c.circuitBreaker.Execute(func() (interface{}, error) {

	})
	fmt.Printf("\ndownstreamErrCount=%d, circuitOpenErrCount=%d, timeoutErrCount=%d", downstreamErrCount, circuitOpenErrCount, timeoutErrCount)

}
