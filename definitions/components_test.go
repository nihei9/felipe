package definitions

import (
	"strings"
	"testing"
)

func TestComponentsDefinition(t *testing.T) {
	tests := []struct {
		caption string
		data    string
		err     error
	}{
		{
			caption: "`components` has a simple component",
			data: `
version: 1
kind: components
components:
- id: c1
`,
		},
		{
			caption: "`components` has a simple components",
			data: `
version: 1
kind: components
components:
- id: c1
  hide:
  base:
  dependencies:
`,
		},
		{
			caption: "`components` has two simple components",
			data: `
version: 1
kind: components
components:
- id: c1
- id: c2
`,
		},
		{
			caption: "`components` has two components",
			data: `
version: 1
kind: components
components:
- id: c1
  hide: true
  base:
  dependencies:
  - id: d1
- id: c2
  hide: false
  base: c1
  dependencies:
  - id: d1
  - id: d2
`,
		},
		{
			caption: "`version` is not specified",
			data: `
kind: components
components:
- id: c1
`,
			err: errorVersionIsMissing,
		},
		{
			caption: "`version` is emnpty",
			data: `
version:
kind: components
components:
- id: c1
`,
			err: errorVersionIsMissing,
		},
		{
			caption: "`kind` is not specified",
			data: `
version: 1
components:
- id: c1
`,
			err: errorKindIsMissing,
		},
		{
			caption: "`kind` is empty",
			data: `
version: 1
kind:
components:
- id: c1
`,
			err: errorKindIsMissing,
		},
		{
			caption: "`kind` is not `components`",
			data: `
version: 1
kind: foo
components:
- id: c1
`,
			err: errorKindIsNotComponents,
		},
		{
			caption: "`components` has no component",
			data: `
version: 1
kind: components
`,
			err: errorComponentsHasNoComponent,
		},
		{
			caption: "`components` has no component",
			data: `
version: 1
kind: components
components:
`,
			err: errorComponentsHasNoComponent,
		},
		{
			caption: "`components[]` includes an empty component",
			data: `
version: 1
kind: components
components:
-
`,
			err: errorComponentsHasEmptyComponent,
		},
		{
			caption: "`components[]` includes an empty component",
			data: `
version: 1
kind: components
components:
- id: c1
-
- id: c2
`,
			err: errorComponentsHasEmptyComponent,
		},
		{
			caption: "`components[].id` is not specified",
			data: `
version: 1
kind: components
components:
- id:
`,
			err: errorComponentIDIsMissing,
		},
		{
			caption: "`components[].dependencies[]` includes an empty dependency",
			data: `
version: 1
kind: components
components:
- id: c1
  dependencies:
  -
`,
			err: errorComponentHasEmptyDependency,
		},
		{
			caption: "`components[].dependencies[]` includes an empty dependency",
			data: `
version: 1
kind: components
components:
- id: c1
  dependencies:
  - id: d1
  -
  - id: d2
`,
			err: errorComponentHasEmptyDependency,
		},
		{
			caption: "`components[].dependencies[].id` is not specified",
			data: `
version: 1
kind: components
components:
- id: c1
  dependencies:
  - id:
`,
			err: errorDependencyIDIsMissing,
		},
	}
	for _, tt := range tests {
		t.Run(tt.caption, func(t *testing.T) {
			_, err := ReadComponentsDefinition(strings.NewReader(tt.data))
			if err != tt.err {
				t.Error(err)
			}
		})
	}
}
