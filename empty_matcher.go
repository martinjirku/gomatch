package gomatch

import (
	"errors"
)

var errNotEmpty = errors.New("expected empty")

// A DateMatcher matches any string
type EmptyMatcher struct {
	pattern string
}

// CanMatch returns true if pattern p can be handled.
func (m *EmptyMatcher) CanMatch(p interface{}) bool {
	return isPattern(p, m.pattern)
}

// Match performs value matching against given pattern.
func (m *EmptyMatcher) Match(p, v interface{}) (bool, error) {
	if v == nil {
		return true, nil
	}
	switch a := v.(type) {
	case string:
		if a == "" {
			return true, nil
		}
	case map[string]interface{}:
		if len(a) == 0 {
			return true, nil
		}
	case []interface{}:
		if len(a) == 0 {
			return true, nil
		}
	}
	return false, errNotEmpty
}

// NewDateMatcher creates StringMatcher.
func NewEmptyMatcher(pattern string) *EmptyMatcher {
	return &EmptyMatcher{pattern}
}
