// Package path contains methods for navigating protobuf paths.
package path

import (
	"fmt"
	"strings"

	pb "github.com/bitcdr/protovalid/valid"
)

// Finding contains the details, level, and custom message.
type Finding struct {
	Details string
	Level   pb.Level
	Msg     string
}

func (f Finding) String() string {
	if len(f.Msg) == 0 {
		return fmt.Sprintf("%s %s", strings.TrimPrefix(f.Level.String(), "LEVEL_"), f.Details)
	}

	return fmt.Sprintf("%s %s %s", strings.TrimPrefix(f.Level.String(), "LEVEL_"), f.Details, f.Msg)
}

// FieldPath is a nested tree of a field path names.
type FieldPath struct {
	FullName string
	Findings []Finding
	Childs   map[string]*FieldPath
}

func (fp FieldPath) String() string {
	return fmt.Sprintf("Full name: %s, findings count: %d, childs count: %d", fp.FullName, len(fp.Findings), len(fp.Childs))
}

// GetChild gets the child with the name for the field path and creates it if it doesn't exist.
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
		Findings: []Finding{},
		Childs:   nil,
	}

	fieldPath.Childs[name] = &child

	return &child
}

// AddFinding adds the finding with the details, level, and message to the field path.
func AddFinding(fieldPath *FieldPath, details string, level pb.Level, msg string) {
	finding := Finding{
		Details: details,
		Level:   level,
		Msg:     msg,
	}

	fieldPath.Findings = append(fieldPath.Findings, finding)
}
