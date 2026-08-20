package withstack

import (
	"errors"
	"fmt"

	cErr "github.com/cockroachdb/errors"
)

type WithStackError struct {
	err error // -> actual WithStack error
}

func (w *WithStackError) Error() string {
	return fmt.Sprintf("%+v", w.err)
}

func (w *WithStackError) Unwrap() error {
	return w.err
}

func Wrap(err error) *WithStackError {
	target := new(WithStackError)
	if errors.As(err, &target) {
		return target
	}

	return &WithStackError{
		err: cErr.WithStackDepth(err, 1),
	}
}

func WrapWithDepth(err error, depth int) *WithStackError {
	return &WithStackError{
		err: cErr.WithStackDepth(err, depth),
	}
}
