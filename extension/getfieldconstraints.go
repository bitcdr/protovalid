// Package extension handles field extensions.
package extension

import (
	pb "github.com/bitcdr/protovalid/valid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func getFieldConstraints(fd protoreflect.FieldDescriptor) (bool, *pb.FieldConstraints) {
	opts := fd.Options().(*descriptorpb.FieldOptions)

	if proto.HasExtension(opts, pb.E_Field) {
		return true, proto.GetExtension(opts, pb.E_Field).(*pb.FieldConstraints)
	}

	return false, nil
}

func GetFieldConstraintsDouble(fd protoreflect.FieldDescriptor) (bool, *pb.FieldConstraints_DoubleConstraints) {
	if ok, xt := getFieldConstraints(fd); ok {
		if res := xt.GetDouble(); res != nil {
			return true, res
		}
	}

	return false, nil
}

func GetFieldConstraintsInt32(fd protoreflect.FieldDescriptor) (bool, *pb.FieldConstraints_Int32Constraints) {
	if ok, xt := getFieldConstraints(fd); ok {
		if res := xt.GetInt32(); res != nil {
			return true, res
		}
	}

	return false, nil
}

func GetFieldConstraintsString(fd protoreflect.FieldDescriptor) (bool, *pb.FieldConstraints_StringConstraints) {
	if ok, xt := getFieldConstraints(fd); ok {
		if res := xt.GetString_(); res != nil {
			return true, res
		}
	}

	return false, nil
}

func GetFieldConstraintsRepeated(fd protoreflect.FieldDescriptor) (bool, *pb.FieldConstraints_RepeatedConstraints) {
	if ok, xt := getFieldConstraints(fd); ok {
		if res := xt.GetRepeated(); res != nil {
			return true, res
		}
	}

	return false, nil
}
