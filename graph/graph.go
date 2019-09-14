package graph

import (
	"fmt"
	"regexp"
	"strings"
)

type Components struct {
	components map[ComponentID]*Component
	tmpList    []*Component
}

func NewComponents() *Components {
	return &Components{
		components: map[ComponentID]*Component{},
		tmpList:    []*Component{},
	}
}

func (cs *Components) Components() map[ComponentID]*Component {
	return cs.components
}

func (cs *Components) Get(cid ComponentID) (*Component, bool) {
	c, ok := cs.components[cid]
	if !ok {
		return newComponent(cid, NewComponentID(""), true), false
	}
	return c, true
}

func (cs *Components) Add(c *Component) {
	cs.tmpList = append(cs.tmpList, c)
}

func (cs *Components) add(c *Component) error {
	if _, ok := cs.components[c.ID()]; ok {
		return fmt.Errorf("component `%v` already exists", c.ID)
	}
	cs.components[c.ID()] = c

	return nil
}

func (cs *Components) Complement() error {
	for _, c := range cs.tmpList {
		re := regexp.MustCompile(`\{.*\}`)
		matches := re.FindSubmatch([]byte(c.ID().String()))
		embeddedValues := []string{}
		for _, m := range matches {
			s := string(m)
			labelName := strings.TrimSpace(s[1 : len(s)-1])
			vs, ok := c.Labels()[labelName]
			if !ok {
				return fmt.Errorf("undefined label `%s` cannot use in `name` directive", labelName)
			}
			if len(vs) != 1 {
				return fmt.Errorf("a label used as the embeded label must have just one value; `%s` has %v values", labelName, len(vs))
			}
			embeddedValues = append(embeddedValues, vs[0])
		}
		id := c.ID().String()
		for i, m := range matches {
			id = strings.Replace(id, string(m), embeddedValues[i], 1)
		}
		c.id = NewComponentID(id)

		err := cs.add(c)
		if err != nil {
			return err
		}
	}
	for _, c := range cs.components {
		err := complement(c, cs, []ComponentID{})
		if err != nil {
			return err
		}
	}

	return nil
}

func complement(c *Component, cs *Components, stack []ComponentID) error {
	if c.complemented {
		return nil
	}

	if c.base.Nil() {
		c.complemented = true

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
