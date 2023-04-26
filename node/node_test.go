package node

import (
	"strconv"
	"strings"
	"testing"
)

const testHTML = `
<html>
  <head>
    <title>Sample "Hello, World" Application</title>
  </head>
  <body bgcolor=white>

    <table border="0" cellpadding="10">
      <tr>
        <td>
          <img src="images/springsource.png">
        </td>
        <td>
          <h1>Sample "Hello, World" Application</h1>
        </td>
      </tr>
    </table>
    <div id="0">
      <div id="1">Just two divs peacing out</div>
    </div>
    check
    <div id="2">One more</div>
    <p>This is the home page for the HelloWorld Web application. </p>
    <p>To prove that they work, you can execute either of the following links:
    <ul>
      <li>To a <a href="hello.jsp">JSP page</a> right?</li>
      <li>To a <a href="hello">servlet</a></li>
    </ul>
    </p>
    <div id="3">
      <div id="4">Last one</div>
    </div>
    <div id="5">
        <h1><span></span></h1>
    </div>
  </body>
</html>
`

const multipleClassesHTML = `
<html>
	<head>
		<title>Sample Application</title>
	</head>
	<body>
		<div class="first second">Multiple classes</div>
		<div class="first">Single class</div>
		<div class="second first third">Multiple classes inorder</div>
		<div>
			<div class="first">Inner single class</div>
			<div class="first second">Inner multiple classes</div>
			<div class="second first">Inner multiple classes inorder</div>
		</div>
	</body>
</html>
`

var (
	doc, _             = ParseHTML(testHTML)
	multipleClasses, _ = ParseHTML(multipleClassesHTML)
)

func TestFind(t *testing.T) {
	if src, _ := doc.Find("img").Attr().Get("src"); src != "images/springsource.png" {
		t.Errorf("expected src %q; got %q", "images/springsource.png", src)
	}
	if text := doc.Find("a", Attr("href", "hello")).Text(); text != "servlet" {
		t.Errorf("expected text %q; got %q", "servlet", text)
	}
	if text := doc.Find("div").Find("div").Text(); text != "Just two divs peacing out" {
		t.Errorf("expected text %q; got %q", "Just two divs peacing out", text)
	}
	if text := multipleClasses.Find("body").Find("").Text(); text != "Multiple classes" {
		t.Errorf("expected text %q; got %q", "Multiple classes", text)
	}
	if text := doc.Find("", Attr("id", "4")).Text(); text != "Last one" {
		t.Errorf("expected text %q; got %q", "Last one", text)
	}
	if id, _ := doc.Find("", Attr("text", "Last one")).Attr().Get("id"); id != "4" {
		t.Errorf("expected id %q; got %q", "4", id)
	}
}

func TestParentElement(t *testing.T) {
	if data := doc.Find("span").Parent().Raw().Data; data != "h1" {
		t.Errorf("expected data %q; got %q", "h1", data)
	}
	if id, _ := doc.Find("span").ParentElement().ParentElement().Attr().Get("id"); id != "5" {
		t.Errorf("expected id %q; got %q", "5", id)
	}
}

func TestNextPrevElement(t *testing.T) {
	if data := strings.TrimSpace(doc.Find("div", Attr("id", "0")).NextSibling().Raw().Data); data != "check" {
		t.Errorf("expected data %q; got %q", "check", data)
	}
	if data := strings.TrimSpace(doc.Find("div", Attr("id", "2")).PrevSibling().Raw().Data); data != "check" {
		t.Errorf("expected data %q; got %q", "check", data)
	}
	if data := doc.Find("table").NextSiblingElement().Raw().Data; data != "div" {
		t.Errorf("expected data %q; got %q", "div", data)
	}
	if data := doc.Find("p").PrevSiblingElement().Raw().Data; data != "div" {
		t.Errorf("expected data %q; got %q", "div", data)
	}
}

func TestFindAll(t *testing.T) {
	for i, div := range doc.FindAll("div") {
		id, _ := div.Attr().Get("id")
		if id, _ := strconv.Atoi(id); id != i {
			t.Errorf("expected id %d; got %d", i, id)
		}
	}
}

func TestFindAllBySingleClass(t *testing.T) {
	if l := len(multipleClasses.FindAll("div", Attr("class", "first"))); l != 6 {
		t.Errorf("expected length %d; got %d", 6, l)
	}
	if l := len(multipleClasses.FindAll("div", Attr("class", "third"))); l != 1 {
		t.Errorf("expected length %d; got %d", 1, l)
	}
}

func TestFindAllByAttribute(t *testing.T) {
	if l := len(doc.FindAll("", Attr("id", "2"))); l != 1 {
		t.Errorf("expected length %d; got %d", 1, l)
	}
}

func TestFindBySingleClass(t *testing.T) {
	if text := multipleClasses.Find("div", Attr("class", "first")).Text(); text != "Multiple classes" {
		t.Errorf("expected text %q; got %q", "Multiple classes", text)
	}
	if text := multipleClasses.Find("div", Attr("class", "third")).Text(); text != "Multiple classes inorder" {
		t.Errorf("expected text %q; got %q", "Multiple classes inorder", text)
	}
}

func TestFindAllStrict(t *testing.T) {
	if l := len(multipleClasses.FindAllStrict("div", Attr("class", "first second"))); l != 2 {
		t.Errorf("expected length %d; got %d", 2, l)
	}
	if l := len(multipleClasses.FindAllStrict("div", Attr("class", "first third second"))); l != 0 {
		t.Errorf("expected length %d; got %d", 0, l)
	}
	if l := len(multipleClasses.FindAllStrict("div", Attr("class", "second first third"))); l != 1 {
		t.Errorf("expected length %d; got %d", 1, l)
	}
}

func TestFindStrict(t *testing.T) {
	if text := multipleClasses.FindStrict("div", Attr("class", "first")).Text(); text != "Single class" {
		t.Errorf("expected text %q; got %q", "Single class", text)
	}
	if node := multipleClasses.FindStrict("div", Attr("class", "third")); node != nil {
		t.Errorf("expected node nil; got %v", node.Raw())
	}
}

func TestText(t *testing.T) {
	if text := doc.Find("ul").Find("li").Text(); text != "To a " {
		t.Errorf("expected text %q; got %q", "To a ", text)
	}
}
func TestFullText(t *testing.T) {
	if text := doc.Find("ul").Find("li").FullText(); text != "To a JSP page right?" {
		t.Errorf("expected text %q; got %q", "To a JSP page right?", text)
	}
}

func TestFullTextEmpty(t *testing.T) {
	if text := doc.Find("div", Attr("id", "5")).Find("h1").FullText(); text != "" {
		t.Errorf("expected text %q; got %q", "", text)
	}
}

func TestHTML(t *testing.T) {
	if html := doc.Find("ul").Find("li").HTML(); html != "<li>To a <a href=\"hello.jsp\">JSP page</a> right?</li>" {
		t.Errorf("expected html %q; got %q", "<li>To a <a href=\"hello.jsp\">JSP page</a> right?</li>", html)
	}
}
