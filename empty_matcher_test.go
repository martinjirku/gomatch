package gomatch

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var emptyMatcherTests = []struct {
	desc string
	p    string
	v    string
	ok   bool
	err  error
}{
	{
		"Succeed if missing key, when pattern is empty",
		`{"a":1,"b":"@empty@"}`,
		`{"a":1}`,
		true,
		nil,
	},
	{
		"Failed if missing key, when pattern is empty",
		`{"a":1,"b":"sdf"}`,
		`{"a":1}`,
		false,
		errMissingKey,
	},
	{
		"Succeed if empty object, when pattern is empty",
		`{"a":1,"b":"@empty@"}`,
		`{"a":1,"b":{}}`,
		true,
		nil,
	},
	{
		"Succeed if empty array, when pattern is empty",
		`{"a":1,"b":"@empty@"}`,
		`{"a":1,"b":[]}`,
		true,
		nil,
	},
	{
		"Succeed if empty string, when pattern is empty",
		`{"a":1,"b":"@empty@"}`,
		`{"a":1,"b":""}`,
		true,
		nil,
	},
	{
		"Fail if empty string, when pattern is empty",
		`{"a":1,"b":"@empty@"}`,
		`{"a":1,"b": "not empty"}`,
		false,
		errNotEmpty,
	},
}

func TestEmptyMatcher(t *testing.T) {
	pattern := "@empty@"

	for _, tt := range emptyMatcherTests {
		t.Run(tt.desc, func(t *testing.T) {
			m := NewJSONMatcher(NewEmptyMatcher(pattern))
			assert.True(t, m.valueMatcher.CanMatch(pattern), "expected to support pattern")

			ok, err := m.Match(tt.p, tt.v)
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
