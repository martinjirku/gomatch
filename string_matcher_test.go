package gomatch

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var stringMatcherTests = []struct {
	desc string
	v    interface{}
	ok   bool
	err  error
}{
	{
		"Should match string",
		"some valid string",
		true,
		nil,
	},
	{
		"Should not match number",
		1234,
		false,
		errNotString,
	},
	{
		"Should not match slice",
		[]interface{}{"a", "b"},
		false,
		errNotString,
	},
}

func TestStringMatcher(t *testing.T) {
	pattern := "@pattern@"

	for _, tt := range stringMatcherTests {
		m := NewStringMatcher(pattern)
		assert.True(t, m.CanMatch(pattern), "expected to support pattern")

		t.Log(tt.desc)

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
