package floodgate

import "context"

func (gate *Gate) Check() {
	ctx, release := context.WithTimeout(gate.baseCtx, gate.timeout)
	defer release()
	for name, state := range gate.services {
		gate.wg.Add(1)
		go func(name string, service Checker) {
			defer gate.wg.Done()
			err := service.Check(ctx)
			gate.mu.Lock()
			if err != nil {
				service.SetError(err)
				service.SetStatus(errorText)
				service.SetMessage(err.Error())
			} else {
				service.SetError(nil)
				service.SetStatus(healthyText)
				service.SetMessage(healthyText)
			}
			gate.mu.Unlock()
			s := service.Clone()
			for _, report := range gate.reporters {
				report(name, s)
			}
		}(name, state)
	}

	gate.wg.Wait()
}
