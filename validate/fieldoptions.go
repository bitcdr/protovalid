// Package validate provides validation of protobuf messages.
package validate

import (
	"fmt"
	"regexp"

	"github.com/bitcdr/protovalid/extension"
	"github.com/bitcdr/protovalid/path"
	pb "github.com/bitcdr/protovalid/valid"
	"github.com/bitcdr/protovalid/wrappers"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Message validates the protobuf message and checks the field options and returns the validation findings.
func Message(msg proto.Message) *path.FieldPath {
	fieldPath := &path.FieldPath{}
	validateNestedMessage(fieldPath, msg.ProtoReflect())
	return fieldPath
}

func validateNestedMessage(fieldPath *path.FieldPath, m protoreflect.Message) {
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		fieldName := string(fd.Name())
		fieldPathChild := path.GetChild(fieldPath, fieldName)

		// Field is repeated / list
		if fd.IsList() {
			validateList(fieldPathChild, fd, v.List())

			// List elements
			for i := 0; i < v.List().Len(); i++ {
				fieldNameWithIndex := fmt.Sprintf("%s[%d]", fieldName, i)
				fieldPathChildWithIndex := path.GetChild(fieldPath, fieldNameWithIndex)
				elemValue := v.List().Get(i)

				switch fd.Kind() {
				case protoreflect.DoubleKind:
					validateDouble(fieldPathChildWithIndex, fd, elemValue.Float())
				case protoreflect.Int32Kind:
					validateInt32(fieldPathChildWithIndex, fd, int32(elemValue.Int()))
				case protoreflect.StringKind:
					validateString(fieldPathChildWithIndex, fd, elemValue.String())
				case protoreflect.MessageKind:
					// Wrapper types
					switch fd.Message().FullName() {
					case "google.protobuf.DoubleValue":
						validateDouble(fieldPathChildWithIndex, fd, wrappers.GetFloat64(fd, elemValue))
					case "google.protobuf.Int32Value":
						validateInt32(fieldPathChildWithIndex, fd, wrappers.GetInt32(fd, elemValue))
					case "google.protobuf.StringValue":
						validateString(fieldPathChildWithIndex, fd, wrappers.GetString(fd, elemValue))
					default:
						// Descend into fields of type Message
						validateNestedMessage(fieldPathChildWithIndex, elemValue.Message())
					}
				}
			}

			return true
		}

		// Field is not repeated / list or map
		switch fd.Kind() {
		case protoreflect.DoubleKind:
			validateDouble(fieldPathChild, fd, v.Float())
		case protoreflect.Int32Kind:
			validateInt32(fieldPathChild, fd, int32(v.Int()))
		case protoreflect.StringKind:
			validateString(fieldPathChild, fd, v.String())
		case protoreflect.MessageKind:
			// Wrapper types
			switch fd.Message().FullName() {
			case "google.protobuf.Int32Value":
				validateInt32(fieldPathChild, fd, wrappers.GetInt32(fd, v))
			case "google.protobuf.DoubleValue":
				validateDouble(fieldPathChild, fd, wrappers.GetFloat64(fd, v))
			case "google.protobuf.StringValue":
				validateString(fieldPathChild, fd, wrappers.GetString(fd, v))
			default:
				// Descend into fields of type Message
				validateNestedMessage(fieldPathChild, v.Message())
			}
		}

		return true
	})
}

func validateDouble(fieldPath *path.FieldPath, fd protoreflect.FieldDescriptor, v float64) {
	if ok, c := extension.GetFieldConstraintsDouble(fd); ok {
		// Min
		switch x := c.Min.(type) {
		case *pb.FieldConstraints_DoubleConstraints_Gte:
			if v < x.Gte-c.Eps {
				path.AddFinding(fieldPath, fmt.Sprintf("value %f is less than %f", v, x.Gte), c.Level, c.Msg)
			}
		case *pb.FieldConstraints_DoubleConstraints_Gt:
			if v <= x.Gt-c.Eps {
				path.AddFinding(fieldPath, fmt.Sprintf("value %f is less than or equal %f", v, x.Gt), c.Level, c.Msg)
			}
		}

		// Max
		switch x := c.Max.(type) {
		case *pb.FieldConstraints_DoubleConstraints_Lte:
			if v > x.Lte+c.Eps {
				path.AddFinding(fieldPath, fmt.Sprintf("value %f is greater than %f", v, x.Lte), c.Level, c.Msg)
			}
		case *pb.FieldConstraints_DoubleConstraints_Lt:
			if v <= x.Lt+c.Eps {
				path.AddFinding(fieldPath, fmt.Sprintf("value %f is greater than or equal %f", v, x.Lt), c.Level, c.Msg)
			}
		}
	}
}

func validateInt32(fieldPath *path.FieldPath, fd protoreflect.FieldDescriptor, v int32) {
	if ok, c := extension.GetFieldConstraintsInt32(fd); ok {
		// Min
		switch x := c.Min.(type) {
		case *pb.FieldConstraints_Int32Constraints_Gte:
			if v < x.Gte {
				path.AddFinding(fieldPath, fmt.Sprintf("value %d is less than %d", v, x.Gte), c.Level, c.Msg)
			}
		case *pb.FieldConstraints_Int32Constraints_Gt:
			if v <= x.Gt {
				path.AddFinding(fieldPath, fmt.Sprintf("value %d is less than or equal %d", v, x.Gt), c.Level, c.Msg)
			}
		}

		// Max
		switch x := c.Max.(type) {
		case *pb.FieldConstraints_Int32Constraints_Lte:
			if v > x.Lte {
				path.AddFinding(fieldPath, fmt.Sprintf("value %d is greater than %d", v, x.Lte), c.Level, c.Msg)
			}
		case *pb.FieldConstraints_Int32Constraints_Lt:
			if v <= x.Lt {
				path.AddFinding(fieldPath, fmt.Sprintf("value %d is greater than or equal %d", v, x.Lt), c.Level, c.Msg)
			}
		}

		// In, not in, equal
		if len(c.In) > 0 {
			hit := false
			for _, in := range c.In {
				if v == in {
					hit = true
				}
			}

			if !hit {
				path.AddFinding(fieldPath, fmt.Sprintf("value %d is not in %v", v, c.In), c.Level, c.Msg)
			}
		} else if len(c.NotIn) > 0 {
			hit := false
			for _, notIn := range c.NotIn {
				if v == notIn {
					hit = true
				}
			}

			if hit {
				path.AddFinding(fieldPath, fmt.Sprintf("value %d is in %v", v, c.NotIn), c.Level, c.Msg)
			}
		} else if c.Eq != 0 && v != c.Eq {
			path.AddFinding(fieldPath, fmt.Sprintf("value %d is not %d", v, c.Eq), c.Level, c.Msg)
		}
	}
}

func validateString(fieldPath *path.FieldPath, fd protoreflect.FieldDescriptor, v string) {
	if ok, c := extension.GetFieldConstraintsString(fd); ok {
		length := uint32(len(v))

		// Empty and not empty
		switch x := c.EmptyType.(type) {
		case *pb.FieldConstraints_StringConstraints_Empty:
			if x.Empty && length > 0 {
				path.AddFinding(fieldPath, fmt.Sprintf("value %s is not empty", v), c.Level, c.Msg)
			}
		case *pb.FieldConstraints_StringConstraints_NotEmpty:
			if x.NotEmpty && length == 0 {
				path.AddFinding(fieldPath, fmt.Sprintf("value %s is empty", v), c.Level, c.Msg)
			}
		}

		// Length equal, gte, lte
		if c.LenEq > 0 && length != c.LenEq {
			path.AddFinding(fieldPath, fmt.Sprintf("length of value %v is not %d", v, c.LenEq), c.Level, c.Msg)
		} else {
			if c.LenGte > 0 && length < c.LenGte {
				path.AddFinding(fieldPath, fmt.Sprintf("length of value %v is less than %d", v, c.LenEq), c.Level, c.Msg)
			}

			if c.LenLte > 0 && length > c.LenLte {
				path.AddFinding(fieldPath, fmt.Sprintf("length of value %v is greater than %d", v, c.LenEq), c.Level, c.Msg)
			}
		}

		// Equal, regex pattern
		if len(c.Eq) > 0 && c.Eq != v {
			path.AddFinding(fieldPath, fmt.Sprintf("value %v is not %s", v, c.Eq), c.Level, c.Msg)
		} else if len(c.Pattern) > 0 {
			regex, err := regexp.Compile(c.Pattern)
			if err != nil {
				path.AddFinding(fieldPath, fmt.Sprintf("invalid regex %s", c.Pattern), c.Level, c.Msg)
			} else if !regex.MatchString(v) {
				path.AddFinding(fieldPath, fmt.Sprintf("value %s doesn't match regex %s", v, c.Pattern), c.Level, c.Msg)
			}
		}
	}
}

func validateList(fieldPath *path.FieldPath, fd protoreflect.FieldDescriptor, v protoreflect.List) {
	if ok, c := extension.GetFieldConstraintsRepeated(fd); ok {
		length := uint32(v.Len())

		// Empty and not empty
		switch x := c.EmptyType.(type) {
		case *pb.FieldConstraints_RepeatedConstraints_Empty:
			if x.Empty && length > 0 {
				path.AddFinding(fieldPath, "list is not empty", c.Level, c.Msg)
			}
		case *pb.FieldConstraints_RepeatedConstraints_NotEmpty:
			if x.NotEmpty && length == 0 {
				path.AddFinding(fieldPath, "list is empty", c.Level, c.Msg)
			}
		}

		// Length equal, gte, lte
		if c.LenEq > 0 && length != c.LenEq {
			path.AddFinding(fieldPath, fmt.Sprintf("length of value list is not %d", c.LenEq), c.Level, c.Msg)
		} else {
			if c.LenGte > 0 && length < c.LenGte {
				path.AddFinding(fieldPath, fmt.Sprintf("length of list is less than %d", c.LenEq), c.Level, c.Msg)
			}

			if c.LenLte > 0 && length > c.LenLte {
				path.AddFinding(fieldPath, fmt.Sprintf("length of list is greater than %d", c.LenEq), c.Level, c.Msg)
			}
		}
	}
}
