package node

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

var _ Node = &node{}

type Node interface {
	Raw() *html.Node

	Parent() Node
	FirstChild() Node
	LastChild() Node
	PrevSibling() Node
	NextSibling() Node

	ParentElement() Node
	Children() []Node
	PrevSiblingElement() Node
	NextSiblingElement() Node

	Attr() Attributes
	Text() string
	HTML() string
	FullText() string

	Find(tag Tag, a ...Attribute) Node
	FindAll(tag Tag, a ...Attribute) []Node
	FindStrict(tag Tag, a ...Attribute) Node
	FindAllStrict(tag Tag, a ...Attribute) []Node
}

func NewNode(node *html.Node) Node {
	return newNode(node)
}

type node struct {
	*html.Node
}

func newNode(n *html.Node) *node {
	return &node{n}
}

func Parse(r io.Reader) (Node, error) {
	return ParseWithOptions(r)
}

func ParseWithOptions(r io.Reader, opts ...html.ParseOption) (Node, error) {
	n, err := html.ParseWithOptions(r, opts...)
	if err != nil {
		return nil, err
	}
	return newNode(n), nil
}

func ParseHTML(s string) (Node, error) {
	return Parse(strings.NewReader(s))
}

func (n *node) Raw() *html.Node {
	return n.Node
}

func (n *node) Parent() Node {
	if n := n.Node.Parent; n != nil {
		return newNode(n)
	}
	return nil
}

func (n *node) ParentElement() Node {
	if node := n.Parent(); node == nil || node.Raw().Type == html.ElementNode {
		return node
	} else {
		return node.ParentElement()
	}
}

func (n *node) Children() (children []Node) {
	child := n.FirstChild()
	for child != nil {
		children = append(children, child)
		child = child.NextSibling()
	}
	return
}

func (n *node) FirstChild() Node {
	if n := n.Node.FirstChild; n != nil {
		return newNode(n)
	}
	return nil
}

func (n *node) LastChild() Node {
	if n := n.Node.LastChild; n != nil {
		return newNode(n)
	}
	return nil
}

func (n *node) PrevSibling() Node {
	if n := n.Node.PrevSibling; n != nil {
		return newNode(n)
	}
	return nil
}

func (n *node) NextSibling() Node {
	if n := n.Node.NextSibling; n != nil {
		return newNode(n)
	}
	return nil
}

func (n *node) PrevSiblingElement() Node {
	if node := n.PrevSibling(); node == nil || node.Raw().Type == html.ElementNode {
		return node
	} else {
		return node.PrevSiblingElement()
	}
}

func (n *node) NextSiblingElement() Node {
	if node := n.NextSibling(); node == nil || node.Raw().Type == html.ElementNode {
		return node
	} else {
		return node.NextSiblingElement()
	}
}

func (n *node) Text() string {
	for node := n.FirstChild(); node != nil; node = node.NextSibling() {
		node := node.Raw()
		if node.Type != html.TextNode {
			continue
		}
		if s := strings.TrimSpace(node.Data); s != "" {
			return node.Data
		}
	}
	return ""
}

func (n *node) HTML() string {
	var b strings.Builder
	html.Render(&b, n.Raw())
	return b.String()
}

func (n *node) FullText() string {
	var b strings.Builder
	var f func(Node)
	f = func(node Node) {
		if node == nil {
			return
		}
		switch raw := node.Raw(); raw.Type {
		case html.TextNode:
			b.WriteString(raw.Data)
		case html.ElementNode:
			f(node.FirstChild())
		}
		if node := node.NextSibling(); node != nil {
			f(node)
		}
	}
	f(n.FirstChild())
	return b.String()
}
