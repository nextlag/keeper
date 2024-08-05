package l

import (
	"fmt"
)

// WrapErr wraps the provided error with the package and function name of the caller.
func WrapErr(err error) error {
	return fmt.Errorf("%w", err)
}
