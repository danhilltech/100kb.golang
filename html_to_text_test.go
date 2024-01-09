package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
)

func TestParserSimple1(t *testing.T) {
	test := `<html>
	<nav>
	  <p>text</p>
	</nav>
	<p>here</p>
	</html>`

	reader := strings.NewReader(test)

	body, _, _, err := HtmlToText(reader)

	if err != nil {
		t.Fatal(err)
	}
	if body != "here" {
		t.Fail()
	}

}

func TestParserSimple2(t *testing.T) {
	test := `<html>
	<nav>
	  <p>text</p>
	</nav>
	<p>here</p>
	<p>there</p>
	</html>`

	reader := strings.NewReader(test)

	body, _, _, err := HtmlToText(reader)

	if err != nil {
		t.Fatal(err)
	}
	if body != "here\nthere" {
		t.Fail()
	}

}

func TestParserSimple3(t *testing.T) {
	test := `<html>
	<nav>

	  <p>text</p>
	</nav>
	<p>here</p>
	<p>there</p>
	<footer>
	  <p>now</p>
	</footer>
	</html>`

	reader := strings.NewReader(test)

	body, _, _, err := HtmlToText(reader)

	if err != nil {
		t.Fatal(err)
	}
	if body != "here\nthere" {
		t.Fail()
	}

}

func TestParserSimple4(t *testing.T) {
	test := `<!DOCTYPE html><html>
	<head>
	<title>Hi Dan</title>
	<meta content='deschere' property='og:description'/>
	</head>
	<body>
	<header>
	  <p>text</p>
	</header>
	<p>here</p>
	<p>there</p>
	<div>bob</div>
	<div><div>
		<div class="share">share</div>
	  <p>now <span class="share">this</span><span>two </span></p>
	</div></div>
	</body>
	</html>`

	reader := strings.NewReader(test)

	body, title, desc, err := HtmlToText(reader)

	if err != nil {
		t.Fatal(err)
	}
	if body != "text\nhere\nthere\nnow two" {
		t.Fail()
	}
	if title != "Hi Dan" {
		t.Fail()
	}
	if desc != "deschere" {
		t.Fail()
	}

}

func TestParserHTML(t *testing.T) {
	for i := 1; i <= 6; i++ {

		reader, err := os.Open(fmt.Sprintf("testcases/%d.html", i))
		if err != nil {
			t.Fatal(err)
		}

		expected, err := os.Open(fmt.Sprintf("testcases/%d.txt", i))
		if err != nil {
			t.Fatal(err)
		}

		expectedStr, err := io.ReadAll(expected)
		if err != nil {
			t.Fatal(err)
		}

		body, _, _, err := HtmlToText(reader)

		if err != nil {
			t.Fatal(err)
		}

		if strings.TrimSpace(body) != strings.TrimSpace(string(expectedStr)) {
			fmt.Printf("HTML %d\n**********\n\n", i)
			fmt.Println(diff.CharacterDiff(strings.TrimSpace(body), strings.TrimSpace(string(expectedStr))))
			t.Fail()
		}
	}
}
