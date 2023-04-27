package node

import (
	"regexp"
	"strings"
)

var (
	_ Option = attribute[string]{}
	_ Option = class[string]{}
	_ Option = classStrict("")
	_ Option = text[string]{}
)

type Option interface {
	Name() string
	IsMatch(node Node) bool
}

type attribute[T Value] struct {
	name  string
	value T
}

func Attr[T Value](name string, value T) Option {
	return attribute[T]{strings.ToLower(name), value}
}

func (attribute attribute[T]) Name() string {
	return attribute.name
}

func (attribute attribute[T]) IsMatch(node Node) bool {
	switch v := (any(attribute.value)).(type) {
	case string:
		if name := attribute.Name(); name == "class" {
			return class[T]{attribute.value}.IsMatch(node)
		} else if value, ok := getAttribute(node, name); !ok {
			return false
		} else {
			return value == v
		}
	case []string:
		if name := attribute.Name(); name == "class" {
			return class[T]{attribute.value}.IsMatch(node)
		} else if value, ok := getAttribute(node, name); !ok {
			return false
		} else {
			for _, v := range v {
				if value == v {
					return true
				}
			}
		}
	case *regexp.Regexp:
		if value, ok := getAttribute(node, attribute.Name()); !ok {
			return false
		} else {
			return v.MatchString(value)
		}
	case bool:
		_, ok := getAttribute(node, attribute.Name())
		return v == ok
	case func(string, Node) bool:
		if value, ok := getAttribute(node, attribute.Name()); !ok {
			return false
		} else {
			return v(value, node)
		}
	}
	return false
}

type class[T Value] struct {
	value T
}

func Class[T Value](v T) Option {
	return class[T]{v}
}

func (class[T]) Name() string {
	return "class"
}

func (cls class[T]) IsMatch(node Node) bool {
	switch v := (any(cls.value)).(type) {
	case string:
		nodeClass, ok := getAttribute(node, cls.Name())
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
		if _, ok := getAttribute(node, cls.Name()); !ok {
			return false
		}
		for _, v := range v {
			if (class[string]{v}).IsMatch(node) {
				return true
			}
		}
		return false
	default:
		return attribute[T]{"class", cls.value}.IsMatch(node)
	}
}

type classStrict string

func ClassStrict(cls string) Option {
	return classStrict(cls)
}

func (classStrict) Name() string {
	return "class"
}

func (classStrict classStrict) IsMatch(node Node) bool {
	nodeClass, ok := getAttribute(node, classStrict.Name())
	if !ok {
		return false
	}
	classA, classB := strings.Fields(nodeClass), strings.Fields(string(classStrict))
	return strings.Join(classA, "|||") == strings.Join(classB, "|||")
}

type text[T Value] struct {
	text T
}

func Text[T Value](t T) Option {
	return text[T]{t}
}

func String[T Value](t T) Option {
	return Text(t)
}

func (text[T]) Name() string {
	return "text"
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
	case bool:
		return node.Text() != ""
	case func(string, Node) bool:
		return v(node.Text(), node)
	}
	return false
}

func getAttribute(node Node, name string) (string, bool) {
	if attr := node.Attr(); attr == nil {
		return "", false
	} else if attr, ok := attr.Get(name); !ok || attr == "" {
		return "", false
	} else {
		return attr, true
	}
}
