package mpatch

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

//go:noinline
func methodA() int {
	x := rand.Int() >> 48
	y := rand.Int() >> 48
	return x + y
}

//go:noinline
func methodB() int {
	x := rand.Int() >> 48
	y := rand.Int() >> 48
	return -(x + y)
}

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
	if methodA() > 0 {
		t.Fatal("The patch did not work")
	}
	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if methodA() < 0 {
		t.Fatal("The unpatch did not work")
	}
}

func TestPatcherUsingReflect(t *testing.T) {
	reflectA := reflect.ValueOf(methodA)
	patch, err := PatchMethodByReflectValue(reflectA, methodB)
	if err != nil {
		t.Fatal(err)
	}
	if methodA() > 0 {
		t.Fatal("The patch did not work")
	}

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if methodA() < 0 {
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
	if methodA() < 0 {
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

func TestRaceTest(t *testing.T) {
	patch, _ := PatchMethod(
		time.Now, func() time.Time {
			loc, _ := time.LoadLocation("America/Sao_Paulo")
			return time.Date(2000, 12, 15, 17, 8, 00, 0, loc)
		})

	defer patch.Unpatch()
}
