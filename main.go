// Feature collection.
package main

import (
	"fmt"
	"github.com/bitcdr/protovalid/examples/geo/build"
	"github.com/bitcdr/protovalid/path"
	"github.com/bitcdr/protovalid/validate"
	"google.golang.org/protobuf/proto"
)

func printFindings(fieldPath *path.FieldPath) {
	for _, f := range fieldPath.Findings {
		fmt.Printf(" - %s: %s\n", fieldPath.FullName, f)
	}

	for _, c := range fieldPath.Childs {
		printFindings(c)
	}
}

func validateMessage(msg proto.Message) {
	fieldPath := validate.Message(msg)
	fmt.Printf("Validation findings for %v:\n", msg)
	printFindings(fieldPath)
}

func main() {
	poiHanover := build.Poi("Hanover", 9.738611, 52.374444)
	validateMessage(poiHanover.Coordinate)
	validateMessage(poiHanover)

	poiInvalid := build.Poi("", 180, -90.1)
	validateMessage(poiInvalid.Coordinate)
	validateMessage(poiInvalid)
}
