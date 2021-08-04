package floodgate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewGate(t *testing.T) {
	wantTimeout, wantInterval := time.Second*10, time.Second*10
	gate := NewGate(wantTimeout, wantInterval)

	assert.Equal(t, wantInterval, gate.interval)
	assert.Equal(t, wantTimeout, gate.timeout)
	assert.NotNil(t, gate.services)
	assert.NotNil(t, gate.baseCtx)
	assert.NotNil(t, gate.mu)
	assert.NotNil(t, gate.once)
	assert.NotNil(t, gate.reporters)
}
