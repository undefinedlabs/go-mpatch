# go-mpatch
Go library for monkey patching

## Compatibility

- **Go version:** tested from `go1.7` to `go1.15-beta`
- **Architectures:** `x86`, `amd64`
- **Operating systems:** tested in `macos`, `linux` and `windows`. 
Write permission to memory pages containing executable code is needed, some operating systems could restrict this access.

## Features

- Can patch package functions, instance functions (by pointer or by value), and create new functions from scratch.

## Limitations

- Target functions could be inlined, making those functions unpatcheables. You can use `//go:noinline` directive or build with the `gcflags=-l`
to disable inlining at compiler level.

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
	patch, err := mPatch.PatchMethodByReflect(reflectA, methodB)
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
	patch, err := PatchMethodWithMakeFunc(reflectA,
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

> Library inspired by the blog post: https://bou.ke/blog/monkey-patching-in-go/
