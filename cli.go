// Copyright (c) 2017 Qian Qiao
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cli provides base constructs to build complex command line
// applications.
//
// By nesting components, one can build command line applications like: the
// go command.
//
package cli

import (
	"bytes"
	"context"
	"flag"
	"io"
	"strings"
	"text/template"
)

// Component represents a command line component
type Component struct {
	// Components are the sub-components of the current component
	Components []*Component

	// Run runs the component
	// args are the arguments after the component name
	Run func(ctx context.Context, comp *Component, args []string)

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the component name
	UsageLine string

	// Short is the short description of the component
	Short string

	// Long is the longer more detailed description of the component
	Long string

	// flagSet is a set of flags specific to this component
	flagSet *flag.FlagSet
}

// FlagSet returns the set of command line flags
func (c *Component) FlagSet() *flag.FlagSet {
	if nil == c.flagSet {
		c.flagSet = flag.NewFlagSet(c.Name(), flag.ExitOnError)
		c.flagSet.Usage = c.Usage
	}

	return c.flagSet
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

// SetOutput sets the destination for usage messages.
// If output is nil, stderr is used
func (c *Component) SetOutput(output io.Writer) {
	c.FlagSet().SetOutput(output)
}

var usageTemplate = `{{if .component.Runnable}}Usage: {{.component.UsageLine}}{{end}}
{{- if ne (len .component.Long) 0}}{{.component.Long | trim}}{{end}}
{{- if ne (len .component.Components) 0}}
The components are:
{{range .component.Components}}{{if .Runnable}}
  {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}{{end}}
{{if ne (len .flags) 0}}
The flags are:
{{.flags}}{{end}}`

// Usage prints out the usage information
func (c *Component) Usage() {
	output := c.flagSet.Output()

	buf := bytes.NewBufferString("")
	c.flagSet.SetOutput(buf)
	c.flagSet.PrintDefaults()

	c.flagSet.SetOutput(output)

	tmpl(output, usageTemplate, map[string]interface{}{
		"component": c,
		"flags":     buf.String(),
	})
}

// Passthrough is a implementation of the Run function that passes the
// execution through the sub commands
func Passthrough(ctx context.Context, component *Component, args []string) {
	if flag.ErrHelp == component.FlagSet().Parse(args) {
		return
	}

	if component.FlagSet().NArg() < 1 {
		component.FlagSet().Usage()
		return
	}

	name := component.FlagSet().Arg(0)

	for _, comp := range component.Components {
		if name == comp.Name() {
			if comp.Runnable() {
				comp.Run(ctx, comp, component.FlagSet().Args()[1:])
				return
			}
		}
	}
	component.FlagSet().Usage()
}

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{
		"trim": strings.TrimSpace,
	})
	template.Must(t.Parse(text))
	t.Execute(w, data)
}
