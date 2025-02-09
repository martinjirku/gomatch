package gomatch

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var emailMatcherTests = []struct {
	desc string
	v    interface{}
	ok   bool
	err  error
}{
	{
		"Should match email",
		"joe.doe@gmail.com",
		true,
		nil,
	},
	{
		"Should match email with IP",
		"joe.doe@192.168.1.5",
		true,
		nil,
	},
	{
		"Should match email with hostname without dot",
		"joe.doe@somehostname",
		true,
		nil,
	},
	{
		"Should not match email with underscore",
		"joe.doe@my_mail.com",
		false,
		errNotEmail,
	},
	{
		"Should not match without hostname",
		"joe.doe@",
		false,
		errNotEmail,
	},
	{
		"Should not match without @",
		"joe.doe[at]gmail.com",
		false,
		errNotEmail,
	},
	{
		"Should not match user/box name",
		"@gmail.com",
		false,
		errNotEmail,
	},
	{
		"Should not match number",
		1234,
		false,
		errNotEmail,
	},
	{
		"Should not match slice",
		[]interface{}{"a", "b"},
		false,
		errNotEmail,
	},
}

func TestEmailMatcher(t *testing.T) {
	pattern := "@pattern@"

	for _, tt := range emailMatcherTests {
		m := NewEmailMatcher(pattern)
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
