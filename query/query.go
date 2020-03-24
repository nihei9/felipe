package query

import (
	"github.com/nihei9/felipe/component"
)

type Query struct {
	Components   *component.Components
	Filter       Filter
	Complementer Complementer
}

func (q Query) Do() (*component.Components, error) {
	filteredComponents, err := q.Filter.Filter(q.Components)
	if err != nil {
		return nil, err
	}

	if q.Complementer != nil {
		result, err := q.Complementer.Complement(filteredComponents)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return filteredComponents, nil
}

type Passer interface {
	Pass(target *component.Component) (bool, error)
}

type Filter interface {
	Passer
	Filter(target *component.Components) (*component.Components, error)
}

func filter(p Passer, target *component.Components) (*component.Components, error) {
	result := component.NewComponents()
	for _, id := range target.GetIDs() {
		c, _ := target.Get(id)
		pass, err := p.Pass(c)
		if err != nil {
			return nil, err
		}
		if pass {
			result.Add(c)
		}
	}

	return result, nil
}

type AllPassFilter struct {
}

func (f AllPassFilter) Filter(target *component.Components) (*component.Components, error) {
	return filter(f, target)
}

func (f AllPassFilter) Pass(target *component.Component) (bool, error) {
	if target.IsHidden() {
		return false, nil
	}
	return true, nil
}

type LabelsFilter struct {
	Labels map[string]string
}

func (f LabelsFilter) Filter(target *component.Components) (*component.Components, error) {
	return filter(f, target)
}

func (f LabelsFilter) Pass(target *component.Component) (bool, error) {
	if target.IsHidden() {
		return false, nil
	}
	for key, value := range f.Labels {
		v, ok := target.Labels[key]
		if !ok || v != value {
			return false, nil
		}
	}
	return true, nil
}

type Complementer interface {
	Complement(target *component.Components) (*component.Components, error)
}

type DependenciesComplementer struct {
	AllComponents *component.Components
	Depth         int
}

func (c DependenciesComplementer) Complement(target *component.Components) (*component.Components, error) {
	result := component.NewComponents()
	for _, id := range target.GetIDs() {
		startingPoint, _ := target.Get(id)
		err := c.complement(0, startingPoint, result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c DependenciesComplementer) complement(depth int, pivot *component.Component, acc *component.Components) error {
	if c.Depth >= 0 && depth > c.Depth {
		return nil
	}
	if _, ok := acc.Get(pivot.ID); ok {
		return nil
	}

	acc.Add(pivot)

	for depID, _ := range pivot.Dependencies {
		dep, _ := c.AllComponents.Get(depID)
		err := c.complement(depth+1, dep, acc)
		if err != nil {
			return err
		}
	}

	return nil
}

type ReverseDependenciesComplementer struct {
	AllComponents *component.Components
	Depth         int
}

func (c ReverseDependenciesComplementer) Complement(target *component.Components) (*component.Components, error) {
	result := component.NewComponents()
	for _, id := range target.GetIDs() {
		startingPoint, _ := target.Get(id)
		err := c.complement(0, startingPoint, result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c ReverseDependenciesComplementer) complement(depth int, pivot *component.Component, acc *component.Components) error {
	if c.Depth >= 0 && depth > c.Depth {
		return nil
	}
	if _, ok := acc.Get(pivot.ID); ok {
		return nil
	}

	acc.Add(pivot)

	for _, rDepID := range c.AllComponents.GetIDs() {
		rDep, _ := c.AllComponents.Get(rDepID)
		for depID, _ := range rDep.Dependencies {
			if depID != pivot.ID {
				continue
			}
			err := c.complement(depth+1, rDep, acc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
