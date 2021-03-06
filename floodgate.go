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
	// Clone must return a copy of the instance with the pointer does pointing to a new address.
	// Clone will be called once by floodgate and distributed to the given report callers.
	// Report callers may use goroutine or anything with the given clone. It will not affect the original in the floodgate.
	Clone() Checker
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
	wg        *sync.WaitGroup
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
		wg:        &sync.WaitGroup{},
	}
}

func (gate *Gate) AddReportCaller(f ...ReportCallerFunc) *Gate {
	gate.mu.Lock()
	gate.reporters = append(gate.reporters, f...)
	gate.mu.Unlock()
	return gate
}

func (gate *Gate) SetInterval(t time.Duration) *Gate {
	gate.mu.Lock()
	gate.interval = t
	gate.mu.Unlock()
	return gate
}

func (gate *Gate) SetTimeout(t time.Duration) *Gate {
	gate.mu.Lock()
	gate.timeout = t
	gate.mu.Unlock()
	return gate
}
