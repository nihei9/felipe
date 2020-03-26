package component

import (
	"fmt"
)

const (
	NilComponentID = ComponentID("")
)

type ComponentID string

func (cid ComponentID) String() string {
	return string(cid)
}

func (cid ComponentID) IsNil() bool {
	return cid == NilComponentID
}

type Components struct {
	set map[ComponentID]*Component
	ids []ComponentID
}

func NewComponents() *Components {
	return &Components{
		set: map[ComponentID]*Component{},
		ids: []ComponentID{},
	}
}

func (cs *Components) Add(c *Component) bool {
	if c.ID.IsNil() {
		return false
	}
	if _, ok := cs.set[c.ID]; !ok {
		cs.ids = append(cs.ids, c.ID)
	}
	cs.set[c.ID] = c
	return true
}

func (cs *Components) Get(id ComponentID) (*Component, bool) {
	if id.IsNil() {
		return nil, false
	}
	c, ok := cs.set[id]
	if !ok {
		return newUndefinedComponent(id), false
	}
	return c, true
}

func newUndefinedComponent(id ComponentID) *Component {
	c := NewComponent(NilComponentID, id)
	c.AddLabel("__undefined__", "true")

	return c
}

func (cs *Components) GetIDs() []ComponentID {
	return cs.ids
}

func (cs *Components) Complement() error {
	for _, c := range cs.set {
		if !c.baseID.IsNil() {
			continue
		}

		err := c.complement(cs)
		if err != nil {
			return err
		}
	}

	for _, c := range cs.set {
		err := c.complement(cs)
		if err != nil {
			return err
		}
	}

	return nil
}

type Relation struct {
	Description string
}

type complementStatus string

const (
	complementStatusNew        = complementStatus("new")
	complementStatusInProgress = complementStatus("in_progress")
	complementStatusDone       = complementStatus("done")
)

type Component struct {
	ID           ComponentID
	Labels       map[string]string
	Dependencies map[ComponentID]*Relation

	baseID           ComponentID
	hidden           bool
	complementStatus complementStatus
}

func NewComponent(baseID ComponentID, id ComponentID) *Component {
	return &Component{
		ID:               id,
		Labels:           map[string]string{},
		Dependencies:     map[ComponentID]*Relation{},
		baseID:           baseID,
		hidden:           false,
		complementStatus: complementStatusNew,
	}
}

func (c *Component) AddLabel(key string, value string) {
	c.Labels[key] = value
}

func (c *Component) DependOn(dependencyID ComponentID, relation *Relation) {
	c.Dependencies[dependencyID] = relation
}

func (c *Component) IsHidden() bool {
	return c.hidden
}

func (c *Component) Hide() {
	c.hidden = true
}

func (c *Component) complement(allComponents *Components) error {
	switch c.complementStatus {
	case complementStatusDone:
		return nil
	case complementStatusInProgress:
		return fmt.Errorf("cyclic inheritance is not allowed")
	}

	c.complementStatus = complementStatusInProgress

	if !c.baseID.IsNil() {
		base, ok := allComponents.Get(c.baseID)
		if !ok {
			return fmt.Errorf("the base component `%s` is undefined", c.baseID)
		}
		err := base.complement(allComponents)
		if err != nil {
			return err
		}

		err = c.inherit(base)
		if err != nil {
			return err
		}
	}

	c.complementStatus = complementStatusDone

	return nil
}

func (c *Component) inherit(base *Component) error {
	if base == nil {
		return nil
	}

	// inherit labels from a base component
	for baseK, baseV := range base.Labels {
		if _, alreadyExists := c.Labels[baseK]; alreadyExists {
			continue
		}
		c.AddLabel(baseK, baseV)
	}

	// inherit dependencies from a base component
	for dep, rel := range base.Dependencies {
		if _, alreadyExists := c.Dependencies[dep]; alreadyExists {
			continue
		}
		c.DependOn(dep, rel)
	}

	return nil
}
