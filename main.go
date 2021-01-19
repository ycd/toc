package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/net/html"
)

func main() {
	var toc TOC
	toc.options.Bulleted = true
	toc.options.Path = "test-markdown/TEST.md"

	resp, _ := toc.readFile()
	toc.parseHTML(resp)

	toc.writeToFile(string(resp))

}

var headers = map[string]int{"h1": 0, "h2": 1, "h3": 2, "h4": 3, "h5": 4, "h6": 5}

// Use 4 spaces for adding tabs
// See Markdown specification
// https://github.github.com/gfm/#tabs
var tab = "    "

type TOCConfig struct {
	Path     string
	Bulleted bool
}

type TOC struct {
	options TOCConfig
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

func (t *TOC) parseHTML(body []byte) {
	var f func(*html.Node)
	var delimiter string

	// Set delimiter
	if t.options.Bulleted == true {
		delimiter = "1."
	} else {
		delimiter = "-"
	}

	parsedMD, err := convertToHTML(body)
	if err != nil {
		log.Fatal(err)
	}

	doc, _ := html.Parse(strings.NewReader(parsedMD))

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

func (t *TOC) readFile() ([]byte, error) {
	file, err := ioutil.ReadFile(t.options.Path)
	if err != nil {
		return []byte{}, err
	}

	return file, nil
}

// Parse the file from path
// convert markdown file to html as string
func convertToHTML(file []byte) (string, error) {
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	if err := md.Convert(file, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// TODO(ycd): make file writing more
// memory efficient and safe
func (t *TOC) writeToFile(markdown string) {
	search := "<!--toc-->"

	idx := strings.Index(markdown, search) + len(search) + 1

	newText := markdown[:idx] + "\n" + t.String() + markdown[idx:]

	err := ioutil.WriteFile(t.options.Path, []byte(newText), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
