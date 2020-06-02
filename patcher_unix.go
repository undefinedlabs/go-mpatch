// +build !windows

package mpatch

import (
	"reflect"
	"syscall"
	"unsafe"
)

var writeAccess = syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC
var readAccess = syscall.PROT_READ | syscall.PROT_EXEC

//go:nosplit
func getMemorySliceFromUintptr(p uintptr, length int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p,
		Len:  length,
		Cap:  length,
	}))
}

//go:nosplit
func callMProtect(addr unsafe.Pointer, length int, prot int) error {
	for p := uintptr(addr) & ^(uintptr(pageSize - 1)); p < uintptr(addr)+uintptr(length); p += uintptr(pageSize) {
		page := getMemorySliceFromUintptr(p, pageSize)
		err := syscall.Mprotect(page, prot)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyDataToPtr(ptr unsafe.Pointer, data []byte) error {
	dataLength := len(data)
	ptrByteSlice := getMemorySliceFromPointer(ptr, len(data))
	err := callMProtect(ptr, dataLength, writeAccess)
	if err != nil {
		return err
	}
	copy(ptrByteSlice, data[:])
	err = callMProtect(ptr, dataLength, readAccess)
	if err != nil {
		return err
	}
	return nil
}
