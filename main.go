package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/net/html"
)

func main() {
	var toc TOC
	toc.Options.Bulleted = true

	path := "test-markdown/TEST.md"

	resp, _ := readParse(path)
	toc.parseHTML(resp)

	fmt.Println(toc.String())

}

var headers = map[string]int{"h1": 0, "h2": 1, "h3": 2, "h4": 3, "h5": 4, "h6": 5}

// Use 4 spaces for adding tabs
// See Markdown specification
// https://github.github.com/gfm/#tabs
var tab = "    "

type TOC struct {
	Options struct {
		Bulleted bool
	}
	Content []string
}

func (t *TOC) String() (s string) {
	for _, v := range t.Content {
		s += v
	}

	return
}

func isHeader(attr string) bool {
	for k := range headers {
		if attr == k {
			return true
		}
	}
	return false
}

func getHeaderValue(header string) int {
	return headers[header]
}

func (t *TOC) parseHTML(body string) {
	var f func(*html.Node)
	var delimiter string

	// Set delimiter
	if t.Options.Bulleted == true {
		delimiter = "1."
	} else {
		delimiter = "-"
	}

	doc, _ := html.Parse(strings.NewReader(body))

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && isHeader(n.Data) {
			t.add(fmt.Sprintf("%s%s [%s](#%s)\n", strings.Repeat(tab, getHeaderValue(n.Data)), delimiter, n.FirstChild.Data, n.Attr[0].Val))
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

}

func (t *TOC) add(content string) {
	t.Content = append(t.Content, content)
}

// Parse the file from path
// convert markdown file to html as string
func readParse(path string) (string, error) {
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	if err := md.Convert(file, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
