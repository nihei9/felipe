package query

import "fmt"

const (
	definitionKindComponents = "components"
)

type definition struct {
	Version    string       `yaml:"version"`
	Kind       string       `yaml:"kind"`
	Components []*component `yaml:"components"`
}

func (def *definition) validateAndComplement() error {
	if def.Version == "" {
		return fmt.Errorf("`version` must be specified")
	}
	if def.Kind == "" {
		return fmt.Errorf("`kind` must be specified")
	}
	if def.Kind != definitionKindComponents {
		return fmt.Errorf("`kind` must be `components`")
	}
	if len(def.Components) <= 0 {
		return fmt.Errorf("`components` must contain at least one content")
	}
	for _, c := range def.Components {
		err := c.validateAndComplement()
		if err != nil {
			return err
		}
	}

	return nil
}

type component struct {
	Name         string                `yaml:"name"`
	Base         string                `yaml:"base,omitempty"`
	Hide         bool                  `yaml:"hide,omitempty"`
	RawLabels    interface{}           `yaml:"labels,omitempty"`
	Labels       map[string][]string   `yaml:"-"`
	Dependencies []*dependentComponent `yaml:"dependencies,omitempty"`
}

func (c *component) validateAndComplement() error {
	c.Labels = map[string][]string{}

	if c.Name == "" {
		return fmt.Errorf("`componets[].name` must be specified")
	}

	if c.RawLabels != nil {
		rawLabels, ok := c.RawLabels.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("`labels` is malformed")
		}
		for rawKey, rawValues := range rawLabels {
			key, ok := rawKey.(string)
			if !ok {
				return fmt.Errorf("a key of `labels` must be string")
			}
			switch values := rawValues.(type) {
			case string:
				c.Labels[key] = []string{values}
			case []interface{}:
				s := []string{}
				for _, value := range values {
					v, ok := value.(string)
					if !ok {
						return fmt.Errorf("a value of `labels` must be string")
					}
					s = append(s, v)
				}
				c.Labels[key] = s
			default:
				return fmt.Errorf("`labels` is malformed")
			}
		}
	}

	for _, dc := range c.Dependencies {
		err := dc.validateAndComplement()
		if err != nil {
			return err
		}
	}

	return nil
}

type dependentComponent struct {
	Name string `yaml:"name"`
}

func (dc *dependentComponent) validateAndComplement() error {
	if dc.Name == "" {
		return fmt.Errorf("`dependencies[].name` must be specified")
	}

	return nil
}
