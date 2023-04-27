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
	Children() []Node
	PrevSibling() Node
	NextSibling() Node

	Attrs() Attributes
	Text() string
	HTML() string
	FullText() string

	Find(Filter, ...Filter) Node
	FindN(int, Filter, ...Filter) []Node
	FindAll(Filter, ...Filter) []Node
	FindParent(Filter, ...Filter) Node
	FindParentsN(int, Filter, ...Filter) []Node
	FindParents(Filter, ...Filter) []Node
	FindPrevSibling(Filter, ...Filter) Node
	FindPrevSiblingsN(int, Filter, ...Filter) []Node
	FindPrevSiblings(Filter, ...Filter) []Node
	FindNextSibling(Filter, ...Filter) Node
	FindNextSiblingsN(int, Filter, ...Filter) []Node
	FindNextSiblings(Filter, ...Filter) []Node
	FindPrevious(Filter, ...Filter) Node
	FindPreviousN(int, Filter, ...Filter) []Node
	FindAllPrevious(Filter, ...Filter) []Node
	FindNext(Filter, ...Filter) Node
	FindNextN(int, Filter, ...Filter) []Node
	FindAllNext(Filter, ...Filter) []Node
}

func NewNode(n *html.Node) Node {
	if n == nil {
		return nil
	}
	return &node{n}
}

type node struct {
	*html.Node
}

func Parse(r io.Reader) (Node, error) {
	return ParseWithOptions(r)
}

func ParseWithOptions(r io.Reader, opts ...html.ParseOption) (Node, error) {
	n, err := html.ParseWithOptions(r, opts...)
	if err != nil {
		return nil, err
	}
	return NewNode(n), nil
}

func ParseHTML(s string) (Node, error) {
	return Parse(strings.NewReader(s))
}

func (n *node) Raw() *html.Node {
	return n.Node
}

func (n *node) Parent() Node {
	return NewNode(n.Node.Parent)
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
	return NewNode(n.Node.FirstChild)
}

func (n *node) LastChild() Node {
	return NewNode(n.Node.LastChild)
}

func (n *node) PrevSibling() Node {
	return NewNode(n.Node.PrevSibling)
}

func (n *node) NextSibling() Node {
	return NewNode(n.Node.NextSibling)
}

func (n *node) Attrs() Attributes {
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
