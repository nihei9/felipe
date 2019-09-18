package definitions

import (
	"fmt"
)

const (
	DefinitionKindComponents = "components"
)

type ComponentsDefinition struct {
	Version    string       `yaml:"version"`
	Kind       string       `yaml:"kind"`
	Components []*Component `yaml:"components"`
}

func (def *ComponentsDefinition) Validate() error {
	if def.Version == "" {
		return fmt.Errorf("`version` must be specified")
	}
	if def.Kind == "" {
		return fmt.Errorf("`kind` must be specified")
	}
	if def.Kind != DefinitionKindComponents {
		return fmt.Errorf("`kind` must be `components`")
	}
	if len(def.Components) <= 0 {
		return fmt.Errorf("`components` must contain at least one content")
	}
	for _, c := range def.Components {
		err := c.validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type Component struct {
	Name         string                `yaml:"name"`
	Base         string                `yaml:"base,omitempty"`
	Hide         bool                  `yaml:"hide,omitempty"`
	Labels       map[string][]string   `yaml:"labels,omitempty"`
	Dependencies []*DependentComponent `yaml:"dependencies,omitempty"`
}

func (c *Component) UnmarshalYAML(unmarshal func(interface{}) error) error {
	comp := &struct {
		Name         string                `yaml:"name"`
		Base         string                `yaml:"base"`
		Hide         bool                  `yaml:"hide"`
		Labels       interface{}           `yaml:"labels"`
		Dependencies []*DependentComponent `yaml:"dependencies"`
	}{}
	err := unmarshal(comp)
	if err != nil {
		return err
	}

	labels := map[string][]string{}
	if comp.Labels != nil {
		rawLabels, ok := comp.Labels.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("`labels` must be map[string]string or map[string][]string")
		}
		for rawKey, rawValues := range rawLabels {
			key, ok := rawKey.(string)
			if !ok {
				return fmt.Errorf("a key of `labels` must be string")
			}
			switch values := rawValues.(type) {
			case string:
				labels[key] = []string{values}
			case []interface{}:
				s := []string{}
				for _, value := range values {
					v, ok := value.(string)
					if !ok {
						return fmt.Errorf("a value of `labels` must be string or []string")
					}
					s = append(s, v)
				}
				labels[key] = s
			default:
				return fmt.Errorf("a value of `labels` must be string or []string")
			}
		}
	}

	c.Name = comp.Name
	c.Base = comp.Base
	c.Hide = comp.Hide
	c.Labels = labels
	c.Dependencies = comp.Dependencies

	return nil
}

func (c *Component) validate() error {
	if c.Name == "" {
		return fmt.Errorf("`componets[].name` must be specified")
	}
	for _, dc := range c.Dependencies {
		err := dc.validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type DependentComponent struct {
	Name string `yaml:"name"`
}

func (dc *DependentComponent) validate() error {
	if dc.Name == "" {
		return fmt.Errorf("`dependencies[].name` must be specified")
	}

	return nil
}
