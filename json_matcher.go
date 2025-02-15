// Package gomatch provides types for pattern based JSON matching.
//
// It provides JSONMatcher type which performs deep comparison of two JSON strings.
// JSONMatcher may be created with a set of ValueMatcher implementations.
// A ValueMatcher is used to make comparison less strict than a regular value comparison.
//
// Use NewDefaultJSONMatcher to create JSONMatcher with a chain of all available ValueMatcher implementations.
//
// Basic usage:
//
//	actual := `
//	{
//		"id": 351,
//		"name": "John Smith",
//		"address": {
//			"city": "Boston"
//		}
//	}
//	`
//	expected := `
//	{
//		"id": "@number@",
//		"name": "John Smith",
//		"address": {
//			"city": "@string@"
//		}
//	}
//	`
//
//	m := gomatch.NewDefaultJSONMatcher()
//	ok, err := m.Match(expected, actual)
//	if ok {
//		fmt.Printf("actual JSON matches expected JSON")
//	} else {
//		fmt.Printf("actual JSON does not match expected JSON: %s", err.Error())
//	}
//
// Use NewJSONMatcher to create JSONMatcher with a custom ValueMatcher implementation.
// Use ChainMatcher to chain multiple ValueMacher implementations.
//
//	m := gomatch.NewJSONMatcher(
//		NewChainMatcher(
//			[]ValueMatcher{
//				NewStringMatcher("@string@"),
//				NewNumberMatcher("@number@"),
//			},
//		)
//	);
package gomatch

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var (
	errInvalidJSON        = errors.New("invalid JSON")
	errInvalidJSONPattern = errors.New("invalid JSON pattern")
	ErrTypesNotEqual      = errors.New("types are not equal")
	errValuesNotEqual     = errors.New("values are not equal")
	errArraysLenNotEqual  = errors.New("arrays sizes are not equal")
	ErrUnexpectedKey      = errors.New("unexpected key")
	ErrMissingKey         = errors.New("missing key")
)

const (
	patternString    = "@string@"
	patternNumber    = "@number@"
	patternBool      = "@bool@"
	patternArray     = "@array@"
	patternUUID      = "@uuid@"
	patternEmail     = "@email@"
	patternWildcard  = "@wildcard@"
	patternDate      = "@date@"
	patternEmpty     = "@empty@"
	patternUnbounded = "@...@"
)

// A ValueMatcher interface should be implemented by any matcher used by JSONMatcher.
type ValueMatcher interface {
	// CanMatch returns true if given pattern can be handled by value matcher implementation.
	CanMatch(p interface{}) bool

	// Match performs the matching of given value v.
	// It also expects pattern p so implementation may handle multiple patterns or some DSL.
	Match(p, v interface{}) (bool, error)
}

// NewDefaultJSONMatcher creates JSONMatcher with default chain of value matchers.
// Default chain contains:
//
// - StringMatcher handling "@string@" pattern
//
// - NumberMatcher handling "@number@" pattern
//
// - BoolMatcher handling "@bool@" pattern
//
// - ArrayMatcher handling "@array@" pattern
//
// - UUIDMatcher handling "@uuid@" pattern
//
// - EmailMatcher handling "@email@" pattern
//
// - DateMatcher handling "@date@" pattern
//
// - EmptyMatcher handling "@empty@" pattern
//
// - WildcardMatcher handling "@wildcard@" pattern
func NewDefaultJSONMatcher() *JSONMatcher {
	return NewJSONMatcher(
		NewChainMatcher(
			[]ValueMatcher{
				NewStringMatcher(patternString),
				NewNumberMatcher(patternNumber),
				NewBoolMatcher(patternBool),
				NewArrayMatcher(patternArray),
				NewUUIDMatcher(patternUUID),
				NewEmailMatcher(patternEmail),
				NewDateMatcher(patternDate),
				NewEmptyMatcher(patternEmpty),
				NewWildcardMatcher(patternWildcard),
			},
		))
}

// NewJSONMatcher creates JSONMatcher with given value matcher.
func NewJSONMatcher(matcher ValueMatcher) *JSONMatcher {
	return &JSONMatcher{matcher}
}

// A JSONMatcher provides Match method to match two JSONs with pattern matching support.
type JSONMatcher struct {
	valueMatcher ValueMatcher
}

// Match performs deep match of given JSON with an expected JSON pattern.
//
// It traverses expected JSON pattern and checks if actual JSON has expected values.
// When traversing it checks if expected value is a pattern supported by internal ValueMatcher.
// In such case it uses the ValueMatcher to match actual value otherwise it compares expected
// value with actual value.
//
// Expected JSON pattern example:
//
//	{
//		"id": "@number@",
//		"name": "John Smith",
//		"address": {
//			"city": "@string@"
//		}
//	}
//
// Matching actual JSON:
//
//	{
//		"id": 351,
//		"name": "John Smith",
//		"address": {
//			"city": "Boston"
//		}
//	}
//
// In above example we assume that ValueMatcher supports "@number@" and "@string@" patterns,
// otherwise matching will fail.
//
// Besides value patterns JSONMatcher supports an "unbounded pattern" - "@...@".
// It can be used at the end of an array to allow any extra array elements:
//
//	[
//		"John Smith",
//		"Joe Doe",
//		"@...@"
//	]
//
// It can be used at the end of an object to allow any extra keys:
//
//	{
//		"id": 351,
//		"name": "John Smith",
//		"@...@": ""
//	}
//
// When matching fails then error message contains a path to invalid value.
func (m *JSONMatcher) Match(expectedJSON, actualJSON string) (bool, error) {
	var expected, actual interface{}
	err := json.Unmarshal([]byte(expectedJSON), &expected)
	if err != nil {
		return false, errInvalidJSONPattern
	}
	err = json.Unmarshal([]byte(actualJSON), &actual)
	if err != nil {
		return false, errInvalidJSON
	}
	err = m.deepMatch(expected, actual, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *JSONMatcher) deepMatch(expected interface{}, actual interface{}, path []interface{}) error {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) && !m.valueMatcher.CanMatch(expected) {
		return NewErrGomatch(ErrTypesNotEqual, path, expected, actual, "")
	}

	switch expected.(type) {
	case []interface{}:
		return m.deepMatchArray(expected.([]interface{}), actual.([]interface{}), path)

	case map[string]interface{}:
		return m.deepMatchMap(expected.(map[string]interface{}), actual.(map[string]interface{}), path)

	default:
		return m.matchValue(expected, actual, path)
	}
}

func (m *JSONMatcher) deepMatchArray(expected, actual, path []interface{}) error {
	unbounded := false
	errs := []error{}
	for i, v := range expected {
		if isUnbounded(v) {
			unbounded = true
			break
		}
		if i == len(actual) {
			break
		}
		errs = append(errs, m.deepMatch(v, actual[i], append(path, i)))
	}
	if !unbounded && len(expected) != len(actual) {
		errs = append(errs, NewErrGomatch(errArraysLenNotEqual, path, expected, actual, ""))
	}
	return errors.Join(errs...)
}

func (m *JSONMatcher) deepMatchMap(expected, actual map[string]interface{}, path []interface{}) error {
	unbounded := false
	errs := []error{}
	for k, v1 := range expected {
		if isUnbounded(k) {
			unbounded = true
			continue
		}
		v2, ok := actual[k]
		if !ok {
			if m.valueMatcher.CanMatch(v1) {
				_, err := m.valueMatcher.Match(v1, nil)
				if err != nil {
					errs = append(errs, NewErrGomatch(err, append(path, k), v1, nil, k))
					continue
				}
				actual[k] = nil
				continue
			}
			errs = append(errs, NewErrGomatch(fmt.Errorf("%w %q", ErrMissingKey, k), path, v1, nil, k))
		} else {
			err := m.deepMatch(v1, v2, append(path, k))
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}
	}
	if !unbounded {
		for k, val := range actual {
			if _, ok := expected[k]; ok {
				continue
			} else {
				errs = append(errs, NewErrGomatch(fmt.Errorf("%w %q", ErrUnexpectedKey, k), path, nil, val, k))
			}
		}
		// errs = append(errs, NewErrGomatch(errUnexpectedKey, path, expected, actual, ""))
	}
	return errors.Join(errs...)
}

func (m *JSONMatcher) matchValue(expected, actual interface{}, path []interface{}) error {
	if m.valueMatcher.CanMatch(expected) {
		_, err := m.valueMatcher.Match(expected, actual)
		return NewErrGomatch(err, path, expected, actual, "")
	}
	if expected != actual {
		return NewErrGomatch(errValuesNotEqual, path, expected, actual, "")
	}
	return nil
}

func isUnbounded(p interface{}) bool {
	return isPattern(p, patternUnbounded)
}

func isPattern(p interface{}, pattern string) bool {
	ps, ok := p.(string)
	return ok && ps == pattern
}
