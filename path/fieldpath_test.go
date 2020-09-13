package path

import (
	"testing"

	pb "github.com/bitcdr/protovalid/valid"
	"github.com/google/go-cmp/cmp"
)

func TestFindingString(t *testing.T) {
	var tests = []struct {
		input Finding
		want  string
	}{
		{Finding{
			Details: "details",
			Level:   pb.Level_LEVEL_INFO,
			Msg:     "msg",
		}, "INFO details msg"},
		{Finding{
			Details: "details",
			Level:   pb.Level_LEVEL_INFO,
			Msg:     "",
		}, "INFO details"},
	}

	for _, test := range tests {
		if got := test.input.String(); got != test.want {
			t.Errorf("got %s, want %s", got, test.want)
		}
	}
}

func TestFieldPathString(t *testing.T) {
	var tests = []struct {
		input FieldPath
		want  string
	}{
		{FieldPath{
			FullName: "name",
		}, "Full name: name, findings count: 0, childs count: 0"},
		{FieldPath{
			FullName: "name",
			Findings: []Finding{
				{
					Details: "details",
					Level:   pb.Level_LEVEL_INFO,
					Msg:     "",
				},
			},
			Childs: map[string]*FieldPath{
				"child": {
					FullName: "childName",
				},
			},
		}, "Full name: name, findings count: 1, childs count: 1"},
	}

	for _, test := range tests {
		if got := test.input.String(); got != test.want {
			t.Errorf("got %s, want %s", got, test.want)
		}
	}
}

func TestGetChild(t *testing.T) {
	fieldPath := FieldPath{
		FullName: "name",
	}

	// Add first child
	firstChild := GetChild(&fieldPath, "firstChildName")

	wantFirstChild := FieldPath{
		FullName: "name.firstChildName",
		Findings: []Finding{},
	}
	if diff := cmp.Diff(firstChild, &wantFirstChild); diff != "" {
		t.Errorf("mismatch for first child\n%s", diff)
	}

	// Add second child
	secondChild := GetChild(&fieldPath, "secondChildName")

	wantSecondChild := FieldPath{
		FullName: "name.secondChildName",
		Findings: []Finding{},
	}
	if diff := cmp.Diff(secondChild, &wantSecondChild); diff != "" {
		t.Errorf("mismatch for second child\n%s", diff)
	}

	if numChilds := len(fieldPath.Childs); numChilds != 2 {
		t.Errorf("got %d number of childs, want 2", numChilds)
	}

	// Get second child
	getSecondChild := GetChild(&fieldPath, "secondChildName")
	if diff := cmp.Diff(getSecondChild, &wantSecondChild); diff != "" {
		t.Errorf("mismatch for get second child\n%s", diff)
	}
}

func TestAddFinding(t *testing.T) {
	fieldPath := FieldPath{}

	want := Finding{
		Details: "details",
		Level:   pb.Level_LEVEL_ERROR,
		Msg:     "msg",
	}

	AddFinding(&fieldPath, want.Details, want.Level, want.Msg)

	if len(fieldPath.Findings) != 1 {
		t.Errorf("fieldPath contains not exactly one finding")
	}

	finding := fieldPath.Findings[0]

	if diff := cmp.Diff(finding, want); diff != "" {
		t.Errorf("mismatch\n%s", diff)
	}
}
