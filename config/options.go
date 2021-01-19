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
    -p, --path <path>        Path for the markdown file.
    -h, --help               Show this message
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
	ShowHelp bool   `json:"show_help"`
}

// ConfigureOptions is a helper function for parsing options
func ConfigureOptions(fs *flag.FlagSet, args []string) (*Options, error) {
	opts := &Options{}

	// Define flags
	fs.BoolVar(&opts.ShowHelp, "h", false, "Show help message")
	fs.BoolVar(&opts.ShowHelp, "help", false, "Show help message")
	fs.BoolVar(&opts.Bulleted, "b", true, "")
	fs.BoolVar(&opts.Bulleted, "bulleted", true, "")
	fs.BoolVar(&opts.Bulleted, "a", true, "Append to markdown after <!--toc--> or write to stdout")
	fs.BoolVar(&opts.Bulleted, "append", true, "Append to markdown after <!--toc--> or write to stdout")
	fs.StringVar(&opts.Path, "p", "", "Path for the markdown file")
	fs.StringVar(&opts.Path, "path", "", "Path for the markdown file")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if (opts.ShowHelp == false) && (opts.Path == "") {
		err := errors.New("path flag is missing")
		return nil, err
	}

	return opts, nil
}
