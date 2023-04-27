package node

import (
	"context"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var _ TagOption = tag[string]{}

type Value interface {
	string | []string | *regexp.Regexp | bool | func(string, Node) bool
}

type TagOption interface {
	IsMatch(string, Node) bool
}

type tag[T Value] struct {
	tag T
}

func Tag[T Value](t T) TagOption {
	return tag[T]{t}
}

func (tag tag[T]) IsMatch(s string, node Node) bool {
	switch v := (any(tag.tag)).(type) {
	case string:
		return strings.ToLower(v) == s
	case []string:
		for _, v := range v {
			if strings.ToLower(v) == s {
				return true
			}
		}
	case *regexp.Regexp:
		return v.MatchString(s)
	case bool:
		return v
	case func(string, Node) bool:
		return v(s, node)
	}
	return false
}

func (n *node) find(tag TagOption, once bool, opts ...Option) (nodes []Node) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var f func(*node, ...Option)
	f = func(node *node, opts ...Option) {
		if ctx.Err() != nil {
			return
		}
		if raw := node.Raw(); n.Raw() != raw && raw.Type == html.ElementNode && (tag == nil || tag.IsMatch(raw.Data, node)) {
			ok := true
			for _, i := range opts {
				if !i.IsMatch(node) {
					ok = false
					break
				}
			}
			if ok {
				nodes = append(nodes, node)
				if once {
					cancel()
				}
			}
		}
		for node := node.FirstChild(); node != nil; node = node.NextSiblingElement() {
			f(newNode(node.Raw()), opts...)
		}
	}
	f(n, opts...)
	return
}

func (n *node) Find(tag TagOption, opts ...Option) Node {
	nodes := n.find(tag, true, opts...)
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

func (n *node) FindAll(tag TagOption, opts ...Option) []Node {
	return n.find(tag, false, opts...)
}
