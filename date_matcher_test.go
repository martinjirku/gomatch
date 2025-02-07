package gomatch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var dateMatcherTests = []struct {
	desc   string
	v      interface{}
	ok     bool
	errMsg string
}{
	{
		"Default date format",
		"some valid date",
		true,
		"",
	},
	{
		"Should not match date",
		"2020-01-01T12:34:56Z",
		false,
		"expected date",
	},
	{
		"Should not match slice",
		[]interface{}{"a", "b"},
		false,
		"expected date",
	},
}

func TestDateMatcher(t *testing.T) {
	pattern := "@pattern@"

	for _, tt := range stringMatcherTests {
		m := NewDateMatcher(pattern)
		assert.True(t, m.CanMatch(pattern), "expected to support pattern")

		t.Logf(tt.desc)

		ok, err := m.Match(pattern, tt.v)

		if tt.ok {
			assert.True(t, ok)
			assert.Nil(t, err)
		} else {
			assert.False(t, ok)
			assert.EqualError(t, err, tt.errMsg)
		}
	}
}
