package floodgate

import (
	"context"
	"net/http"
	"sync"
	"time"
)

const (
	errorText   = "error"
	healthyText = "healthy"
)

type (
	// Function that is used to handle checking. When error is returned, service is considered unhealthy.
	CheckerFunc      = func(ctx context.Context) error
	ResReq           = func(http.ResponseWriter, *http.Request)
	ReportCallerFunc = func(name string, state State)
)

type State interface {
	// Gets the state error. Nil if there is no error.
	Err() error
	// Gets the state status. Returns "healthy" or "error"
	Status() string
	// Gets the message. Error message if there is error. "ok" if there is no error
	Message() string
}

type Checker interface {
	State
	// Error to be set.
	SetError(error)
	// Status to be set.
	SetStatus(string)
	// Message to be set.
	SetMessage(string)
	// Given context is already with timeout set. No need to set it again.
	Check(context.Context) error
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Gate struct {
	services  map[string]Checker
	timeout   time.Duration
	mu        *sync.RWMutex
	baseCtx   context.Context
	interval  time.Duration
	once      *sync.Once
	client    Doer
	reporters []ReportCallerFunc
}

// Creates Gate instance. Not recommended to be used simultaneously with top level Gate instance.
func NewGate(timeout time.Duration, interval time.Duration) *Gate {
	return &Gate{
		services:  make(map[string]Checker),
		timeout:   timeout,
		mu:        &sync.RWMutex{},
		baseCtx:   context.Background(),
		interval:  interval,
		once:      &sync.Once{},
		client:    http.DefaultClient,
		reporters: []ReportCallerFunc{},
	}
}
