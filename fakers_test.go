package stubs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerators_SimpleGenerators(t *testing.T) {
	g, err := newGenerator("")
	assert.NoError(t, err)
	for i := 0; i < 20; i++ {
		b, err := g.bool(nil)
		assert.NoError(t, err)
		t.Logf("%d - bool(): %v", i, b)

		n64, err := g.numGenInt64(nil)
		assert.NoError(t, err)
		t.Logf("%d - numGenInt64(): %v", i, n64)

		f64, err := g.numGenFloat64(nil)
		assert.NoError(t, err)
		t.Logf("%d - numGenFloat64(): %v", i, f64)

		f32, err := g.numGenFloat32(nil)
		assert.NoError(t, err)
		t.Logf("%d - numGenFloat32(): %v", i, f32)

		d, err := g.dateGen(nil)
		assert.NoError(t, err)
		t.Logf("%d - dateGen(): %v", i, d)

		dt, err := g.dateTimeGen(nil)
		assert.NoError(t, err)
		t.Logf("%d - dateTimeGen(): %v", i, dt)

		du, err := g.durationGen(nil)
		assert.NoError(t, err)
		t.Logf("%d - durationGen(): %v", i, du)

	}
}
