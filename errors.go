package gomatch

import (
	"bytes"
	"fmt"
)

func NewErrGomatch(err error, path []interface{}) error {
	if err == nil {
		return nil
	}
	return ErrGomatch{path: path, err: err}
}

type ErrGomatch struct {
	path []interface{}
	err  error
}

func (e ErrGomatch) Error() string {
	return fmt.Sprintf("%s at %q", e.err, pathToString(e.path))
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
