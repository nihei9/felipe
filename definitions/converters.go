package definitions

import (
	"github.com/nihei9/felipe/component"
)

func MakeComponentsDefinition(cs *component.Components) *ComponentsDefinition {
	components := []*Component{}
	for _, id := range cs.GetIDs() {
		c, _ := cs.Get(id)
		deps := []*DependentComponent{}
		for dep, rel := range c.Dependencies {
			deps = append(deps, &DependentComponent{
				ID:       dep.String(),
				Relation: rel.Description,
			})
		}

		components = append(components, &Component{
			ID:           c.ID.String(),
			Hide:         false,
			Labels:       c.Labels,
			Dependencies: deps,
		})
	}

	return &ComponentsDefinition{
		Version:    "1",
		Kind:       DefinitionKindComponents,
		Components: components,
	}
}

func MakeComponentEntity(def *Component) *component.Component {
	baseID := component.ComponentID(def.Base)
	id := component.ComponentID(def.ID)
	c := component.NewComponent(baseID, id)
	for k, v := range def.Labels {
		c.AddLabel(k, v)
	}
	for _, dDef := range def.Dependencies {
		rel := &component.Relation{
			Description: dDef.Relation,
		}
		c.DependOn(component.ComponentID(dDef.ID), rel)
	}
	if def.Hide {
		c.Hide()
	}

	return c
}
