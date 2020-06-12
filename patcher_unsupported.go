// +build !386
// +build !amd64

package mpatch

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

// Gets the jump function rewrite bytes
//go:nosplit
func getJumpFuncBytes(to unsafe.Pointer) ([]byte, error) {
	return nil, errors.New(fmt.Sprintf("unsupported architecture: %s", runtime.GOARCH))
}
