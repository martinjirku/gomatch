package gomatch

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var wildcardMatcherTests = []struct {
	desc string
	v    interface{}
	ok   bool
	err  error
}{
	{
		"Should match everything - string",
		"some string",
		true,
		nil,
	},
	{
		"Should match everything - array",
		[]interface{}{1, 2, 3},
		true,
		nil,
	},
	{
		"Should match everything - number",
		100.,
		true,
		nil,
	},
	{
		"Should match everything - null",
		nil,
		true,
		nil,
	},
	{
		"Should match everything - map",
		map[string]interface{}{"key": "value"},
		true,
		nil,
	},
}

func TestWildcardMatcher(t *testing.T) {
	pattern := "@pattern@"

	for _, tt := range wildcardMatcherTests {
		t.Run(tt.desc, func(t *testing.T) {
			m := NewWildcardMatcher(pattern)
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
