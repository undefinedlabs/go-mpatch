# go-mpatch
Go library for monkey patching

## Compatibility

- **Go version:** tested from `go1.7` to `go1.19`
- **Architectures:** `x86`, `amd64`
- **Operating systems:** tested in `macos`, `linux` and `windows`. 

## Features

- Can patch package functions, instance functions (by pointer or by value), and create new functions from scratch.

## Limitations

- Target functions could be inlined, making those functions unpatcheables. You can use `//go:noinline` directive or build with the `gcflags=-l`
to disable inlining at compiler level.

- Write permission to memory pages containing executable code is needed, some operating systems could restrict this access.

- Not thread safe.

## Usage

### Patching a func
```go
//go:noinline
func methodA() int { return 1 }

//go:noinline
func methodB() int { return 2 }

func TestPatcher(t *testing.T) {
	patch, err := mpatch.PatchMethod(methodA, methodB)
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
```

### Patching using `reflect.ValueOf`
```go
//go:noinline
func methodA() int { return 1 }

//go:noinline
func methodB() int { return 2 }

func TestPatcherUsingReflect(t *testing.T) {
	reflectA := reflect.ValueOf(methodA)
	patch, err := mPatch.PatchMethodByReflectValue(reflectA, methodB)
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
```

### Patching creating a new func at runtime
```go
//go:noinline
func methodA() int { return 1 }

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
```

### Patching an instance func
```go
type myStruct struct {
}

//go:noinline
func (s *myStruct) Method() int {
	return 1
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
```

### Patching an instance func by Value
```go
type myStruct struct {
}

//go:noinline
func (s myStruct) ValueMethod() int {
	return 1
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
```

> Library inspired by the blog post: https://bou.ke/blog/monkey-patching-in-go/
