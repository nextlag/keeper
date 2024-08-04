package l

import (
	"fmt"
	"runtime"
	"strings"
)

// WrapErr wraps the provided error with the package and function name of the caller.
func WrapErr(err error) error {
	if err == nil {
		return nil
	}

	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return err
	}
	fullFuncName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullFuncName, "/")
	funcName := parts[len(parts)-1]
	packageAndFunc := strings.Split(funcName, ".")
	packageName := strings.Join(packageAndFunc[:len(packageAndFunc)-1], ".")
	funcName = packageAndFunc[len(packageAndFunc)-1]
	return fmt.Errorf("%s.%s: %w", packageName, funcName, err)
}
