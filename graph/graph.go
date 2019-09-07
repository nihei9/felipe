package graph

import (
	"fmt"
	"strings"
)

type Components struct {
	Components map[ComponentID]*Component
}

func NewComponents() *Components {
	return &Components{
		Components: map[ComponentID]*Component{},
	}
}

func (cs *Components) AddComponent(c *Component) {
	cs.Components[c.ID] = c
}

func (cs *Components) Query(label string) ([]*Component, error) {
	if label == "" {
		allCs := []*Component{}
		for _, c := range cs.Components {
			allCs = append(allCs, c)
		}

		return allCs, nil
	}

	l := strings.Split(label, "=")
	if len(l) != 2 {
		return nil, fmt.Errorf("query label is malformed; got: %v", label)
	}
	queryK := strings.TrimSpace(l[0])
	queryV := strings.TrimSpace(l[1])

	matches := []*Component{}
	for _, c := range cs.Components {
		v, ok := c.Labels[queryK]
		if ok && v == queryV {
			matches = append(matches, c)
		}
	}

	return matches, nil
}

type ComponentID string

func (cid ComponentID) String() string {
	return string(cid)
}

func newComponentID(name string) ComponentID {
	return ComponentID(name)
}

type Component struct {
	ID           ComponentID
	Name         string
	Labels       map[string]string
	Dependencies map[ComponentID]*Component
}

func NewComponent(name string) *Component {
	return &Component{
		ID:           newComponentID(name),
		Name:         name,
		Labels:       map[string]string{},
		Dependencies: map[ComponentID]*Component{},
	}
}

func (c *Component) Label(key string, value string) {
	c.Labels[key] = value
}

func (c *Component) DependOn(d *Component) {
	c.Dependencies[d.ID] = d
}
