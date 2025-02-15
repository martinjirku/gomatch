package gomatch

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func NewErrGomatch(err error, path []interface{}, expected, actual interface{}, key string) error {
	if err == nil {
		return nil
	}
	return ErrGomatch{Path: path, err: err, Expected: expected, Provided: actual, Key: key}
}

type ErrGomatch struct {
	Path     []interface{}
	Key      string
	Expected any
	Provided any
	err      error
}

func (e ErrGomatch) Error() string {
	expected := "null"
	if e.Expected != nil {
		expected = valueOf(e.Expected)
	}
	provided := "null"
	if e.Provided != nil {
		provided = valueOf(e.Provided)
	}
	return fmt.Sprintf("%s at %q. expected: %s, provided: %s", e.err, pathToString(e.Path), expected, provided)
}
func (e ErrGomatch) Unwrap() error {
	return e.err
}

func pathToString(path []interface{}) string {
	var b bytes.Buffer
	b.WriteRune('.')
	for _, p := range path {
		switch v := p.(type) {
		case int:
			b.WriteString(fmt.Sprintf("[%d]", v))
		default:
			if b.Len() > 1 {
				b.WriteRune('.')
			}
			b.WriteString(v.(string))
		}
	}
	return b.String()
}

func valueOf(v interface{}) string {
	val, _ := json.Marshal(v)
	return string(val)
}
