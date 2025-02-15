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
		ErrNotEmail,
	},
	{
		"Should not match without hostname",
		"joe.doe@",
		false,
		ErrNotEmail,
	},
	{
		"Should not match without @",
		"joe.doe[at]gmail.com",
		false,
		ErrNotEmail,
	},
	{
		"Should not match user/box name",
		"@gmail.com",
		false,
		ErrNotEmail,
	},
	{
		"Should not match number",
		1234,
		false,
		ErrNotEmail,
	},
	{
		"Should not match slice",
		[]interface{}{"a", "b"},
		false,
		ErrNotEmail,
	},
}

func TestEmailMatcher(t *testing.T) {
	pattern := "@pattern@"

	for _, tt := range emailMatcherTests {
		t.Run(tt.desc, func(t *testing.T) {
			m := NewEmailMatcher(pattern)
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
