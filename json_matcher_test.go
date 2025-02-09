package gomatch

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsonMatcherTests = []struct {
	desc string
	p    string
	v    string
	ok   bool
	err  error
}{
	{
		"Should fail if invalid JSON pattern given",
		`{"foo":}`,
		`{"foo": "bar"}`,
		false,
		errInvalidJSONPattern,
	},
	{
		"Should fail if invalid JSON given",
		`{"foo": "bar"}`,
		`{"foo":}`,
		false,
		errInvalidJSON,
	},
	{
		"Should succeed if strings are equal",
		`"John"`,
		`"John"`,
		true,
		nil,
	},
	{
		"Should succeed if numbers are equal",
		`123`,
		`123`,
		true,
		nil,
	},
	{
		"Should succeed if bools are equal",
		`true`,
		`true`,
		true,
		nil,
	},
	{
		"Should fail if types are not equal",
		`"John"`,
		`true`,
		false,
		errTypesNotEqual,
	},
	{
		"Should fail if values are not equal",
		`100`,
		`200`,
		false,
		errValuesNotEqual,
	},
	{
		"Should succeed if objects are equal",
		`
	{
		"id": 123,
		"name": "John Smith"
	}
	`,
		`
	{
		"id": 123,
		"name": "John Smith"
	}
	`,
		true,
		nil,
	},
	{
		"Should succeed if objects are equal but with different key order",
		`
	{
		"id": 123,
		"name": "John Smith"
	}
	`,
		`
	{
		"name": "John Smith",
		"id": 123
	}
	`,
		true,
		nil,
	},
	{
		"Should succeed if arrays are equal",
		"[1,2,3]",
		"[1,2,3]",
		true,
		nil,
	},
	{
		"Should fail if array values in different order",
		"[1,2,3]",
		"[1,3,2]",
		false,
		errValuesNotEqual, //"values are not equal at path: [1]",
	},
	{
		"Should fail if has same keys but values differ",
		`
	{
		"id": 123,
		"name": "John Smith"
	}
	`,
		`
	{
		"id": 999,
		"name": "John Smith"
	}
	`,
		false,
		errValuesNotEqual, //"values are not equal at path: id",
	},
	{
		"Should succeed if nested objects are equal",
		`
	{
		"id": 123,
		"name": "John Smith",
		"phones": [
			{
				"type": "home",
				"phone": "111-111-111"
			},
			{
				"type": "work",
				"phone": "222-222-222"
			}
		]
	}
	`,
		`
	{
		"id": 123,
		"name": "John Smith",
		"phones": [
			{
				"type": "home",
				"phone": "111-111-111"
			},
			{
				"type": "work",
				"phone": "222-222-222"
			}
		]
	}
	`,
		true,
		nil,
	},
	{
		"Should fail if nested objects are not equal",
		`
	{
		"id": 123,
		"name": "John Smith",
		"phones": [
			{
				"type": "home",
				"phone": "111-111-111"
			},
			{
				"type": "work",
				"phone": "222-222-222"
			}
		]
	}
	`,
		`
	{
		"id": 123,
		"name": "John Smith",
		"phones": [
			{
				"type": "home",
				"phone": "111-111-111"
			},
			{
				"type": "mobile",
				"phone": "222-222-222"
			}
		]
	}
	`,
		false,
		errValuesNotEqual, //"values are not equal at path: phones[1].type",
	},
	{
		"Should succeed if values matches patterns",
		`
	{
		"id": "@wildcard@",
		"name": "@wildcard@",
		"phones": [
			{
				"type": "home",
				"phone": "@wildcard@"
			},
			"@wildcard@"
		]
	}
	`,
		`
	{
		"id": 123,
		"name": "John Smith",
		"phones": [
			{
				"type": "home",
				"phone": "111-111-111"
			},
			{
				"type": "office",
				"phone": "222-222-222"
			}
		]
	}
	`,
		true,
		nil,
	},
	{
		"Should fail if object does not have all expected keys",
		`
	{
		"id": 1,
		"name": "John Smith"
	}
	`,
		`
	{
		"id": 1
	}
	`,
		false,
		errMissingKey, //`expected key "name"`,
	},
	{
		"Should fail if object has unexpected keys",
		`
	{
		"id": 1
	}
	`,
		`
	{
		"id": 1,
		"name": "John Smith"
	}
	`,
		false,
		errUnexpectedKey,
	},
	{
		"Should succeed if object has unexpected keys but unbounded pattern was used",
		`
	{
		"id": 1,
		"@...@": ""
	}
	`,
		`
	{
		"id": 1,
		"name": "John Smith"
	}
	`,
		true,
		nil,
	},
	{
		"Should fail if array has unexpected extra values",
		"[1,2,3]",
		"[1,2,3,4]",
		false,
		errArraysLenNotEqual,
	},
	{
		"Should fail if array misses some values",
		"[1,2,3]",
		"[1,2]",
		false,
		errArraysLenNotEqual,
	},
	{
		"Should fail if array misses some values but unbounded pattern was used",
		`[1,2,"@...@"]`,
		"[1]",
		false,
		errArraysLenNotEqual,
	},
	{
		"Should succeed if array has unexpected extra values but unbounded pattern was used",
		`[1,2,3,"@...@"]`,
		"[1,2,3,4]",
		true,
		nil,
	},
	{
		"Should succeed if nested object has extra values but unbounded patterns were used",
		`
	{
		"name": "John Smith",
		"phones": [
			{
				"phone": "111-111-111",
				"@...@": ""
			},
			"@...@"
		],
		"@...@": ""
	}
	`,
		`
	{
		"id": 1,
		"name": "John Smith",
		"phones": [
			{
				"type": "home",
				"phone": "111-111-111"
			},
			{
				"type": "office",
				"phone": "222-222-222"
			}
		]
	}
	`,
		true,
		nil,
	},
	{
		"Should fail if unknown value pattern was used (only @wildcard@ was setup in this suite)",
		`
	{
		"id": "@wildcard@",
		"name": "@string@"
	}
	`,
		`
	{
		"id": 1,
		"name": "John Smith"
	}
	`,
		false,
		errValuesNotEqual, //"values are not equal at path: name",
	},
}

func TestJSONMatcher(t *testing.T) {
	for _, tt := range jsonMatcherTests {
		m := NewJSONMatcher(NewWildcardMatcher(patternWildcard))
		t.Run(tt.desc, func(t *testing.T) {
			ok, err := m.Match(tt.p, tt.v)
			if tt.ok {
				assert.Nil(t, err)
				assert.True(t, ok)
			} else {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.err))
				assert.False(t, ok)
			}
		})
	}
}

func TestJSONMatcherWithDefaultMatchers(t *testing.T) {
	p := `
	{
		"id": "@number@",
		"uuid": "@uuid@",
		"name": "@string@",
		"isActive": "@bool@",
		"createdAt": "@wildcard@",
		"phones": "@array@",
		"email": "@email@",
		"date": "@date@",
		"empty": "@empty@",
		"@...@": ""
	}
	`
	v := `
	{
		"id": 1,
		"uuid": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"name": "John Smith",
		"isActive": true,
		"createdAt": "2018-05-27T12:00:00Z",
		"date": "2018-05-27T12:00:00Z",
		"phones": [
			{
				"type": "home",
				"phone": "111-111-111"
			},
			{
				"type": "office",
				"phone": "222-222-222"
			}
		],
		"email": "john.smith@gmail.com",
		"isVip": false
	}
	`

	m := NewDefaultJSONMatcher()
	ok, err := m.Match(p, v)

	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestJSONMatcherWithErrorReporting(t *testing.T) {
	p := `
	{
		"id": "@number@",
		"uuid": "@uuid@",
		"name": "@string@",
		"isActive": "@bool@",
		"date": "@date@",
		"phones": "@array@",
		"email": "@email@",
		"empty": "@empty@",
		"missingKey": "this_will_miss",
		"deep": {
			"nested": {
				"key": "@string@"
			}
		},
		"arr": [{"key": "@string@"}]
	}
	`
	v := `
	{
		"id": "not_a_number",
		"uuid": "invalid_uuid",
		"name": 666,
		"isActive": "not_a_bool",
		"date": "invalid_2018-05-27T12:00:00Z",
		"phones": {
			"type": "home",
			"phone": "111-111-111"
		},
		"email": "invalid_email",
		"isVip": false,
		"deep": {
			"nested": {
				"key": 123
			}
		},
		"unexpectedKey": null,
		"arr": [{"key": 123}]
	}
	`

	m := NewDefaultJSONMatcher()
	ok, err := m.Match(p, v)
	assert.False(t, ok)
	assert.True(t, errors.Is(err, errNotUUID))
	assert.True(t, errors.Is(err, errNotString))
	assert.True(t, errors.Is(err, errNotBool))
	assert.True(t, errors.Is(err, errNotEmail))
	assert.True(t, errors.Is(err, errNotDate))
	assert.True(t, errors.Is(err, errNotArray))
	assert.True(t, errors.Is(err, errUnexpectedKey))
	assert.True(t, errors.Is(err, errMissingKey))

	assert.True(t, strings.Contains(err.Error(), `expected UUID at "uuid"`))
	assert.True(t, strings.Contains(err.Error(), `expected email at "email"`))
	assert.True(t, strings.Contains(err.Error(), `expected string at "key.nested.deep"`))
	assert.True(t, strings.Contains(err.Error(), `expected number at "id"`))
	assert.True(t, strings.Contains(err.Error(), `expected string at "name"`))
	assert.True(t, strings.Contains(err.Error(), `expected date at "date"`))
	assert.True(t, strings.Contains(err.Error(), `expected array at "phones"`))
	assert.True(t, strings.Contains(err.Error(), `missing key at "missingKey"`))
	assert.True(t, strings.Contains(err.Error(), `types are not equal`))
	assert.True(t, strings.Contains(err.Error(), `unexpected key at ""`))

}
