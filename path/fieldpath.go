// Package path contains methods for navigating protobuf paths.
package path

import (
	"fmt"
	"strings"
)

// FieldPath is a nested tree of a field path names.
type FieldPath struct {
	FullName string
	Findings []string
	Childs   map[string]*FieldPath
}

func (fp FieldPath) String() string {
	return fmt.Sprintf("Full name: %s, findings count: %d, childs count: %d", fp.FullName, len(fp.Findings), len(fp.Childs))
}

func GetChild(fieldPath *FieldPath, name string) *FieldPath {
	if child, ok := fieldPath.Childs[name]; ok {
		return child
	}

	return addChild(fieldPath, name)
}

func childFullName(fieldPath *FieldPath, name string) string {
	if len(fieldPath.FullName) > 0 {
		return strings.Join([]string{fieldPath.FullName, name}, ".")
	}

	return fmt.Sprintf("%s", name)
}

func addChild(fieldPath *FieldPath, name string) *FieldPath {
	if fieldPath.Childs == nil {
		fieldPath.Childs = make(map[string]*FieldPath)
	}

	child := FieldPath{
		FullName: childFullName(fieldPath, name),
		Findings: []string{},
		Childs:   nil,
	}

	fieldPath.Childs[name] = &child

	return &child
}

func AddFinding(fieldPath *FieldPath, finding string) {
	fieldPath.Findings = append(fieldPath.Findings, finding)
}
