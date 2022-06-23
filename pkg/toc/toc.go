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
	"regexp"
	"strings"

	"github.com/ycd/toc/config"

	"github.com/fatih/color"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/net/html"
)

var (
	startRe  = regexp.MustCompile(`(?i)\<!--\s*toc\s*--\>`)
	finishRe = regexp.MustCompile(`(?i)\<!--\s*(end of toc|tocstop|/TOC)\s*--\>`)
)

var skip = 0

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
	toc.Options.Skip = opts.Skip
	toc.Options.Depth = opts.Depth

	toc.logic()
}

func (t *toc) logic() {
	red := color.New(color.FgRed, color.Bold).PrintlnFunc()

	resp, err := t.readFile()
	if err != nil {
		red("ERROR: " + err.Error())
		os.Exit(1)
	}

	err = t.parseHTML(resp)
	if err != nil {
		red("ERROR: " + err.Error())
		os.Exit(1)
	}

	if !t.Options.Append {
		fmt.Print(t.String())
		return
	}

	if err = t.writeToFile(string(resp)); err != nil {
		red("ERROR: " + err.Error())
		os.Exit(1)
	}
	color.Green("✔ Table of contents generated successfully")
}

func (t *toc) String() (s string) {
	if len(t.Content) == 0 {
		color.Red("ERROR: skip value is bigger than the length of table of contents")
		os.Exit(1)
	}
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

func (t *toc) Last() int {
	if len(t.Content) >= 1 {
		switch spaceCount(t.Content[len(t.Content)-1]) {
		case 0:
			return 0
		case 4:
			return 1
		case 8:
			return 2
		case 12:
			return 3
		case 16:
			return 4
		case 20:
			return 5
		}
	}
	return 0
}

func spaceCount(s string) int {
	return len(s) - len(strings.TrimLeft(s, " "))
}

func (t *toc) getDelimiter(header int) string {
	// Set delimiter
	if !t.Options.Bulleted {
		return "1."
	}

	if header >= 1 {
		return "*"
	}

	return "-"
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

			if headerVal >= t.Options.Skip && headerVal < t.Options.Depth {
				t.add(fmt.Sprintf("%s%s [%s](#%s)\n", strings.Repeat(tab, headerVal-t.Options.Skip), t.getDelimiter(headerVal-t.Options.Skip), n.FirstChild.Data, n.Attr[0].Val))
			}
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
		return nil, fmt.Errorf("path '%s' doesn't exists", t.Options.Path)
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

// reformatMarkdown loads the entire string in the memory,
// finds the end and starting position for pos
// deletes the older one and creates a new.
//
// if you are concerned about the performance, usually markdown files
// are smaller than 3MB. So it would be pretty fast.
func (t *toc) reformatMarkdown(markdown string) (string, error) {
	// Get indexes of ending position of <!--toc-->
	// get the ending position of finish if exists.
	startIdx := startRe.FindStringIndex(markdown)
	finishIdx := finishRe.FindStringIndex(markdown)

	if startIdx == nil {
		return "", errors.New("ERROR: toc path is missing, add '<!--toc-->' to your markdown")
	}

	idx := startIdx[1] // end of <!--toc--> string

	finish := "<!-- tocstop -->" // default finish string
	if finishIdx != nil {
		finish = markdown[finishIdx[0]:finishIdx[1]]
		markdown = (markdown[:idx]) + markdown[finishIdx[1]:]
	}

	markdown = markdown[:idx] + "\n" + t.String() + "\n" + finish + markdown[idx:]

	return markdown, nil
}

func (t *toc) writeToFile(markdown string) error {
	markdown, err := t.reformatMarkdown(markdown)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(t.Options.Path, []byte(markdown), 0644)
	if err != nil {
		return err
	}

	return nil
}
