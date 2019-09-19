package graph

func Query(cs *Components, cond *Condition, _ *ComplementMethod) (*Components, error) {
	result := NewComponents()
	for _, c := range cs.components {
		ok, err := Match(c, cond)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		err = result.add(c)
		if err != nil {
			return nil, err
		}
	}

	// complement a result
	ds := []ComponentID{}
	for _, c := range result.Components() {
		for _, d := range c.Dependencies() {
			ds = append(ds, d)
		}
	}
	for len(ds) > 0 {
		dds := []ComponentID{}
		for _, d := range ds {
			if _, ok := result.Get(d); ok {
				continue
			}
			dc, _ := cs.Get(d)
			result.add(dc)

			for _, dd := range dc.Dependencies() {
				if _, ok := result.Get(dd); ok {
					continue
				}
				dds = append(dds, dd)
			}
		}
		ds = dds
	}

	return result, nil
}

func Match(c *Component, cond *Condition) (bool, error) {
	if !c.queryable {
		return false, nil
	}
	for _, m := range cond.matchers {
		ok, err := m.Match(c)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	return false, nil
}

type Condition struct {
	matchers []Matcher
}

func NewCondition() *Condition {
	return &Condition{
		matchers: []Matcher{},
	}
}

func (cond *Condition) AddMatcher(m Matcher) {
	cond.matchers = append(cond.matchers, m)
}

type Matcher interface {
	Match(*Component) (bool, error)
}

// AnyMatcher is implementation of the Matcher interface.
type AnyMatcher struct {
}

func NewAnyMatcher() *AnyMatcher {
	return &AnyMatcher{}
}

func (m *AnyMatcher) Match(c *Component) (bool, error) {
	return true, nil
}

// LabelsMatcher is implementation of the Matcher interface.
type LabelsMatcher struct {
	labels map[string]string
}

func NewLabelsMatcher(labels map[string]string) *LabelsMatcher {
	m := &LabelsMatcher{
		labels: map[string]string{},
	}
	for k, v := range labels {
		m.labels[k] = v
	}

	return m
}

func (m *LabelsMatcher) Match(c *Component) (bool, error) {
	for condK, condV := range m.labels {
		vs, ok := c.Labels()[condK]
		for _, v := range vs {
			if ok && v == condV {
				return true, nil
			}
		}
	}

	return false, nil
}

type ComplementMethod struct {
}
