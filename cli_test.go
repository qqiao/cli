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

package cli

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
)

const UsageLine = `test [-i input]`
const Long = `Long usage line for the application designed to test formatting.`

func TestName(t *testing.T) {
	c := Component{
		UsageLine: UsageLine,
	}

	if "test" != c.Name() {
		t.Errorf("Expected '%s', got '%s'", "test", c.Name())
	}
}

func TestRunnable(t *testing.T) {
	r := Component{
		Run: func(context.Context, *Component, []string) {},
	}
	if !r.Runnable() {
		t.Errorf("Expected '%t', got '%t'", true, r.Runnable())
	}

	nr := Component{}
	if nr.Runnable() {
		t.Errorf("Expected '%t', got '%t'", false, nr.Runnable())
	}
}

func TestUsageFlags(t *testing.T) {
	var buf bytes.Buffer

	c := Component{
		UsageLine: UsageLine,
		Run:       func(context.Context, *Component, []string) {},
	}
	c.SetOutput(&buf)
	c.FlagSet().String("i", "", "input of the test component")
	c.Usage()

	expected := `Usage: test [-i input]

The flags are:
  -i string
    	input of the test component
`

	if buf.String() != expected {
		t.Errorf("Expected '%s'. got '%s'", expected, buf.String())
	}
}

func TestUsageFlagsWithLong(t *testing.T) {
	var buf bytes.Buffer

	c := Component{
		UsageLine: UsageLine,
		Long:      Long,
		Run:       func(context.Context, *Component, []string) {},
	}
	c.SetOutput(&buf)
	c.FlagSet().String("i", "", "input of the test component")
	c.Usage()

	expected := `Usage: test [-i input]
Long usage line for the application designed to test formatting.

The flags are:
  -i string
    	input of the test component
`

	if buf.String() != expected {
		t.Errorf("Expected '%s'. got '%s'", expected, buf.String())
	}
}

func TestUsageRunnable(t *testing.T) {
	expectedUsageLine := fmt.Sprintf("Usage: %s", UsageLine)

	var buf bytes.Buffer
	c := Component{
		UsageLine: UsageLine,
		Long:      "This is the long description of the test component.",
	}
	c.SetOutput(&buf)

	c.Usage()
	usage := buf.String()
	if strings.HasPrefix(usage, expectedUsageLine) {
		t.Error("Non-runnable component shouldn't have a usage line")
	}

	buf.Reset()
	c.Run = func(context.Context, *Component, []string) {}
	c.Usage()
	usage = buf.String()
	if !strings.HasPrefix(usage, expectedUsageLine) {
		t.Error("Usage line missing")
	}
}

func TestUsageSubComponent(t *testing.T) {
	var buf bytes.Buffer

	c := Component{
		UsageLine: UsageLine,
		Long:      Long,
		Run:       func(context.Context, *Component, []string) {},
		Components: []*Component{
			&Component{
				UsageLine: "subcomponent1",
				Short:     "description of subcomponent 1",
				Run:       func(context.Context, *Component, []string) {},
			},
			&Component{
				UsageLine: "subcomponent2",
				Short:     "description of subcomponent 2",
				Run:       func(context.Context, *Component, []string) {},
			},
		},
	}
	c.SetOutput(&buf)
	c.FlagSet().String("i", "", "input of the test component")
	c.Usage()

	expected := `Usage: test [-i input]
Long usage line for the application designed to test formatting.

The components are:
  subcomponent1 description of subcomponent 1
  subcomponent2 description of subcomponent 2

The flags are:
  -i string
    	input of the test component
`

	if buf.String() != expected {
		t.Errorf("Expected '%s'. got '%s'", expected, buf.String())
	}
}
