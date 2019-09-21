package definitions

import (
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

const (
	DefinitionKindComponents = "components"
)

var (
	errorVersionIsMissing            = errors.New("`version` must be specified")
	errorKindIsMissing               = errors.New("`kind` must be specified")
	errorKindIsNotComponents         = errors.New("`kind` must be `components`")
	errorComponentsHasNoComponent    = errors.New("`components` must contain at least one content")
	errorComponentsHasEmptyComponent = errors.New("`components[]` includes empty components")
	errorComponentNameIsMissing      = errors.New("`components[].name` must be specified")
	errorComponentHasEmptyDependency = errors.New("`dependencies[]` includes empty components")
	errorDependencyNameIsMissing     = errors.New("`dependencies[].name` must be specified")
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
		return errorComponentNameIsMissing
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
	Name string `yaml:"name"`
}

func (dc *DependentComponent) validate() error {
	if dc.Name == "" {
		return errorDependencyNameIsMissing
	}

	return nil
}
