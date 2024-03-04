//go:build !javascript
// +build !javascript

package util

import (
	"errors"
	"runtime/debug"
)

// RecoverPanic recovers a panic.
func RecoverPanic(err *error) {
	if e := recover(); nil != e {
		stack := debug.Stack()
		errMsg := ""
		switch x := e.(type) {
		case error:
			errMsg = x.Error()
		case string:
			errMsg = x
		default:
			errMsg = "unknown panic"
		}
		if nil != err {
			*err = errors.New("PANIC RECOVERED: " + errMsg + "\n\t" + string(stack) + "\n")
		}
	}
}
