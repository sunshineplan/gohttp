package node

import (
	"context"

	"golang.org/x/net/html"
)

type findMethod int

const (
	down findMethod = iota
	up
	prevSibling
	nextSibling
	prev
	next
)

func (n *node) find(method findMethod, limit int, tag Filter, filters ...Filter) (nodes []Node) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var f func(Node, ...Filter)
	f = func(node Node, filters ...Filter) {
		if ctx.Err() != nil || node == nil {
			return
		}
		if raw := node.Raw(); n.Raw() != raw && raw.Type == html.ElementNode && (tag == nil || tag.IsMatch(node)) {
			ok := true
			for _, i := range filters {
				if !i.IsMatch(node) {
					ok = false
					break
				}
			}
			if ok {
				nodes = append(nodes, node)
				if len(nodes) == limit {
					cancel()
				}
			}
		}
		switch method {
		case down:
			for node := node.FirstChild(); node != nil; node = node.NextSibling() {
				f(NewNode(node.Raw()), filters...)
			}
		case up:
			for node := node.Parent(); node != nil; node = node.Parent() {
				f(NewNode(node.Raw()), filters...)
			}
		case prevSibling:
			for node := node.PrevSibling(); node != nil; node = node.PrevSibling() {
				f(NewNode(node.Raw()), filters...)
			}
		case nextSibling:
			for node := node.NextSibling(); node != nil; node = node.NextSibling() {
				f(NewNode(node.Raw()), filters...)
			}
		case prev:
			for {
				prev := node.PrevSibling()
				if prev == nil {
					if prev = node.Parent(); prev == nil {
						return
					}
				}
				f(NewNode(prev.Raw()), filters...)
			}
		case next:
			for {
				next := node.FirstChild()
				if next == nil {
					if next = node.NextSibling(); next == nil {
						if parent := node.Parent(); parent != nil {
							if next = parent.NextSibling(); next == nil {
								return
							}
						}
					}
				}
				f(NewNode(next.Raw()), filters...)
			}
		}
	}
	f(n, filters...)
	return
}

func (n *node) findOnce(method findMethod, tag Filter, filters ...Filter) Node {
	nodes := n.find(method, 1, tag, filters...)
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

func (n *node) findN(method findMethod, limit int, tag Filter, filters ...Filter) []Node {
	if limit <= 0 {
		return nil
	}
	return n.find(method, limit, tag, filters...)
}

func (n *node) Find(tag Filter, filters ...Filter) Node {
	return n.findOnce(down, tag, filters...)
}

func (n *node) FindN(limit int, tag Filter, filters ...Filter) []Node {
	return n.findN(down, limit, tag, filters...)
}

func (n *node) FindAll(tag Filter, filters ...Filter) []Node {
	return n.find(down, 0, tag, filters...)
}

func (n *node) FindParent(tag Filter, filters ...Filter) Node {
	return n.findOnce(up, tag, filters...)
}

func (n *node) FindParentsN(limit int, tag Filter, filters ...Filter) []Node {
	return n.findN(up, limit, tag, filters...)
}

func (n *node) FindParents(tag Filter, filters ...Filter) []Node {
	return n.find(up, 0, tag, filters...)
}

func (n *node) FindPrevSibling(tag Filter, filters ...Filter) Node {
	return n.findOnce(prevSibling, tag, filters...)
}

func (n *node) FindPrevSiblingsN(limit int, tag Filter, filters ...Filter) []Node {
	return n.findN(prevSibling, limit, tag, filters...)
}

func (n *node) FindPrevSiblings(tag Filter, filters ...Filter) []Node {
	return n.find(prevSibling, 0, tag, filters...)
}

func (n *node) FindNextSibling(tag Filter, filters ...Filter) Node {
	return n.findOnce(nextSibling, tag, filters...)
}

func (n *node) FindNextSiblingsN(limit int, tag Filter, filters ...Filter) []Node {
	return n.findN(nextSibling, limit, tag, filters...)
}

func (n *node) FindNextSiblings(tag Filter, filters ...Filter) []Node {
	return n.find(nextSibling, 0, tag, filters...)
}

func (n *node) FindPrevious(tag Filter, filters ...Filter) Node {
	return n.findOnce(prev, tag, filters...)
}

func (n *node) FindPreviousN(limit int, tag Filter, filters ...Filter) []Node {
	return n.findN(prev, limit, tag, filters...)
}

func (n *node) FindAllPrevious(tag Filter, filters ...Filter) []Node {
	return n.find(prev, 0, tag, filters...)
}

func (n *node) FindNext(tag Filter, filters ...Filter) Node {
	return n.findOnce(next, tag, filters...)
}

func (n *node) FindNextN(limit int, tag Filter, filters ...Filter) []Node {
	return n.findN(next, limit, tag, filters...)
}

func (n *node) FindAllNext(tag Filter, filters ...Filter) []Node {
	return n.find(next, 0, tag, filters...)
}
