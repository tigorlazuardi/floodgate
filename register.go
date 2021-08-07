package floodgate

import (
	"context"
	"errors"
	"net/http"
)

func (gate *Gate) RegisterService(name string, f CheckerFunc) *Gate {
	safe := func(ctx context.Context) (err error) {
		defer func() {
			if rec := recover(); rec != nil {
				if e, ok := rec.(error); ok {
					err = newGateError(name, e)
					return
				}
				if e, ok := rec.(string); ok {
					err = newGateError(name, errors.New(e))
					return
				}
				err = newGateError(name, errors.New("panic caught on health checking service: "+name))
			}
		}()
		err = f(ctx)
		return
	}
	return gate.RegisterHandler(name, &service{check: safe})
}

func (gate *Gate) RegisterHandler(name string, c Checker) *Gate {
	gate.mu.Lock()
	gate.services[name] = c
	gate.mu.Unlock()
	return gate
}

/*
RegisterHTTPService checks an endpoint and see if the returned endpoint returned status code under codeTreshold.

If there's an error on attempting connection, it will be marked as unhealthy.

If status code is equal or higher than codeTreshold, it will be marked as unhealthy.

if client is nil, it will use http.DefaultClient.

Request is always using method GET and no request body is sent.
*/
func (gate *Gate) RegisterHTTPService(name string, url string, codeTreshold int, client Doer) *Gate {
	if client == nil {
		client = http.DefaultClient
	}
	gate.mu.Lock()
	gate.services[name] = &httpservice{
		client:       client,
		CodeTreshold: codeTreshold,
		Url:          url,
	}
	gate.mu.Unlock()
	return gate
}


