package definitions

import (
	"strings"
	"testing"
)

func TestFacesDefinition(t *testing.T) {
	tests := []struct {
		caption string
		data    string
		err     error
	}{
		{
			caption: "`faces` has a face",
			data: `
version: 1
kind: faces
faces:
- targets:
    match_labels:
      l1: foo
  attributes:
    fontcolor: red
`,
		},
		{
			caption: "`faces` has some faces",
			data: `
version: 1
kind: faces
faces:
- targets:
    match_labels:
      l1: foo
  attributes:
    fontcolor: white
    fillcolor: black
    style: filled
`,
		},
		{
			caption: "`version` is not specified",
			data: `
kind: faces
faces:
- targes:
    match_labels:
      l1: foo
  attributes:
    fontcolor: red
`,
			err: errorVersionIsMissing,
		},
		{
			caption: "`version` is emnpty",
			data: `
version:
kind: faces
faces:
- targes:
    match_labels:
      l1: foo
  attributes:
    fontcolor: red
`,
			err: errorVersionIsMissing,
		},
		{
			caption: "`kind` is not specified",
			data: `
version: 1
faces:
- targes:
    match_labels:
      l1: foo
  attributes:
    fontcolor: red
`,
			err: errorKindIsMissing,
		},
		{
			caption: "`kind` is empty",
			data: `
version: 1
kind:
faces:
- targes:
    match_labels:
      l1: foo
  attributes:
    fontcolor: red
`,
			err: errorKindIsMissing,
		},
		{
			caption: "`kind` is not `faces`",
			data: `
version: 1
kind: foo
faces:
- targes:
    match_labels:
      l1: foo
  attributes:
    fontcolor: red
`,
			err: errorKindIsNotFaces,
		},
		{
			caption: "`faces` has no face",
			data: `
version: 1
kind: faces
`,
			err: errorFacesHasNoFace,
		},
		{
			caption: "`faces` has no face",
			data: `
version: 1
kind: faces
faces:
`,
			err: errorFacesHasNoFace,
		},
		{
			caption: "`faces[]` includes an empty face",
			data: `
version: 1
kind: faces
faces:
-
`,
			err: errorFacesHasEmptyFace,
		},
		{
			caption: "`faces[]` includes an empty face",
			data: `
version: 1
kind: faces
faces:
- targets:
    match_labels:
      l1: foo
  attributes:
    fontcolor: red
-
- targets:
    match_labels:
      l2: bar
  attributes:
    shape: cylinder
`,
			err: errorFacesHasEmptyFace,
		},
		{
			caption: "`faces[].targets` is not specified",
			data: `
version: 1
kind: faces
faces:
- attributes:
    fontcolor: red
`,
			err: errorFaceTargetIsMissing,
		},
		{
			caption: "`faces[].targets` is empty",
			data: `
version: 1
kind: faces
faces:
- targets:
  attributes:
    fontcolor: red
`,
			err: errorFaceTargetIsMissing,
		},
		{
			caption: "`faces[].targets` is empty",
			data: `
version: 1
kind: faces
faces:
- targets:
    match_labels:
  attributes:
    fontcolor: red
`,
			err: errorFaceTargetIsEmpty,
		},
		{
			caption: "`faces[].targets.match_labels[]` includes an empty entry",
			data: `
version: 1
kind: faces
faces:
- targets:
    match_labels:
      "": foo
  attributes:
    fontcolor: red
`,
			err: errorFaceMatchLabelsTargetHasEmptyEntry,
		},
		{
			caption: "`faces[].targets.match_labels[]` includes an empty entry",
			data: `
version: 1
kind: faces
faces:
- targets:
    match_labels:
      l1: foo
      "": bar
      l2: baz
  attributes:
    fontcolor: red
`,
			err: errorFaceMatchLabelsTargetHasEmptyEntry,
		},
		{
			caption: "`faces[].attributes[]` is not specified",
			data: `
version: 1
kind: faces
faces:
- targets:
    match_labels:
      l1: foo
`,
			err: errorFaceAttributesHasNoAttribute,
		},
		{
			caption: "`faces[].attributes[]` has no attribute",
			data: `
version: 1
kind: faces
faces:
- targets:
    match_labels:
      l1: foo
  attributes:
`,
			err: errorFaceAttributesHasNoAttribute,
		},
		{
			caption: "`faces[].attributes[]` includes an empty face",
			data: `
version: 1
kind: faces
faces:
- targets:
    match_labels:
      l1: foo
  attributes:
    "":
`,
			err: errorFaceAttributesHasEmptyAttribute,
		},
		{
			caption: "`faces[].attributes[]` includes an empty face",
			data: `
version: 1
kind: faces
faces:
- targets:
    match_labels:
      l1: foo
  attributes:
    fillcolor: black
    "":
    fontcolor: white
`,
			err: errorFaceAttributesHasEmptyAttribute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.caption, func(t *testing.T) {
			_, err := ReadFacesDefinition(strings.NewReader(tt.data))
			if err != tt.err {
				t.Error(err)
			}
		})
	}
}
