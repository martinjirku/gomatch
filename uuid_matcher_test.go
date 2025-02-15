package gomatch

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var uuidMatcherTests = []struct {
	desc string
	v    interface{}
	ok   bool
	err  error
}{
	{
		"Should match UUID",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		true,
		nil,
	},
	{
		"Should not match invalid UUID",
		"6ba7b810-9dad-XXXX-80b4-00c04fd430c8",
		false,
		ErrNotUUID,
	},
	{
		"Should not match if value is not a string",
		123,
		false,
		ErrNotUUID,
	},
}

func TestUUIDMatcher(t *testing.T) {
	pattern := "@pattern@"

	for _, tt := range uuidMatcherTests {
		t.Run(tt.desc, func(t *testing.T) {
			m := NewUUIDMatcher(pattern)
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
