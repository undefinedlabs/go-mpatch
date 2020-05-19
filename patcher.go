package mpatch // import "github.com/undefinedlabs/go-mpatch"

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"syscall"
	"unsafe"
)

type (
	Patch struct {
		targetBytes []byte
		target      *reflect.Value
		redirection *reflect.Value
	}
	sliceHeader struct {
		Data unsafe.Pointer
		Len  int
		Cap  int
	}
)

//go:linkname getInternalPtrFromValue reflect.(*Value).pointer
func getInternalPtrFromValue(v *reflect.Value) unsafe.Pointer

var (
	patchLock = sync.Mutex{}
	patches   = make(map[unsafe.Pointer]*Patch)
	pageSize  = syscall.Getpagesize()
)

func PatchMethod(target, redirection interface{}) (*Patch, error) {
	tValue := getValueFrom(target)
	rValue := getValueFrom(redirection)
	err := isPatchable(&tValue, &rValue)
	if err != nil {
		return nil, err
	}
	patch := &Patch{target: &tValue, redirection: &rValue}
	err = applyPatch(patch)
	if err != nil {
		return nil, err
	}
	return patch, nil
}
func PatchInstanceMethodByName(target reflect.Type, methodName string, redirection interface{}) (*Patch, error) {
	method, ok := target.MethodByName(methodName)
	if !ok && target.Kind() == reflect.Struct {
		target = reflect.PtrTo(target)
		method, ok = target.MethodByName(methodName)
	}
	if !ok {
		return nil, errors.New(fmt.Sprintf("Method '%v' not found", methodName))
	}
	return PatchMethodByReflect(method, redirection)
}
func PatchMethodByReflect(target reflect.Method, redirection interface{}) (*Patch, error) {
	tValue := &target.Func
	rValue := getValueFrom(redirection)
	err := isPatchable(tValue, &rValue)
	if err != nil {
		return nil, err
	}
	patch := &Patch{target: tValue, redirection: &rValue}
	err = applyPatch(patch)
	if err != nil {
		return nil, err
	}
	return patch, nil
}
func PatchMethodWithMakeFunc(target reflect.Method, fn func(args []reflect.Value) (results []reflect.Value)) (*Patch, error) {
	rValue := reflect.MakeFunc(target.Type, fn)
	return PatchMethodByReflect(target, rValue)
}

func (p *Patch) Patch() error {
	if p == nil {
		return errors.New("patch is nil")
	}
	err := isPatchable(p.target, p.redirection)
	if err != nil {
		return err
	}
	err = applyPatch(p)
	if err != nil {
		return err
	}
	return nil
}
func (p *Patch) Unpatch() error {
	if p == nil {
		return errors.New("patch is nil")
	}
	return applyUnpatch(p)
}

func isPatchable(target, redirection *reflect.Value) error {
	if target.Kind() != reflect.Func || redirection.Kind() != reflect.Func {
		return errors.New("the target and/or redirection is not a Func")
	}
	if target.Type() != redirection.Type() {
		return errors.New(fmt.Sprintf("the target and/or redirection doesn't have the same type: %s != %s", target.Type(), redirection.Type()))
	}
	if _, ok := patches[getSafePointer(target)]; ok {
		return errors.New("the target is already patched")
	}
	return nil
}

func applyPatch(patch *Patch) error {
	patchLock.Lock()
	defer patchLock.Unlock()
	tPointer := getSafePointer(patch.target)
	rPointer := getInternalPtrFromValue(patch.redirection)
	rPointerJumpBytes, err := getJumpFuncBytes(rPointer)
	if err != nil {
		return err
	}
	tPointerBytes := getMemorySliceFromPointer(tPointer, len(rPointerJumpBytes))
	targetBytes := make([]byte, len(tPointerBytes))
	copy(targetBytes, tPointerBytes)
	err = copyDataToPtr(tPointer, rPointerJumpBytes)
	if err != nil {
		return err
	}
	patch.targetBytes = targetBytes
	patches[tPointer] = patch
	return nil
}

func applyUnpatch(patch *Patch) error {
	patchLock.Lock()
	defer patchLock.Unlock()
	if patch.targetBytes == nil || len(patch.targetBytes) == 0 {
		return errors.New("the target is not patched")
	}
	tPointer := getSafePointer(patch.target)
	if _, ok := patches[tPointer]; !ok {
		return errors.New("the target is not patched")
	}
	delete(patches, tPointer)
	err := copyDataToPtr(tPointer, patch.targetBytes)
	if err != nil {
		return err
	}
	return nil
}

func getValueFrom(data interface{}) reflect.Value {
	if cValue, ok := data.(reflect.Value); ok {
		return cValue
	} else {
		return reflect.ValueOf(data)
	}
}

func getMemorySliceFromPointer(p unsafe.Pointer, length int) []byte {
	return *(*[]byte)(unsafe.Pointer(&sliceHeader{
		Data: p,
		Len:  length,
		Cap:  length,
	}))
}

func getSafePointer(value *reflect.Value) unsafe.Pointer {
	p := getInternalPtrFromValue(value)
	if p != nil {
		p = *(*unsafe.Pointer)(p)
	}
	return p
}
