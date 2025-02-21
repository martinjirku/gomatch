package gomatch_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/martinjirku/gomatch"
)

func TestSyncGoldenJSON(t *testing.T) {
	testcase := []struct {
		title  string
		golden string
		actual string
		result string
		error  bool
	}{
		{
			title: "Different matchers",
			golden: `{
	"normalValue": 1,
	"string": "@string@",
	"date": "@date@",
	"array": [1, 2, 3],
	"mapMatcher": {"a": 1, "b": 2},
	"bool": "@bool@",
	"empty": "@empty@",
	"email": "@email@",
	"uuid": "@uuid@"
}`,
			actual: `{
	"missingInGolden": 1,
	"normalValue": 5,
	"string": "@string@",
	"date": "@date@",
	"array": [3, 2, 1],
	"mapMatcher": {"a": 2, "b": 1},
	"bool": "@bool@",
	"empty": "@empty@",
	"email": "@email@",
	"uuid": "@uuid@"
}`,
			result: `{
	"missingInGolden": 1,
	"normalValue": 5,
	"string": "@string@",
	"date": "@date@",
	"array": [3, 2, 1],
	"mapMatcher": {"a": 2, "b": 1},
	"bool": "@bool@",
	"empty": "@empty@",
	"email": "@email@",
	"uuid": "@uuid@"
}`,
		},
		{
			title:  "empty golden and new field",
			golden: `{}`,
			actual: `{"a": 1}`,
			result: `{"a": 1}`,
		},
		{
			title:  "full golden and empty actual",
			golden: `{"a": 1}`,
			actual: `{}`,
			result: `{}`,
		},
		{
			title:  "full golden and new field",
			golden: `{"a": 1, "b": 2}`,
			actual: `{"a": 3}`,
			result: `{"a": 3}`,
		},
		{
			title:  "partial golden and new field",
			golden: `{"a": 1}`,
			actual: `{"a": 3, "b": 4}`,
			result: `{"a": 3, "b": 4}`,
		},
		{
			title:  "array with different length",
			golden: `{"a": 1, "b": null, "c": 3.14}`,
			actual: `{"a": 5, "b": null, "c": 2.71}`,
			result: `{"a": 5, "b": null, "c": 2.71}`,
		},
		{
			title:  "invalid actual",
			golden: `{"a": 1}`,
			actual: `invalid json`,
			result: `{"a": 1}`,
			error:  true,
		},
		{
			title:  "invalid golden",
			golden: `invalid json`,
			actual: `{"a": 1}`,
			result: `{"a": 1}`,
		},
		{
			title:  "bounded object",
			golden: `{"a": 1, "@...@": null}`,
			actual: `{"a": 1, "@...@": null, "b": 2}`,
			result: `{"a": 1, "@...@": null}`,
		},
		{
			title:  "bounded array",
			golden: `[1, 2, 3, "@...@"]`,
			actual: `[1, 2, 3, 4, 5]`,
			result: `[1, 2, 3, "@...@"]`,
		},
		{
			title:  "unbounded array with different length",
			golden: `[1, 2, 3, 4]`,
			actual: `[1, 2, 3, 4, 5, 6]`,
			result: `[1, 2, 3, 4, 5, 6]`,
		},
		{
			title:  "unbounded array with less items than golden",
			golden: `[1, 2, 3, 4, 5, 6, 7]`,
			actual: `[1, 2, 3, 4, 5]`,
			result: `[1, 2, 3, 4, 5]`,
		},
		{
			title:  "different types",
			golden: `{"a": 1}`,
			actual: `{"a": "text"}`,
			result: `{"a": "text"}`,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.title, func(t *testing.T) {
			goldenJSONSync := gomatch.NewGoldenJSONSync()
			goldenJSONSync.Marshaler(func(v any) ([]byte, error) {
				return json.Marshal(v)
			})
			matcher := gomatch.NewJSONMatcher(gomatch.NewChainMatcher([]gomatch.ValueMatcher{}))
			// Given
			res, err := goldenJSONSync.Sync(tc.golden, tc.actual)
			if tc.error && err == nil {
				t.Errorf("Expected error %v, got %v", tc.error, err)
			}
			if _, err := matcher.Match(res, tc.result); err != nil {
				t.Errorf("Expected result %v, got %v. Differents: %s", tc.result, res, err.Error())
			}
		})
	}
}

func TestSyncGoldenJSON_Marshaler(t *testing.T) {
	goldenJSONSync := gomatch.NewGoldenJSONSync()
	goldenJSONSync.Marshaler(func(v any) ([]byte, error) {
		return []byte{}, fmt.Errorf("error")
	})
	actual := `{"a": 1}`
	golden := `{"a": 2}`
	res, err := goldenJSONSync.Sync(golden, actual)
	if err == nil {
		t.Errorf("Expected error, got %v", res)
	}
	if res != golden {
		t.Errorf("Expected result %v, got %v", actual, res)
	}
}
