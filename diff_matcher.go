package gomatch

import (
	"encoding/json"
	"fmt"
)

// DiffResult represents the result of a JSON diff operation
type DiffResult struct {
	Path     string      // The path where the difference was found
	Expected interface{} // The expected value
	Actual   interface{} // The actual value
	Type     DiffType    // The type of difference
}

// DiffType represents the type of difference found
type DiffType int

const (
	DiffTypeValueMismatch DiffType = iota
	DiffTypeTypeMismatch
	DiffTypeMissingKey
	DiffTypeUnexpectedKey
	DiffTypeArrayLengthMismatch
)

// JSONDiff provides functionality to find differences between two JSON structures
type JSONDiff struct {
	matcher ValueMatcher
}

// Diff compares two JSON strings and returns a slice of differences.
// If either the expected or actual JSON is invalid, an error is returned.
func (d *JSONDiff) Diff(expected, actual string) ([]DiffResult, error) {
	var expectedObj, actualObj interface{}

	if err := json.Unmarshal([]byte(expected), &expectedObj); err != nil {
		return nil, fmt.Errorf("invalid expected JSON: %w", err)
	}

	if err := json.Unmarshal([]byte(actual), &actualObj); err != nil {
		return nil, fmt.Errorf("invalid actual JSON: %w", err)
	}

	differences := make([]DiffResult, 0)
	d.diffValues(expectedObj, actualObj, "", &differences)

	return differences, nil
}

// diffValues recursively compares the expected and actual values, appending any differences to the diffs slice.
func (d *JSONDiff) diffValues(expected, actual interface{}, path string, diffs *[]DiffResult) {
	switch exp := expected.(type) {
	case map[string]interface{}:
		if act, ok := actual.(map[string]interface{}); ok {
			d.diffObjects(exp, act, path, diffs)
		} else {
			*diffs = append(*diffs, DiffResult{
				Path:     path,
				Expected: expected,
				Actual:   actual,
				Type:     DiffTypeTypeMismatch,
			})
		}
	case []interface{}:
		if act, ok := actual.([]interface{}); ok {
			d.diffArrays(exp, act, path, diffs)
		} else {
			*diffs = append(*diffs, DiffResult{
				Path:     path,
				Expected: expected,
				Actual:   actual,
				Type:     DiffTypeTypeMismatch,
			})
		}
	default:
		if !d.valuesEqual(expected, actual) {
			*diffs = append(*diffs, DiffResult{
				Path:     path,
				Expected: expected,
				Actual:   actual,
				Type:     DiffTypeValueMismatch,
			})
		}
	}
}

// diffObjects compares the expected and actual objects, appending any differences to the diffs slice.
func (d *JSONDiff) diffObjects(expected, actual map[string]interface{}, path string, diffs *[]DiffResult) {
	for key, expectedVal := range expected {
		actualVal, exists := actual[key]
		currentPath := d.joinPath(path, key)

		if !exists {
			*diffs = append(*diffs, DiffResult{
				Path:     currentPath,
				Expected: expectedVal,
				Actual:   nil,
				Type:     DiffTypeMissingKey,
			})
			continue
		}

		d.diffValues(expectedVal, actualVal, currentPath, diffs)
	}

	for key, actualVal := range actual {
		if _, exists := expected[key]; !exists {
			*diffs = append(*diffs, DiffResult{
				Path:     d.joinPath(path, key),
				Expected: nil,
				Actual:   actualVal,
				Type:     DiffTypeUnexpectedKey,
			})
		}
	}
}

// diffArrays compares the expected and actual arrays, appending any differences to the diffs slice.
func (d *JSONDiff) diffArrays(expected, actual []interface{}, path string, diffs *[]DiffResult) {
	if len(expected) != len(actual) {
		*diffs = append(*diffs, DiffResult{
			Path:     path,
			Expected: len(expected),
			Actual:   len(actual),
			Type:     DiffTypeArrayLengthMismatch,
		})
	}

	minLen := len(expected)
	if len(actual) < minLen {
		minLen = len(actual)
	}

	for i := 0; i < minLen; i++ {
		currentPath := fmt.Sprintf("%s[%d]", path, i)
		d.diffValues(expected[i], actual[i], currentPath, diffs)
	}
}

// valuesEqual compares the expected and actual values using the provided ValueMatcher.
func (d *JSONDiff) valuesEqual(expected, actual interface{}) bool {
	if d.matcher.CanMatch(expected) {
		ok, _ := d.matcher.Match(expected, actual)
		return ok
	}
	return expected == actual
}

// joinPath joins the base path and the key to create a full path.
func (d *JSONDiff) joinPath(base, key string) string {
	if base == "" {
		return key
	}
	return base + "." + key
}
