# Gomatch

<img align="right" width="147px" src="https://raw.github.com/martinjirku/gomatch/master/logo.png">

[![GoDoc](https://godoc.org/github.com/martinjirku/gomatch?status.svg)](https://godoc.org/github.com/martinjirku/gomatch)
[![Go Report Card](https://goreportcard.com/badge/github.com/martinjirku/gomatch)](https://goreportcard.com/report/github.com/martinjirku/gomatch)
[![codecov](https://codecov.io/gh/martinjirku/gomatch/graph/badge.svg?token=CQP9LNE417)](https://codecov.io/gh/martinjirku/gomatch)

Library created for testing JSON against patterns. The goal was to be able to validate JSON focusing only on parts essential in given test case so tests are more expressive and less fragile. It can be used with both unit tests and functional tests.

When used with Gherkin driven BDD tests it makes scenarios more compact and readable. See [Gherkin example](#gherkin-example)

## Contests

- [Installation](#installation)
- [Basic usage](#basic-usage)
- [Available patterns](#available-patterns)
- [Custom Matchers](#custom-matchers)
- [Golden JSON Sync](#golden-json-sync)
- [Gherkin example](#gherkin-example)
- [License](#license)
- [Credits](#credits)

## Installation

```shell
go get github.com/martinjirku/gomatch
```

## Basic usage

```go

actual := `
{
  "id": 351,
  "name": "John Smith",
  "address": {
    "city": "Boston"
  }
}
`
expected := `
{
  "id": "@number@",
  "name": "John Smith",
  "address": {
    "city": "@string@"
  }
}
`

m := gomatch.NewDefaultJSONMatcher()
ok, err := m.Match(expected, actual)
if ok {
  fmt.Printf("actual JSON matches expected JSON")
} else {
  fmt.Printf("actual JSON does not match expected JSON: %s", err.Error())
}

```

## Available patterns

- `@string@`
- `@number@`
- `@bool@`
- `@array@`
- `@uuid@`
- `@email@`
- `@wildcard@`
- `@date@`
- `@empty@` - checks if the value is empty (null, undefined, empty string, slice, or map or not present)
- `@...@` - unbounded array or object

### Unbounded pattern

It can be used at the end of an array to allow any extra array elements:

```json
["John Smith", "Joe Doe", "@...@"]
```

It can be used at the end of an object to allow any extra keys:

```json
{
  "id": 351,
  "name": "John Smith",
  "@...@": ""
}
```

## Custom Matchers

You can extend gomatch with your own matchers by implementing the ValueMatcher interface:

```go
type ValueMatcher interface {
    // CanMatch returns true if given pattern can be handled by value matcher implementation.
    CanMatch(p interface{}) bool

    // Match performs the matching of given value v.
    // It also expects pattern p so implementation may handle multiple patterns or some DSL.
    Match(p, v interface{}) (bool, error)
}
```

Then, you can create a new JSONMatcher with a chain of your custom matchers:

```go
// Example custom matcher for matching a specific value.
type SpecificValueMatcher struct{}

func (m SpecificValueMatcher) CanMatch(p interface{}) bool {
    _, ok := p.(string)
    return ok && p.(string) == "@specific_value@"
}

func (m SpecificValueMatcher) Match(p, v interface{}) (bool, error) {
    return v == "my_specific_value", nil
}


matchers := []gomatch.ValueMatcher{SpecificValueMatcher{}}
matcher := gomatch.NewJSONMatcher(gomatch.NewChainMatcher(matchers))

actual := `{"value": "my_specific_value"}`
expected := `{"value": "@specific_value@"}`

ok, err := matcher.Match(expected, actual)
// ...
```

## Golden JSON Sync

`goldenJSONSync.Sync` helps to synchronize expected JSON (golden file) with actual JSON. It merges the structure of the actual JSON into the golden JSON, preserving the matcher patterns from the golden file. This is particularly useful for updating expected results in tests when the structure of the actual data changes but the matching criteria remain the same.

The `goldenJSONSync.Sync` merges the values from actual into golden while keeping the matcher patterns. The Marshaler function is used to format the resulting JSON. This helps keep your golden files up-to-date with the actual data structure without losing the flexibility of your matchers.

Golden JSON:

```json
{
  "normalValue": 1,
  "string": "@string@",
  "date": "@date@",
  "array": [1, 2, 3],
  "mapMatcher": { "a": 1, "b": 2 },
  "bool": "@bool@",
  "empty": "@empty@",
  "email": "@email@",
  "uuid": "@uuid@"
}
```

Actual JSON:

```json
{
  "normalValue": 5,
  "string": "test_string",
  "date": "2024-10-27",
  "array": [3, 2, 1],
  "mapMatcher": { "a": 2, "b": 1 },
  "bool": true,
  "empty": null,
  "email": "[email address removed]",
  "uuid": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
}
```

New Golden JSON:

```json
{
  "normalValue": 5,
  "string": "@string@",
  "date": "@date@",
  "array": [3, 2, 1],
  "mapMatcher": { "a": 2, "b": 1 },
  "bool": "@bool@",
  "empty": "@empty@",
  "email": "@email@",
  "uuid": "@uuid@"
}
```

## Gherkin example

Gomatch was created to use it together with tools like [GODOG](https://github.com/DATA-DOG/godog).
The goal was to be able to validate JSON response focusing only on parts essential in given scenario.

```gherkin
Feature: User management API
  In order to provide GUI for user management
  As a frontent developer
  I need to be able to create, retrive, update and delete users

  Scenario: Get list of users sorted by username ascending
    Given the database contains users:
    | Username   | Email                  |
    | john.smith | john.smith@example.com |
    | alvin34    | alvin34@example.com    |
    | mike1990   | mike.jones@example.com |
    When I send "GET" request to "/v1/users?sortBy=username&sortDir=asc"
    Then the response code should be 200
    And the response body should match json:
    """
    {
      "items": [
        {
          "username": "alvin34",
          "@...@": ""
        },
        {
          "username": "john.smith",
          "@...@": ""
        },
        {
          "username": "mike1990",
          "@...@": ""
        }
      ],
      "@...@": ""
    }
    """
```

## License

This library is distributed under the MIT license. Please see the LICENSE file.

## Credits

This library was inspired by [PHP Matcher](https://github.com/coduo/php-matcher).

This library is also fork of the [jfilipczyk/gomatch](https://github.com/jfilipczyk/gomatch).

### Logo

The Go gopher was designed by Renee French. (http://reneefrench.blogspot.com/).
Gomatch logo was based on a gopher created by Takuya Ueda (https://twitter.com/tenntenn). Licensed under the [Creative Commons 3.0 Attributions license](http://creativecommons.org/licenses/by/3.0/deed.en). Gopher eyes were changed.
