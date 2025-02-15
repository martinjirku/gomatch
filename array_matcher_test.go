package gomatch

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var arrayMatcherTests = []struct {
	desc string
	v    interface{}
	ok   bool
	err  error
}{
	{
		"Should match slice",
		[]interface{}{1, 2, 3},
		true,
		nil,
	},
	{
		"Should match empty slice",
		[]interface{}{},
		true,
		nil,
	},
	{
		"Should not match string",
		"some string",
		false,
		ErrNotArray,
	},
	{
		"Should not match nil",
		nil,
		false,
		ErrNotArray,
	},
}

func TestArrayMatcher(t *testing.T) {
	pattern := "@pattern@"

	for _, tt := range arrayMatcherTests {
		t.Run(tt.desc, func(t *testing.T) {
			m := NewArrayMatcher(pattern)
			assert.True(t, m.CanMatch(pattern), "expected to support pattern")

			ok, err := m.Match(pattern, tt.v)
			if tt.ok {
				assert.True(t, ok)
				assert.Nil(t, err)
			} else {
				assert.False(t, ok)
				assert.True(t, errors.Is(err, tt.err))
			}
		})
	}
}
