package definitions

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

const (
	DefinitionKindComponents = "components"
)

func ReadComponentsDefinition(r io.Reader) (*ComponentsDefinition, error) {
	def := &ComponentsDefinition{}
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

type ComponentsDefinition struct {
	Version    string       `yaml:"version"`
	Kind       string       `yaml:"kind"`
	Components []*Component `yaml:"components"`
}

func (def *ComponentsDefinition) validate() error {
	if def.Version == "" {
		return errorVersionIsMissing
	}
	if def.Kind == "" {
		return errorKindIsMissing
	}
	if def.Kind != DefinitionKindComponents {
		return errorKindIsNotComponents
	}
	if len(def.Components) <= 0 {
		return errorComponentsHasNoComponent
	}
	for _, c := range def.Components {
		if c == nil {
			return errorComponentsHasEmptyComponent
		}

		err := c.validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type Component struct {
	ID           string                `yaml:"id"`
	Base         string                `yaml:"base,omitempty"`
	Hide         bool                  `yaml:"hide,omitempty"`
	Labels       map[string]string     `yaml:"labels,omitempty"`
	Dependencies []*DependentComponent `yaml:"dependencies,omitempty"`
}

func (c *Component) UnmarshalYAML(unmarshal func(interface{}) error) error {
	comp := &struct {
		ID           string                `yaml:"id"`
		Base         string                `yaml:"base"`
		Hide         bool                  `yaml:"hide"`
		Labels       interface{}           `yaml:"labels"`
		Dependencies []*DependentComponent `yaml:"dependencies"`
	}{}
	err := unmarshal(comp)
	if err != nil {
		return err
	}

	labels := map[string]string{}
	if comp.Labels != nil {
		rawLabels, ok := comp.Labels.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("`labels` must be map[string]string or map[string][]string")
		}
		for rawKey, rawValue := range rawLabels {
			key, ok := rawKey.(string)
			if !ok {
				return fmt.Errorf("a key of `labels` must be string")
			}
			value, ok := rawValue.(string)
			if !ok {
				return fmt.Errorf("a value of `labels` must be string or []string")
			}
			labels[key] = value
		}
	}

	c.ID = comp.ID
	c.Base = comp.Base
	c.Hide = comp.Hide
	c.Labels = labels
	c.Dependencies = comp.Dependencies

	return nil
}

func (c *Component) validate() error {
	if c.ID == "" {
		return errorComponentIDIsMissing
	}
	for _, dc := range c.Dependencies {
		if dc == nil {
			return errorComponentHasEmptyDependency
		}

		err := dc.validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type DependentComponent struct {
	ID       string `yaml:"id"`
	Relation string `yaml:"relation"`
}

func (dc *DependentComponent) validate() error {
	if dc.ID == "" {
		return errorDependencyIDIsMissing
	}

	return nil
}
