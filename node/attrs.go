package node

var _ Attributes = attributes{}

type Attributes interface {
	Range(func(k, v string) bool)
	Get(string) (string, bool)
}

type attributes map[string]string

func (attrs attributes) Range(f func(string, string) bool) {
	for k, v := range attrs {
		if !f(k, v) {
			break
		}
	}
}

func (attrs attributes) Get(attr string) (string, bool) {
	v, ok := attrs[attr]
	return v, ok
}
