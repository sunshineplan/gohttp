package node

import (
	"regexp"
	"strings"
)

var (
	_ Filter = everything{}
	_ Filter = tag[string]{}
	_ Filter = attribute[string]{}
	_ Filter = class[string]{}
	_ Filter = classStrict("")
	_ Filter = text[string]{}
)

type Value interface {
	string | []string | *regexp.Regexp | everything | func(string, Node) bool
}

type Filter interface {
	IsMatch(node Node) bool
}

var True everything

type everything struct{}

func (everything) IsMatch(Node) bool {
	return true
}

type tag[T Value] struct {
	tag T
}

func Tag[T Value](t T) Filter {
	return tag[T]{t}
}

func (tag tag[T]) IsMatch(node Node) bool {
	switch v := (any(tag.tag)).(type) {
	case string:
		return strings.ToLower(v) == node.Raw().Data
	case []string:
		for _, v := range v {
			if strings.ToLower(v) == node.Raw().Data {
				return true
			}
		}
	case *regexp.Regexp:
		return v.MatchString(node.Raw().Data)
	case everything:
		return true
	case func(string, Node) bool:
		return v(node.Raw().Data, node)
	}
	return false
}

type attribute[T Value] struct {
	name  string
	value T
}

func Attr[T Value](name string, value T) Filter {
	return attribute[T]{strings.ToLower(name), value}
}

func (attribute attribute[T]) IsMatch(node Node) bool {
	switch v := (any(attribute.value)).(type) {
	case string:
		if attribute.name == "class" {
			return class[T]{attribute.value}.IsMatch(node)
		} else if value, ok := getAttribute(node, attribute.name); !ok {
			return false
		} else {
			return value == v
		}
	case []string:
		if attribute.name == "class" {
			return class[T]{attribute.value}.IsMatch(node)
		} else if value, ok := getAttribute(node, attribute.name); !ok {
			return false
		} else {
			for _, v := range v {
				if value == v {
					return true
				}
			}
		}
	case *regexp.Regexp:
		if value, ok := getAttribute(node, attribute.name); !ok {
			return false
		} else {
			return v.MatchString(value)
		}
	case everything:
		_, ok := getAttribute(node, attribute.name)
		return ok
	case func(string, Node) bool:
		if value, ok := getAttribute(node, attribute.name); !ok {
			return false
		} else {
			return v(value, node)
		}
	}
	return false
}

type class[T Value] struct {
	class T
}

func Class[T Value](v T) Filter {
	return class[T]{v}
}

func (cls class[T]) IsMatch(node Node) bool {
	switch v := (any(cls.class)).(type) {
	case string:
		nodeClass, ok := getAttribute(node, "class")
		if !ok {
			return false
		}
		classA, classB := strings.Fields(nodeClass), strings.Fields(v)
		for _, i := range classB {
			var b bool
			for _, ii := range classA {
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
	case []string:
		if _, ok := getAttribute(node, "class"); !ok {
			return false
		}
		for _, v := range v {
			if (class[string]{v}).IsMatch(node) {
				return true
			}
		}
		return false
	default:
		return attribute[T]{"class", cls.class}.IsMatch(node)
	}
}

type classStrict string

func ClassStrict(cls string) Filter {
	return classStrict(cls)
}

func (classStrict classStrict) IsMatch(node Node) bool {
	nodeClass, ok := getAttribute(node, "class")
	if !ok {
		return false
	}
	classA, classB := strings.Fields(nodeClass), strings.Fields(string(classStrict))
	return strings.Join(classA, "|||") == strings.Join(classB, "|||")
}

type text[T Value] struct {
	text T
}

func Text[T Value](t T) Filter {
	return text[T]{t}
}

func String[T Value](t T) Filter {
	return Text(t)
}

func (text text[T]) IsMatch(node Node) bool {
	switch v := (any(text.text)).(type) {
	case string:
		return node.Text() == v
	case []string:
		for _, v := range v {
			if node.Text() == v {
				return true
			}
		}
	case *regexp.Regexp:
		return v.MatchString(node.Text())
	case everything:
		return node.Text() != ""
	case func(string, Node) bool:
		return v(node.Text(), node)
	}
	return false
}

func getAttribute(node Node, name string) (string, bool) {
	if attr := node.Attrs(); attr == nil {
		return "", false
	} else {
		attr, ok := attr.Get(name)
		return attr, ok
	}
}
