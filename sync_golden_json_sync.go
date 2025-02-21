package gomatch

import (
	"encoding/json"
	"reflect"
)

// JSONMarshalFn is a function type that defines JSON marshaling behavior.
// It takes any value and returns the JSON byte representation and an error.
// This allows for custom JSON marshaling implementations to be used.
type JSONMarshalFn func(v any) ([]byte, error)

// GoldenJSONSync provides functionality to synchronize and match JSON content
// against golden (expected) patterns. It supports pattern matching for various
// data types and structures while maintaining the original JSON structure.
type GoldenJSONSync struct {
	valueMatcher ValueMatcher
	marshaler    JSONMarshalFn
}

// NewGoldenJSONSync creates a new GoldenJSONSync instance with default pattern matchers.
// It initializes the sync with a chain of predefined matchers for common data types:
//   - String patterns (using patternString)
//   - Number patterns (using patternNumber)
//   - Boolean patterns (using patternBool)
//   - Array patterns (using patternArray)
//   - UUID patterns (using patternUUID)
//   - Email patterns (using patternEmail)
//   - Date patterns (using patternDate)
//   - Empty patterns (using patternEmpty)
//   - Wildcard patterns (using patternWildcard)
//
// Returns a pointer to the configured GoldenJSONSync instance.
func NewGoldenJSONSync() *GoldenJSONSync {
	return NewGoldenJSON(
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

// NewGoldenJSON creates a new GoldenJSONSync instance with a custom matcher.
// This constructor allows for more flexibility in defining matching behavior
// compared to NewGoldenJSONSync().
//
// Parameters:
//   - matcher: A ValueMatcher implementation that defines how values should be matched
//
// Returns a pointer to a GoldenJSONSync instance configured with the provided matcher
// and the default json.Marshal function.
func NewGoldenJSON(matcher ValueMatcher) *GoldenJSONSync {
	return &GoldenJSONSync{matcher, json.Marshal}
}

func (g *GoldenJSONSync) Marshaler(m JSONMarshalFn) {
	g.marshaler = m
}

// Sync synchronizes a golden (expected) JSON with a new JSON string while preserving
// pattern matching expressions from the golden JSON.
//
// The function performs the following steps:
//  1. Unmarshals both golden and new JSON strings into interface{} types
//  2. Performs deep matching between the golden and new JSON structures
//  3. Preserves pattern matching expressions from the golden JSON where applicable
//  4. Incorporates new values from the new JSON where patterns don't exist
//  5. Marshals the resulting structure back to a JSON string
//
// Parameters:
//   - goldenJSON: The expected JSON string containing pattern matching expressions
//   - newJSON: The new JSON string to sync with the golden JSON
//
// Returns:
//   - string: The synchronized JSON string
//   - error: An error if any step fails (invalid JSON, marshaling errors)
//
// If the golden JSON is invalid, returns the new JSON as-is
// If the new JSON is invalid, returns the golden JSON and an error
func (g *GoldenJSONSync) Sync(goldenJSON, newJSON string) (string, error) {
	var golden, actual interface{}
	err := json.Unmarshal([]byte(goldenJSON), &golden)
	if err != nil {
		return newJSON, nil
	}
	err = json.Unmarshal([]byte(newJSON), &actual)
	if err != nil {
		return goldenJSON, errInvalidJSON
	}
	newGolden := g.deepMatch(golden, actual)
	newGoldenJSON, err := g.marshaler(newGolden)
	if err != nil {
		return goldenJSON, err
	}
	return string(newGoldenJSON), nil
}

func (g *GoldenJSONSync) deepMatch(golden interface{}, actual interface{}) interface{} {
	if reflect.TypeOf(golden) != reflect.TypeOf(actual) && !g.valueMatcher.CanMatch(golden) {
		return actual
	}

	switch golden.(type) {
	case []interface{}:
		return g.deepMatchArray(golden.([]interface{}), actual.([]interface{}))

	case map[string]interface{}:
		return g.deepMatchMap(golden.(map[string]interface{}), actual.(map[string]interface{}))

	default:
		return g.matchValue(golden, actual)
	}
}

func (g *GoldenJSONSync) deepMatchArray(golden, actual []interface{}) []interface{} {
	unbounded := false
	results := []interface{}{}
	for i, goldenVal := range golden {
		if isUnbounded(goldenVal) {
			unbounded = true
			results = append(results, goldenVal)
			break
		}
		if i >= len(actual) {
			break
		}
		results = append(results, g.deepMatch(goldenVal, actual[i]))
	}
	if !unbounded && len(golden) < len(actual) {
		results = append(results, actual[len(golden):]...)
	}
	return results
}

func (g *GoldenJSONSync) deepMatchMap(golden, actual map[string]interface{}) map[string]interface{} {
	unbounded := false
	results := map[string]interface{}{}
	for k, goldenVal := range golden {
		if isUnbounded(k) {
			unbounded = true
			results[k] = nil
			continue
		}
		if actualVal, ok := actual[k]; ok {
			results[k] = g.deepMatch(goldenVal, actualVal)
		}
	}
	if !unbounded {
		for k, v2 := range actual {
			if _, ok := golden[k]; !ok {
				results[k] = v2
			}
		}
	}
	return results
}

func (g *GoldenJSONSync) matchValue(golden, actual interface{}) interface{} {
	if g.valueMatcher.CanMatch(golden) {
		return golden
	}
	return actual
}
