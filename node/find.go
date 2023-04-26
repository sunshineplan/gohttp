package node

import (
	"context"
	"strings"

	"golang.org/x/net/html"
)

type Tag string

func (tag Tag) Equal(s string) bool {
	return tag == "" || strings.ToLower(string(tag)) == s
}

func (n *node) find(tag Tag, once, strict bool, a ...Attribute) (nodes []Node) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var f func(*node, ...Attribute)
	f = func(node *node, a ...Attribute) {
		if ctx.Err() != nil {
			return
		}
		if raw := node.Raw(); n.Raw() != raw && raw.Type == html.ElementNode && tag.Equal(raw.Data) {
			if node.compareAttribute(strict, a...) {
				nodes = append(nodes, node)
				if once {
					cancel()
				}
			}
		}
		for node := node.FirstChild(); node != nil; node = node.NextSiblingElement() {
			f(newNode(node.Raw()), a...)
		}
	}
	f(n, a...)
	return
}

func (n *node) Find(tag Tag, a ...Attribute) Node {
	nodes := n.find(tag, true, false, a...)
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

func (n *node) FindAll(tag Tag, a ...Attribute) []Node {
	return n.find(tag, false, false, a...)
}

func (n *node) FindStrict(tag Tag, a ...Attribute) Node {
	nodes := n.find(tag, true, true, a...)
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

func (n *node) FindAllStrict(tag Tag, a ...Attribute) []Node {
	return n.find(tag, false, true, a...)
}
