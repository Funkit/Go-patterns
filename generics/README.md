# Generics

## Generic type builder 

### The problem

When having a function building a specific type, it can be useful to have a generic implementation that would only build types that satisfy a given interface.

For instance, consider the following:

- We have an interface called `AntennaController`, with a `LoadConfig` method. This method mutates the object, so the implementations have pointer receivers.
- Two distinct structs `ControllerOne` and `ControllerTwo` satisfy the interface.
- In the main, depending on a string value (from a configuration file for instance), we have to build a specific implementation of `AntennaController`, either `ControllerOne` or `ControllerTwo`.

```
type AntennaController interface {
	LoadConfig(filePath string)
}

type ControllerOne struct{}

func (c1 *ControllerOne) LoadConfig(filePath string) {
    //Load data from file, mutate object
}

type ControllerTwo struct{}

func (c2 *ControllerTwo) LoadConfig(filePath string) {
    //Load data from file, mutate object
}

func main() {
    var acu AntennaController

    var acuType string
    acuType := "ControllerTwo"

    switch acuType {
        case "ControllerOne":
            var c1 ControllerOne{}
            c1.LoadConfig("path/to/config/file")
            acu = c1
        case "ControllerTwo":
            var c2 ControllerTwo{}
            c2.LoadConfig("path/to/config/file")
            acu = c2
    }
}
```

Now let's say we want to write a generic `NewAntennaController` function to save us writing three lines for every case. The naive implementation would be:

```
NewAntennaController[T AntennaController](filePath string) AntennaController {
    var t T
    t.LoadConfig(filepath)

    return t
}
```

Unfortunately, because `ControllerOne` and `ControllerTwo` use pointer receivers for their methods, you cannot write `NewAntennaController[ControllerOne](filepath)` since it is not the object itself but the **pointer** that actually satisfies the interface.

You might then be tempted to write `NewAntennaController[*ControllerOne](filepath)`, but when doing that, **the pointer is not initialized**. As such, the naive implementation above will fail.

### The solution

It is possible to use an intermediate interface, `AntennaControllerConstraint`, and to rewrite `NewAntennaController` like this :

```
type AntennaControllerConstraint[P any] interface {
	AntennaController
	*P
}

func NewAntennaController[T any, P AntennaControllerConstraint[T]](filePath string) AntennaController {
	var acuPointer P
	var acu T
	acuPointer = &acu

	acuPointer.LoadConfig(filePath)

	return acuPointer
}
```

`AntennaControllerConstraint`, for a given type T, is an interface that is at the same time a pointer to T and an `AntennaController`. This version of `NewAntennaController` takes as type parameter any type that match this interface.

You can now call `NewAntennaController[ControllerOne](filepath)`, as long as `ControllerOne` has pointer receivers for its methods.

The full pseudo code:

```
type AntennaController interface {
	LoadConfig(filePath string)
	Type() string
}

type ControllerOne struct{}

func (c1 *ControllerOne) LoadConfig(filePath string) {}

func (c1 *ControllerOne) Type() string {
	return "ControllerOne"
}

type ControllerTwo struct{}

func (c2 *ControllerTwo) LoadConfig(filePath string) {}

func (c2 *ControllerTwo) Type() string {
	return "ControllerTwo"
}

type AntennaControllerConstraint[P any] interface {
	AntennaController
	*P
}

func NewAntennaController[T any, P AntennaControllerConstraint[T]](filePath string) AntennaController {
	var acuPointer P
	var acu T
	acuPointer = &acu

	acuPointer.LoadConfig(filePath)

	return acuPointer
}
```



