package generics

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
