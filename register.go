package floodgate

import (
	"context"
	"errors"
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
