package gomatch

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var numberMatcherTests = []struct {
	desc string
	v    interface{}
	ok   bool
	err  error
}{
	{
		// json package uses float64 when unmarshals to interface{}
		"Should match float64",
		100.,
		true,
		nil,
	},
	{
		"Should not match string",
		"100",
		false,
		errNotNumber,
	},
	{
		"Should not match bool",
		true,
		false,
		errNotNumber,
	},
}

func TestNumberMatcher(t *testing.T) {
	pattern := "@pattern@"

	for _, tt := range numberMatcherTests {
		m := NewNumberMatcher(pattern)
		assert.True(t, m.CanMatch(pattern), "expected to support pattern")

		t.Logf(tt.desc)

		ok, err := m.Match(pattern, tt.v)

		if tt.ok {
			assert.True(t, ok)
			assert.Nil(t, err)
		} else {
			assert.False(t, ok)
			assert.True(t, errors.Is(err, tt.err))
		}
	}
}
