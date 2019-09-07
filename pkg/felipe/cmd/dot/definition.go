package dot

import "fmt"

const (
	definitionKindComponents = "components"
)

type definition struct {
	Version    string       `yaml:"version"`
	Kind       string       `yaml:"kind"`
	Components []*component `yaml:"components"`
}

func (def *definition) validate() error {
	if def.Version == "" {
		return fmt.Errorf("`version` must be specified")
	}
	if def.Kind == "" {
		return fmt.Errorf("`kind` must be specified")
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

type component struct {
	Name         string                `yaml:"name"`
	Labels       map[string]string     `yaml:"labels"`
	Dependencies []*dependentComponent `yaml:"dependencies"`
}

func (c *component) validate() error {
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

type dependentComponent struct {
	Name string `yaml:"name"`
}

func (dc *dependentComponent) validate() error {
	if dc.Name == "" {
		return fmt.Errorf("`dependencies[].name` must be specified")
	}

	return nil
}