package gomatch

import (
	"errors"
	"time"
)

var ErrNotDate = errors.New("expected date")

// A DateMatcher matches any string
type DateMatcher struct {
	pattern string
}

// CanMatch returns true if pattern p can be handled.
func (m *DateMatcher) CanMatch(p interface{}) bool {
	return isPattern(p, m.pattern)
}

// Match performs value matching against given pattern.
func (m *DateMatcher) Match(p, v interface{}) (bool, error) {
	value, ok := v.(string)
	if !ok {
		return ok, ErrNotDate
	}
	_, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return false, ErrNotDate
	}
	return ok, nil
}

// NewDateMatcher creates StringMatcher.
func NewDateMatcher(pattern string) *DateMatcher {
	return &DateMatcher{pattern}
}
