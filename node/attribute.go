package node

import "strings"

var (
	_ Attribute  = attribute{}
	_ Attributes = attributes{}
)

type Attribute interface {
	Key() string
	Value() string
}

type Attributes interface {
	Range(func(Attribute) bool)
	Get(string) (string, bool)
}

type attribute [2]string

func Attr(key, value string) Attribute {
	return attribute{key, value}
}

func (attr attribute) Key() string {
	return attr[0]
}
func (attr attribute) Value() string {
	return attr[1]
}

type attributes map[string]string

func (attrs attributes) Range(f func(Attribute) bool) {
	for k, v := range attrs {
		if !f(Attr(k, v)) {
			break
		}
	}
}

func (attrs attributes) Get(key string) (string, bool) {
	v, ok := attrs[key]
	return v, ok
}

func (n *node) Attr() Attributes {
	if len(n.Node.Attr) == 0 {
		return nil
	}
	attrs := make(attributes)
	for _, i := range n.Node.Attr {
		if _, ok := attrs[i.Key]; !ok {
			attrs[i.Key] = i.Val
		}
	}
	return attrs
}

func (n *node) compareAttribute(strict bool, a ...Attribute) bool {
	for _, i := range a {
		var b bool
		if key := strings.ToLower(i.Key()); key == "string" || key == "text" {
			// https://www.crummy.com/software/BeautifulSoup/bs4/doc/#the-string-argument
			b = n.Text() == i.Value()
		} else if attr := n.Attr(); attr != nil {
			if v, ok := attr.Get(key); ok && compare(v, i.Value(), strict) {
				b = true
			}
		}
		if !b {
			return false
		}
	}
	return true
}

func compare(a, b string, strict bool) bool {
	av, bv := strings.Fields(a), strings.Fields(b)
	if strict {
		return strings.Join(av, "|||") == strings.Join(bv, "|||")
	}
	for _, i := range bv {
		var b bool
		for _, ii := range av {
			if i == ii {
				b = true
				break
			}
		}
		if !b {
			return false
		}
	}
	return true
}
