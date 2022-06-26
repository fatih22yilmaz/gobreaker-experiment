package main

import (
	"errors"
	"github.com/sony/gobreaker"
	"gobreaker-experiment/infrastructure/http"
	"log"
	netHttp "net/http"
	"time"
)

var downstreamErrCount int
var circuitOpenErrCount int
var timeoutErrCount int

var httpClient = http.NewHTTPClient()

type clientStruct struct {
	circuitBreaker *gobreaker.CircuitBreaker
}

func newClientStruct(circuitBreaker *gobreaker.CircuitBreaker) *clientStruct {
	return &clientStruct{circuitBreaker: circuitBreaker}
}

func main() {
	downstreamErrCount = 0
	circuitOpenErrCount = 0
	timeoutErrCount = 0

	circuitBreaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "GoBreaker",
		MaxRequests: 3,
		Interval:    0,
		Timeout:     1 * time.Second,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Gobreaker's state changed from: %v, to: %v", from, to)
		},
		IsSuccessful: func(err error) bool {
			if err != nil {
				log.Printf("err: %v", err.Error())
			}
			return false
		},
	})

	c := newClientStruct(circuitBreaker)

	netHttp.HandleFunc("/", c.handlerFunc)

	log.Fatal(netHttp.ListenAndServe(":8080", nil))
}

func (c *clientStruct) handlerFunc(w netHttp.ResponseWriter, r *netHttp.Request) {
	response, err := c.circuitBreaker.Execute(func() (interface{}, error) {
		responseByteArray, err := httpClient.Get("http://localhost:8081", "/server-path")
		if err != nil {
			return nil, err
		}

		return responseByteArray, nil
	})

	if err == nil {
		w.Write(response.([]byte))
	}

	// ErrOpenState can be handled as fallback method.
	if err != nil && errors.Is(err, gobreaker.ErrOpenState) {
		log.Printf("error open state: %v", err.Error())
	}

	// ErrTooManyRequests can be handled when too many requests
	if err != nil && errors.Is(err, gobreaker.ErrTooManyRequests) {
		log.Printf("error too many requests: %v", err.Error())
	}
}
