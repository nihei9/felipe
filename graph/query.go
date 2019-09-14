package graph

func Query(cs *Components, cond *Condition, _ *ComplementMethod) (*Components, error) {
	matches := NewComponents()
	for _, c := range cs.Components {
		ok, err := Match(c, cond)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		err = matches.Add(c)
		if err != nil {
			return nil, err
		}
	}

	return matches, nil
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
