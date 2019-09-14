package dot

import "fmt"

const (
	definitionKindComponents = "components"
	definitionKindFaces      = "faces"
)

type definition struct {
	Version string `yaml:"version"`
	Kind    string `yaml:"kind"`

	// kind: components
	Components []*component `yaml:"components"`

	// kind: faces
	Faces []*face `yaml:"faces"`
}

func (def *definition) validateAndComplement() error {
	if def.Version == "" {
		return fmt.Errorf("`version` must be specified")
	}
	if def.Kind == "" {
		return fmt.Errorf("`kind` must be specified")
	}

	switch def.Kind {
	case definitionKindComponents:
		if len(def.Components) <= 0 {
			return fmt.Errorf("`components` must contain at least one content")
		}

		for _, c := range def.Components {
			err := c.validateAndComplement()
			if err != nil {
				return err
			}
		}
	case definitionKindFaces:
		if len(def.Faces) <= 0 {
			return fmt.Errorf("`faces` must contain at least one face")
		}

		for _, f := range def.Faces {
			err := f.validateAndComplement()
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("`kind` must be `components` or `faces`")
	}

	return nil
}

type component struct {
	Name         string                `yaml:"name"`
	Base         string                `yaml:"base"`
	Hide         bool                  `yaml:"hide"`
	RawLabels    interface{}           `yaml:"labels"`
	Labels       map[string][]string   `yaml:"-"`
	Dependencies []*dependentComponent `yaml:"dependencies"`
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

type face struct {
	Targets    targets           `yaml:"targets"`
	Attributes map[string]string `yaml:"attributes"`
}

func (f *face) validateAndComplement() error {
	return nil
}

type targets struct {
	MatchLabels map[string]string `yaml:"match_labels"`
}
