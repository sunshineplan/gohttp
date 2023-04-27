package node

import (
	"strconv"
	"testing"
)

func TestFind(t *testing.T) {
	if src, _ := doc.Find(Tag("img")).Attrs().Get("src"); src != "images/springsource.png" {
		t.Errorf("expected src %q; got %q", "images/springsource.png", src)
	}
	if text := doc.Find(Tag("a"), Attr("href", "hello")).Text(); text != "servlet" {
		t.Errorf("expected text %q; got %q", "servlet", text)
	}
	if text := doc.Find(Tag("div")).Find(Tag("div")).Text(); text != "Just two divs peacing out" {
		t.Errorf("expected text %q; got %q", "Just two divs peacing out", text)
	}
	if text := multipleClasses.Find(Tag("body")).Find(nil).Text(); text != "Multiple classes" {
		t.Errorf("expected text %q; got %q", "Multiple classes", text)
	}
	if text := doc.Find(True, Attr("id", "4")).Text(); text != "Last one" {
		t.Errorf("expected text %q; got %q", "Last one", text)
	}
	if id, _ := doc.Find(True, Text("Last one")).Attrs().Get("id"); id != "4" {
		t.Errorf("expected id %q; got %q", "4", id)
	}
}

func TestFindAll(t *testing.T) {
	for i, div := range doc.FindAll(Tag("div")) {
		id, _ := div.Attrs().Get("id")
		if id, _ := strconv.Atoi(id); id != i {
			t.Errorf("expected id %d; got %d", i, id)
		}
	}
}

func TestFindAllBySingleClass(t *testing.T) {
	if l := len(multipleClasses.FindAll(Tag("div"), Class("first"))); l != 6 {
		t.Errorf("expected length %d; got %d", 6, l)
	}
	if l := len(multipleClasses.FindAll(Tag("div"), Class("third"))); l != 1 {
		t.Errorf("expected length %d; got %d", 1, l)
	}
}

func TestFindAllByAttribute(t *testing.T) {
	if l := len(doc.FindAll(nil, Attr("id", "2"))); l != 1 {
		t.Errorf("expected length %d; got %d", 1, l)
	}
}

func TestFindBySingleClass(t *testing.T) {
	if text := multipleClasses.Find(Tag("div"), Class("first")).Text(); text != "Multiple classes" {
		t.Errorf("expected text %q; got %q", "Multiple classes", text)
	}
	if text := multipleClasses.Find(Tag("div"), Class("third")).Text(); text != "Multiple classes inorder" {
		t.Errorf("expected text %q; got %q", "Multiple classes inorder", text)
	}
}

func TestFindAllStrict(t *testing.T) {
	if l := len(multipleClasses.FindAll(Tag("div"), ClassStrict("first second"))); l != 2 {
		t.Errorf("expected length %d; got %d", 2, l)
	}
	if l := len(multipleClasses.FindAll(Tag("div"), ClassStrict("first third second"))); l != 0 {
		t.Errorf("expected length %d; got %d", 0, l)
	}
	if l := len(multipleClasses.FindAll(Tag("div"), ClassStrict("second first third"))); l != 1 {
		t.Errorf("expected length %d; got %d", 1, l)
	}
}

func TestFindStrict(t *testing.T) {
	if text := multipleClasses.Find(Tag("div"), ClassStrict("first")).Text(); text != "Single class" {
		t.Errorf("expected text %q; got %q", "Single class", text)
	}
	if node := multipleClasses.Find(Tag("div"), ClassStrict("third")); node != nil {
		t.Errorf("expected node nil; got %v", node.Raw())
	}
}
