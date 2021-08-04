package floodgate

import (
	"encoding/json"
	"fmt"
)

type GateError struct {
	Name      string
	OriginErr error
	Err       string
}

func newGateError(name string, err error) GateError {
	return GateError{
		Name:      name,
		OriginErr: err,
		Err:       err.Error(),
	}
}

func (ge GateError) Error() string {
	return ge.Err
}

func (ge GateError) String() string {
	return fmt.Sprintf("service '%s' returned error: '%s'", ge.Name, ge.Err)
}

func (ge GateError) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["name"] = ge.Name
	if val, _ := json.Marshal(ge.OriginErr); string(val) == "{}" || string(val) == "" {
		m["cause"] = ge.Err
	} else {
		m["cause"] = string(val)
	}
	m["error"] = ge.Err
	return json.Marshal(m)
}
