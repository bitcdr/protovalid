// Package wrappers provides access to protobuf wrapped types.
package wrappers

import "google.golang.org/protobuf/reflect/protoreflect"

// GetFloat64 gets the double value of a google.protobuf.DoubleValue message type.
func GetFloat64(fd protoreflect.FieldDescriptor, v protoreflect.Value) float64 {
	var value float64

	switch fd.Message().FullName() {
	case "google.protobuf.DoubleValue":
		v.Message().Range(func(wrappedFd protoreflect.FieldDescriptor, wrappedValue protoreflect.Value) bool {
			value = wrappedValue.Float()
			return false
		})
	default:
		panic("Unsupported message type, require DoubleValue")
	}

	return value
}

// GetInt32 gets the int32 value of a google.protobuf.Int32Value message type.
func GetInt32(fd protoreflect.FieldDescriptor, v protoreflect.Value) int32 {
	var value int32

	switch fd.Message().FullName() {
	case "google.protobuf.Int32Value":
		v.Message().Range(func(wrappedFd protoreflect.FieldDescriptor, wrappedValue protoreflect.Value) bool {
			value = int32(wrappedValue.Int())
			return false
		})
	default:
		panic("Unsupported message type, require Int32Value")
	}

	return value
}

// GetInt64 gets the int64 value of a google.protobuf.Int64Value message type.
func GetInt64(fd protoreflect.FieldDescriptor, v protoreflect.Value) int64 {
	var value int64

	switch fd.Message().FullName() {
	case "google.protobuf.Int64Value":
		v.Message().Range(func(wrappedFd protoreflect.FieldDescriptor, wrappedValue protoreflect.Value) bool {
			value = int64(wrappedValue.Int())
			return false
		})
	default:
		panic("Unsupported message type, require Int64Value")
	}

	return value
}

// GetString gets the string value of a google.protobuf.StringValue message type.
func GetString(fd protoreflect.FieldDescriptor, v protoreflect.Value) string {
	var value string

	switch fd.Message().FullName() {
	case "google.protobuf.StringValue":
		v.Message().Range(func(wrappedFd protoreflect.FieldDescriptor, wrappedValue protoreflect.Value) bool {
			value = wrappedValue.String()
			return false
		})
	default:
		panic("Unsupported message type, require StringValue")
	}

	return value
}
