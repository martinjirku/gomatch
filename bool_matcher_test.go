package gomatch

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var boolMatcherTests = []struct {
	desc string
	v    interface{}
	ok   bool
	err  error
}{
	{
		"Should match true",
		true,
		true,
		nil,
	},
	{
		"Should match false",
		false,
		true,
		nil,
	},
	{
		"Should not match string",
		"false",
		false,
		errNotBool,
	},
}

func TestBoolMatcher(t *testing.T) {
	pattern := "@pattern@"

	for _, tt := range boolMatcherTests {
		t.Run(tt.desc, func(t *testing.T) {
			m := NewBoolMatcher(pattern)
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
