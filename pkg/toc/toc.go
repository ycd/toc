package toc

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"toc/config"

	"github.com/fatih/color"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/net/html"
)

// Run handles the application logic.
func Run() {
	var toc toc
	// Create a FlagSet and sets the usage
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)

	// Configure the options from the flags/config file
	opts, err := config.ConfigureOptions(fs, os.Args[1:])
	if err != nil {
		config.UsageAndExit(err)
	}

	if opts.ShowHelp {
		config.HelpAndExit()
	}

	toc.Options.Path = opts.Path
	toc.Options.Bulleted = opts.Bulleted
	toc.Options.Append = opts.Append

	toc.logic()
}

func (t *toc) logic() {
	resp, err := t.readFile()
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	err = t.parseHTML(resp)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	if t.Options.Append == true {
		err = t.writeToFile(string(resp))
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Print(t.String())
	}

	color.HiGreen("✔ Table of contents generated successfully.")
}

func (t *toc) String() (s string) {
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

func (t *toc) getDelimiter(header int) string {
	// Set delimiter
	if t.Options.Bulleted == true {
		if header >= 1 {
			return "*"
		}
		return "-"
	}
	return "1."
}

func (t *toc) parseHTML(body []byte) error {
	var f func(*html.Node)

	parsedMD, err := convertToHTML(body)
	if err != nil {
		log.Fatal(err)
	}

	doc, _ := html.Parse(strings.NewReader(parsedMD))

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && isHeader(n.Data) {
			headerVal := getHeaderValue(n.Data)
			val := fmt.Sprintf("%s%s [%s](#%s)\n", strings.Repeat(tab, headerVal), t.getDelimiter(headerVal), n.FirstChild.Data, n.Attr[0].Val)
			t.add(val)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return nil
}

func (t *toc) add(content string) {
	t.Content = append(t.Content, content)
}

func (t *toc) readFile() ([]byte, error) {
	if _, err := os.Stat(t.Options.Path); os.IsNotExist(err) {
		return nil, fmt.Errorf(fmt.Sprintf("path (%s) doesn't exists", t.Options.Path))
	}

	file, err := ioutil.ReadFile(t.Options.Path)
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
func (t *toc) writeToFile(markdown string) error {
	search := "<!--toc-->"

	idx := strings.Index(markdown, search) + len(search) + 1
	if idx == 10 {
		return errors.New("toc path is missing, add '<!--toc--->' to your markdown")
	}

	newText := markdown[:idx] + "\n" + t.String() + markdown[idx:]

	err := ioutil.WriteFile(t.Options.Path, []byte(newText), 0644)
	if err != nil {
		return err
	}

	return nil
}
