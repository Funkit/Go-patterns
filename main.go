package main

import (
	"fmt"

	"github.com/Funkit/Go-patterns/generics"
)

func main() {

	fmt.Println("-------------------------Generics : builder pattern-------------------------")
	var acu generics.AntennaController

	var acuType string
	acuType = "ControllerTwo"

	switch acuType {
	case "ControllerOne":
		acu = generics.NewAntennaController[generics.ControllerOne]("path/to/config/file")
	case "ControllerTwo":
		acu = generics.NewAntennaController[generics.ControllerTwo]("path/to/config/file")
	}

	fmt.Println(acu.Type())
}
