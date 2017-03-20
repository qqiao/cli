package cli

import (
	"bufio"
	"flag"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

// Component represents a command line component
type Component struct {
	// Components are the sub-components of the current component
	Components []*Component

	// Flag is a set of flags specific to this component
	Flag flag.FlagSet

	// Run runs the component
	// args are the arguments after the component name
	Run func(comp *Component, args []string)

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the component name
	UsageLine string

	// Short is the short description of the component
	Short string

	// Long is the longer more detailed description of the component
	Long string
}

// Name returns the name of the component: the first word in the UsageLine
func (c *Component) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// Runnable returns whether this component is runnable or pure informational
func (c *Component) Runnable() bool {
	return nil != c.Run
}

var usageTemplate = ` {{if .Runnable}}Usage: {{.UsageLine}}
{{end}}{{.Long | trim}}
{{if ne (len .Components) 0}}
The components are:
{{range .Components}}{{if .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}{{end}}
`

// Usage prints out the usage information
func (c *Component) Usage() {
	bw := bufio.NewWriter(os.Stdout)
	tmpl(bw, usageTemplate, c)
	bw.Flush()

	c.Flag.PrintDefaults()
}

func capitalize(s string) string {
	if "" == s {
		return s
	}

	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{
		"trim":       strings.TrimSpace,
		"capitalize": capitalize,
	})
	template.Must(t.Parse(text))
	t.Execute(w, data)
}
