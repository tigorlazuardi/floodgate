package floodgate

func (gate *Gate) RegisterService(name string, f CheckerFunc) *Gate {

	return gate
}

func (gate *Gate) RegisterHandler(name string, c Checker) *Gate {
	gate.mu.Lock()
	gate.services[name] = c
	gate.mu.Unlock()
	return gate
}
