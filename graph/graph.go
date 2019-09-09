package graph

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
	return c, ok
}

func (cs *Components) AddComponent(c *Component) {
	cs.Components[c.ID] = c
}

type ComponentID string

func (cid ComponentID) String() string {
	return string(cid)
}

func NewComponentID(name string) ComponentID {
	return ComponentID(name)
}

type Component struct {
	ID           ComponentID
	Name         string
	Labels       map[string][]string
	Dependencies []ComponentID
}

func NewComponent(name string) *Component {
	return &Component{
		ID:           NewComponentID(name),
		Name:         name,
		Labels:       map[string][]string{},
		Dependencies: []ComponentID{},
	}
}

func (c *Component) Label(key string, value string) {
	if _, ok := c.Labels[key]; !ok {
		c.Labels[key] = []string{value}
	} else {
		c.Labels[key] = append(c.Labels[key], value)
	}
}

func (c *Component) DependOn(cid ComponentID) {
	c.Dependencies = append(c.Dependencies, cid)
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
