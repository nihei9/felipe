package definitions

import "github.com/nihei9/felipe/component"

func MakeComponentEntity(def *Component) *component.Component {
	baseID := component.ComponentID(def.Base)
	id := component.ComponentID(def.ID)
	c := component.NewComponent(baseID, id)
	for k, v := range def.Labels {
		c.AddLabel(k, v)
	}
	for _, dDef := range def.Dependencies {
		c.DependOn(component.ComponentID(dDef.ID))
	}
	if def.Hide {
		c.Hide()
	}

	return c
}
