// +build arm64

package mpatch

import (
	"errors"
	"runtime"
)

// Gets the jump function rewrite bytes
func getJumpFuncBytes(to uintptr) ([]byte, error) {
	return nil, errors.New(fmt.Sprintf("Unsupported architecture: %s", runtime.GOOS))
}