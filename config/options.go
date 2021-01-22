package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

var usageStr = `
Usage: toc [options]
Options:
	-p, --path      <path>   Path for the markdown file.                               [REQUIRED]
	-a, --append    <bool>   Append toc after <!--toc-->, or write to stdout.          [Default: true]
	-b, --bulleted  <bool>   Write as bulleted, or write as numbered list.             [Default: true] 
	-s, --skip      <int>    Skip the first given number of headers.                   [Default: 0]
	-d, --depth     <int>    Set the number of maximum heading level to be included.   [Default: 6]
	-h, --help               Show this message and exit.
`

// UsageAndExit prints usage and exists the program.
func UsageAndExit(err error) {
	color.Red(err.Error())
	fmt.Println(usageStr)
	os.Exit(1)
}

// HelpAndExit , prints helps and exists the program.
func HelpAndExit() {
	fmt.Println(usageStr)
	os.Exit(0)
}

// Options is main value holder agentgo-server flags.
type Options struct {
	Path     string `json:"path"`
	Append   bool   `json:"append"`
	Bulleted bool   `json:"bulleted"`
	Skip     int    `json:"skip"`
	Depth    int    `json:"depth"`
	ShowHelp bool   `json:"show_help"`
}

// ConfigureOptions is a helper function for parsing options
func ConfigureOptions(fs *flag.FlagSet, args []string) (*Options, error) {
	opts := &Options{}

	// Define flags
	fs.StringVar(&opts.Path, "p", "", "Path for the markdown file")
	fs.StringVar(&opts.Path, "path", "", "Path for the markdown file")
	fs.BoolVar(&opts.Append, "a", true, "Append to markdown after <!--toc--> or write to stdout")
	fs.BoolVar(&opts.Append, "append", true, "Append to markdown after <!--toc--> or write to stdout")
	fs.BoolVar(&opts.Bulleted, "b", true, "Write as bulleted, or write as numbered list")
	fs.BoolVar(&opts.Bulleted, "bulleted", true, "Write as bulleted, or write as numbered list")
	fs.IntVar(&opts.Skip, "s", 0, "Skip the first given number of headers ")
	fs.IntVar(&opts.Skip, "skip", 0, "Skip the first given number of headers ")
	fs.IntVar(&opts.Depth, "d", 6, "Set the number of maximum heading level to be included")
	fs.IntVar(&opts.Depth, "depth", 6, "Set the number of maximum heading level to be included")
	fs.BoolVar(&opts.ShowHelp, "h", false, "Show help message")
	fs.BoolVar(&opts.ShowHelp, "help", false, "Show help message")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if (opts.ShowHelp == false) && (opts.Path == "") {
		err := errors.New("path flag is missing")
		return nil, err
	}

	return opts, nil
}
