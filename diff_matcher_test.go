package gomatch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONDiff(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		actual   string
		wantDiff []DiffResult
		wantErr  bool
	}{
		{
			name:     "same json",
			expected: `{"name": "John"}`,
			actual:   `{"name": "John"}`,
			wantDiff: []DiffResult{},
		},
		{
			name:     "simple value difference",
			expected: `{"name": "John"}`,
			actual:   `{"name": "Jane"}`,
			wantDiff: []DiffResult{
				{
					Path:     "name",
					Expected: "John",
					Actual:   "Jane",
					Type:     DiffTypeValueMismatch,
				},
			},
		},
		{
			name:     "missing key",
			expected: `{"name": "John", "age": 30}`,
			actual:   `{"name": "John"}`,
			wantDiff: []DiffResult{
				{
					Path:     "age",
					Expected: float64(30),
					Actual:   nil,
					Type:     DiffTypeMissingKey,
				},
			},
		},
		{
			name:     "pattern matching",
			expected: `{"id": "@number@", "email": "@email@"}`,
			actual:   `{"id": "not-a-number", "email": "invalid-email"}`,
			wantDiff: []DiffResult{
				{
					Path:     "id",
					Expected: "@number@",
					Actual:   "not-a-number",
					Type:     DiffTypeValueMismatch,
				},
				{
					Path:     "email",
					Expected: "@email@",
					Actual:   "invalid-email",
					Type:     DiffTypeValueMismatch,
				},
			},
		},
		{
			name:     "nested object difference",
			expected: `{"user": {"name": "John", "details": {"age": 30}}}`,
			actual:   `{"user": {"name": "John", "details": {"age": 25}}}`,
			wantDiff: []DiffResult{
				{
					Path:     "user.details.age",
					Expected: float64(30),
					Actual:   float64(25),
					Type:     DiffTypeValueMismatch,
				},
			},
		},
		{
			name:     "array length mismatch",
			expected: `{"items": [1, 2, 3]}`,
			actual:   `{"items": [1, 2]}`,
			wantDiff: []DiffResult{
				{
					Path:     "items",
					Expected: 3,
					Actual:   2,
					Type:     DiffTypeArrayLengthMismatch,
				},
			},
		},
		{
			name:     "array element mismatch",
			expected: `{"scores": [85, 90, 95]}`,
			actual:   `{"scores": [85, 92, 95]}`,
			wantDiff: []DiffResult{
				{
					Path:     "scores[1]",
					Expected: float64(90),
					Actual:   float64(92),
					Type:     DiffTypeValueMismatch,
				},
			},
		},
		{
			name:     "multiple pattern matchers",
			expected: `{"timestamp": "@date@", "temperature": "@number@", "uuid": "@uuid@"}`,
			actual:   `{"timestamp": "invalid-date", "temperature": "25.5", "uuid": "invalid-uuid"}`,
			wantDiff: []DiffResult{
				{
					Path:     "timestamp",
					Expected: "@date@",
					Actual:   "invalid-date",
					Type:     DiffTypeValueMismatch,
				},
				{
					Path:     "temperature",
					Expected: "@number@",
					Actual:   "25.5",
					Type:     DiffTypeValueMismatch,
				},
				{
					Path:     "uuid",
					Expected: "@uuid@",
					Actual:   "invalid-uuid",
					Type:     DiffTypeValueMismatch,
				},
			},
		},
		{
			name:     "extra key in actual",
			expected: `{"name": "John"}`,
			actual:   `{"name": "John", "extra": "field"}`,
			wantDiff: []DiffResult{
				{
					Path:     "extra",
					Expected: nil,
					Actual:   "field",
					Type:     DiffTypeUnexpectedKey,
				},
			},
		},
		{
			name:     "null value handling",
			expected: `{"optional": null, "required": "value"}`,
			actual:   `{"optional": "not-null", "required": null}`,
			wantDiff: []DiffResult{
				{
					Path:     "optional",
					Expected: nil,
					Actual:   "not-null",
					Type:     DiffTypeValueMismatch,
				},
				{
					Path:     "required",
					Expected: "value",
					Actual:   nil,
					Type:     DiffTypeValueMismatch,
				},
			},
		},
		{
			name:     "type mismatch - object vs array",
			expected: `{"data": {"key": "value"}}`,
			actual:   `{"data": ["value"]}`,
			wantDiff: []DiffResult{
				{
					Path:     "data",
					Expected: map[string]interface{}{"key": "value"},
					Actual:   []interface{}{"value"},
					Type:     DiffTypeTypeMismatch,
				},
			},
		},
		{
			name:     "type mismatch - array vs object",
			expected: `{"data": ["first", "second"]}`,
			actual:   `{"data": {"0": "first", "1": "second"}}`,
			wantDiff: []DiffResult{
				{
					Path:     "data",
					Expected: []interface{}{"first", "second"},
					Actual:   map[string]interface{}{"0": "first", "1": "second"},
					Type:     DiffTypeTypeMismatch,
				},
			},
		},
		{
			name:     "type mismatch - array vs string",
			expected: `{"data": ["value"]}`,
			actual:   `{"data": "value"}`,
			wantDiff: []DiffResult{
				{
					Path:     "data",
					Expected: []interface{}{"value"},
					Actual:   "value",
					Type:     DiffTypeTypeMismatch,
				},
			},
		},
		{
			name:     "type mismatch - object vs string",
			expected: `{"data": {"key": "value"}}`,
			actual:   `{"data": "value"}`,
			wantDiff: []DiffResult{
				{
					Path:     "data",
					Expected: map[string]interface{}{"key": "value"},
					Actual:   "value",
					Type:     DiffTypeTypeMismatch,
				},
			},
		},
		{
			name:     "invalid expected json",
			expected: `{"broken": "json"`,
			actual:   `{"valid": "json"}`,
			wantErr:  true,
		},
		{
			name:     "invalid actual json",
			expected: `{"valid": "json"}`,
			actual:   `{"broken": "json"`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			differ := JSONDiff{matcher: NewDefaultChainMatcher()}
			got, err := differ.Diff(tt.expected, tt.actual)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assertDiffResultsEqual(t, tt.wantDiff, got)
		})
	}
}

func assertDiffResultsEqual(t *testing.T, expected, actual []DiffResult) {
	assert.Equal(t, len(expected), len(actual), "Different number of diff results")

	// Create maps to compare results by path
	expectedMap := make(map[string]DiffResult)
	actualMap := make(map[string]DiffResult)

	for _, e := range expected {
		expectedMap[e.Path] = e
	}

	for _, a := range actual {
		actualMap[a.Path] = a
	}

	// Compare the maps
	assert.Equal(t, expectedMap, actualMap, "Diff results don't match")
}
