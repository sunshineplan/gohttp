package node

import (
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

func TestParentElement(t *testing.T) {
	if data := doc.Find(Tag("span")).Parent().Raw().Data; data != "h1" {
		t.Errorf("expected data %q; got %q", "h1", data)
	}
	if id, _ := doc.Find(Tag("span")).ParentElement().ParentElement().Attr().Get("id"); id != "5" {
		t.Errorf("expected id %q; got %q", "5", id)
	}
}

func TestNextPrevElement(t *testing.T) {
	if data := strings.TrimSpace(doc.Find(Tag("div"), Attr("id", "0")).NextSibling().Raw().Data); data != "check" {
		t.Errorf("expected data %q; got %q", "check", data)
	}
	if data := strings.TrimSpace(doc.Find(Tag("div"), Attr("id", "2")).PrevSibling().Raw().Data); data != "check" {
		t.Errorf("expected data %q; got %q", "check", data)
	}
	if data := doc.Find(Tag("table")).NextSiblingElement().Raw().Data; data != "div" {
		t.Errorf("expected data %q; got %q", "div", data)
	}
	if data := doc.Find(Tag("p")).PrevSiblingElement().Raw().Data; data != "div" {
		t.Errorf("expected data %q; got %q", "div", data)
	}
}

func TestText(t *testing.T) {
	if text := doc.Find(Tag("ul")).Find(Tag("li")).Text(); text != "To a " {
		t.Errorf("expected text %q; got %q", "To a ", text)
	}
}
func TestFullText(t *testing.T) {
	if text := doc.Find(Tag("ul")).Find(Tag("li")).FullText(); text != "To a JSP page right?" {
		t.Errorf("expected text %q; got %q", "To a JSP page right?", text)
	}
}

func TestFullTextEmpty(t *testing.T) {
	if text := doc.Find(Tag("div"), Attr("id", "5")).Find(Tag("h1")).FullText(); text != "" {
		t.Errorf("expected text %q; got %q", "", text)
	}
}

func TestHTML(t *testing.T) {
	if html := doc.Find(Tag("ul")).Find(Tag("li")).HTML(); html != "<li>To a <a href=\"hello.jsp\">JSP page</a> right?</li>" {
		t.Errorf("expected html %q; got %q", "<li>To a <a href=\"hello.jsp\">JSP page</a> right?</li>", html)
	}
}
