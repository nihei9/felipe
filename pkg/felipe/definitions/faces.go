package definitions

import "fmt"

const (
	DefinitionKindFaces = "faces"
)

type FacesDefinition struct {
	Version string  `yaml:"version"`
	Kind    string  `yaml:"kind"`
	Faces   []*Face `yaml:"faces"`
}

func (def *FacesDefinition) Validate() error {
	if def.Version == "" {
		return fmt.Errorf("`version` must be specified")
	}
	if def.Kind == "" {
		return fmt.Errorf("`kind` must be specified")
	}
	if def.Kind != DefinitionKindFaces {
		return fmt.Errorf("`kind` must be `faces`")
	}
	if len(def.Faces) <= 0 {
		return fmt.Errorf("`faces` must contain at least one face")
	}
	for _, f := range def.Faces {
		err := f.validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type Face struct {
	Targets    Targets           `yaml:"targets"`
	Attributes map[string]string `yaml:"attributes"`
}

func (f *Face) validate() error {
	return nil
}

type Targets struct {
	MatchLabels map[string]string `yaml:"match_labels"`
}
