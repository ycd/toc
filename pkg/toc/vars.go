package toc

var headers = map[string]int{"h1": 0, "h2": 1, "h3": 2, "h4": 3, "h5": 4, "h6": 5}

// Use 4 spaces for adding tabs
// See Markdown specification
// https://github.github.com/gfm/#tabs
var tab = "    "

type tocConfig struct {
	Path     string
	Bulleted bool
	Append   bool
	Skip     int
	Depth    int
}

type toc struct {
	Options tocConfig
	Content []string
}
