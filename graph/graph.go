package graph

import "fmt"

type Components struct {
	Components map[ComponentID]*Component
}

func NewComponents() *Components {
	return &Components{
		Components: map[ComponentID]*Component{},
	}
}

func (cs *Components) Get(cid ComponentID) (*Component, bool) {
	c, ok := cs.Components[cid]
	if !ok {
		return newComponent(cid, NewComponentID(""), true), false
	}
	return c, true
}

func (cs *Components) Add(c *Component) error {
	if _, ok := cs.Components[c.ID()]; ok {
		return fmt.Errorf("component `%v` already exists", c.ID)
	}
	cs.Components[c.ID()] = c

	return nil
}

func (cs *Components) Complement() error {
	for _, c := range cs.Components {
		err := complement(c, cs, []ComponentID{})
		if err != nil {
			return err
		}
	}

	return nil
}

func complement(c *Component, cs *Components, stack []ComponentID) error {
	if c.complemented || c.base.Nil() {
		return nil
	}
	for _, cid := range stack {
		if cid == c.ID() {
			return fmt.Errorf("cyclic inheritance is not allowed\n%v", stack)
		}
	}
	base, ok := cs.Get(c.base)
	if !ok {
		return fmt.Errorf("it is not allowed to inherit the undefined component; %v", c.base)
	}
	complement(base, cs, append(stack, c.ID()))
	for k, vs := range base.Labels() {
		for _, v := range vs {
			c.Label(k, v)
		}
	}
	for _, v := range base.Dependencies() {
		c.DependOn(v)
	}
	c.complemented = true

	return nil
}

type ComponentID string

func (cid ComponentID) String() string {
	return string(cid)
}

func (cid ComponentID) Nil() bool {
	if cid.String() == "" {
		return true
	}
	return false
}

func NewComponentID(name string) ComponentID {
	return ComponentID(name)
}

type Component struct {
	id           ComponentID
	base         ComponentID
	queryable    bool
	labels       map[string][]string
	dependencies []ComponentID
	complemented bool
}

func NewComponent(name string, base string, queryable bool) *Component {
	return newComponent(NewComponentID(name), NewComponentID(base), queryable)
}

func newComponent(cid ComponentID, base ComponentID, queryable bool) *Component {
	return &Component{
		id:           cid,
		base:         base,
		queryable:    queryable,
		labels:       map[string][]string{},
		dependencies: []ComponentID{},
		complemented: false,
	}
}

func (c *Component) ID() ComponentID {
	return c.id
}

func (c *Component) Labels() map[string][]string {
	return c.labels
}

func (c *Component) Label(key string, value string) {
	if _, ok := c.labels[key]; !ok {
		c.labels[key] = []string{value}
	} else {
		c.labels[key] = append(c.labels[key], value)
	}
}

func (c *Component) Dependencies() []ComponentID {
	return c.dependencies
}

func (c *Component) DependOn(cid ComponentID) {
	c.dependencies = append(c.dependencies, cid)
}

type Face struct {
	Condition  *Condition
	Attributes map[string]string
}

func NewFace() *Face {
	return &Face{
		Condition:  NewCondition(),
		Attributes: map[string]string{},
	}
}

func (f *Face) AddTarget(m Matcher) {
	f.Condition.AddMatcher(m)
}

func (f *Face) AddAttributes(attrs map[string]string) {
	for k, v := range attrs {
		f.Attributes[k] = v
	}
}
