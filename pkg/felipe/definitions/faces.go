package definitions

import (
	"io"

	"gopkg.in/yaml.v2"
)

const (
	DefinitionKindFaces = "faces"
)

func ReadFacesDefinition(r io.Reader) (*FacesDefinition, error) {
	def := &FacesDefinition{}
	err := yaml.NewDecoder(r).Decode(def)
	if err != nil {
		return nil, err
	}

	err = def.validate()
	if err != nil {
		return nil, err
	}

	return def, nil
}

type FacesDefinition struct {
	Version string  `yaml:"version"`
	Kind    string  `yaml:"kind"`
	Faces   []*Face `yaml:"faces"`
}

func (def *FacesDefinition) validate() error {
	if def.Version == "" {
		return errorVersionIsMissing
	}
	if def.Kind == "" {
		return errorKindIsMissing
	}
	if def.Kind != DefinitionKindFaces {
		return errorKindIsNotFaces
	}
	if len(def.Faces) <= 0 {
		return errorFacesHasNoFace
	}
	for _, f := range def.Faces {
		if f == nil {
			return errorFacesHasEmptyFace
		}

		err := f.validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type Face struct {
	Targets    *Targets          `yaml:"targets"`
	Attributes map[string]string `yaml:"attributes"`
}

func (f *Face) validate() error {
	if f.Targets == nil {
		return errorFaceTargetIsMissing
	}
	err := f.Targets.validate()
	if err != nil {
		return err
	}
	if len(f.Attributes) <= 0 {
		return errorFaceAttributesHasNoAttribute
	}
	for k := range f.Attributes {
		if k == "" {
			return errorFaceAttributesHasEmptyAttribute
		}
	}

	return nil
}

type Targets struct {
	MatchLabels map[string]string `yaml:"match_labels"`
}

func (t *Targets) validate() error {
	if len(t.MatchLabels) <= 0 {
		return errorFaceTargetIsEmpty
	}
	for k := range t.MatchLabels {
		if k == "" {
			return errorFaceMatchLabelsTargetHasEmptyEntry
		}
	}

	return nil
}
