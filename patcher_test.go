package mpatch

import (
	"reflect"
	"runtime"
	"testing"
	"time"
)

//go:noinline
func methodA() int { return 1 }

//go:noinline
func methodB() int { return 2 }

type myStruct struct {
}

//go:noinline
func (s *myStruct) Method() int {
	return 1
}

//go:noinline
func (s myStruct) ValueMethod() int {
	return 1
}

func TestPatcher(t *testing.T) {
	patch, err := PatchMethod(methodA, methodB)
	if err != nil {
		t.Fatal(err)
	}
	if methodA() != 2 {
		t.Fatal("The patch did not work")
	}

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if methodA() != 1 {
		t.Fatal("The unpatch did not work")
	}
}

func TestPatcherUsingReflect(t *testing.T) {
	reflectA := reflect.ValueOf(methodA)
	patch, err := PatchMethodByReflectValue(reflectA, methodB)
	if err != nil {
		t.Fatal(err)
	}
	if methodA() != 2 {
		t.Fatal("The patch did not work")
	}

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if methodA() != 1 {
		t.Fatal("The unpatch did not work")
	}
}

func TestPatcherUsingMakeFunc(t *testing.T) {
	reflectA := reflect.ValueOf(methodA)
	patch, err := PatchMethodWithMakeFuncValue(reflectA,
		func(args []reflect.Value) (results []reflect.Value) {
			return []reflect.Value{reflect.ValueOf(42)}
		})
	if err != nil {
		t.Fatal(err)
	}
	if methodA() != 42 {
		t.Fatal("The patch did not work")
	}

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if methodA() != 1 {
		t.Fatal("The unpatch did not work")
	}
}

func TestInstancePatcher(t *testing.T) {
	mStruct := myStruct{}

	var patch *Patch
	var err error
	patch, err = PatchInstanceMethodByName(reflect.TypeOf(mStruct), "Method", func(m *myStruct) int {
		patch.Unpatch()
		defer patch.Patch()
		return 41 + m.Method()
	})
	if err != nil {
		t.Fatal(err)
	}

	if mStruct.Method() != 42 {
		t.Fatal("The patch did not work")
	}
	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if mStruct.Method() != 1 {
		t.Fatal("The unpatch did not work")
	}
}

func TestInstanceValuePatcher(t *testing.T) {
	mStruct := myStruct{}

	var patch *Patch
	var err error
	patch, err = PatchInstanceMethodByName(reflect.TypeOf(mStruct), "ValueMethod", func(m myStruct) int {
		patch.Unpatch()
		defer patch.Patch()
		return 41 + m.Method()
	})
	if err != nil {
		t.Fatal(err)
	}

	if mStruct.ValueMethod() != 42 {
		t.Fatal("The patch did not work")
	}
	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if mStruct.ValueMethod() != 1 {
		t.Fatal("The unpatch did not work")
	}
}

var slice []int

//go:noinline
func TestGarbageCollectorExperiment(t *testing.T) {

	for i := 0; i < 10000000; i++ {
		slice = append(slice, i)
	}
	go func() {
		var sl []int
		for i := 0; i < 10000000; i++ {
			sl = append(slice, i)
		}
		_ = sl
	}()
	<-time.After(time.Second)

	aVal := methodA
	ptr01 := reflect.ValueOf(aVal).Pointer()
	slice = nil
	runtime.GC()
	for i := 0; i < 10000000; i++ {
		slice = append(slice, i)
	}
	go func() {
		var sl []int
		for i := 0; i < 10000000; i++ {
			sl = append(slice, i)
		}
		_ = sl
	}()
	<-time.After(time.Second)
	slice = nil
	runtime.GC()
	ptr02 := reflect.ValueOf(aVal).Pointer()

	if ptr01 != ptr02 {
		t.Fail()
	}
}
