package graph

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
	Dependencies map[ComponentID]*Component
}

func NewComponent(name string) *Component {
	return &Component{
		ID:           newComponentID(name),
		Name:         name,
		Dependencies: map[ComponentID]*Component{},
	}
}

func (c *Component) DependOn(d *Component) {
	c.Dependencies[d.ID] = d
}
