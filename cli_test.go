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
	"testing"
)

const UsageLine = `test [-i input]`
const Long = `Long usage line for the application designed to test formatting.`

func TestComponent_Runnable(t *testing.T) {
	tests := []struct {
		name string
		c    *Component
		want bool
	}{
		{
			name: "Not Runnable",
			c:    &Component{},
			want: false,
		},
		{
			name: "Runnable",
			c: &Component{
				Run: Passthrough,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Runnable(); got != tt.want {
				t.Errorf("Component.Runnable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComponent_Name(t *testing.T) {
	tests := []struct {
		name string
		c    *Component
		want string
	}{
		{
			name: "test",
			c: &Component{
				UsageLine: UsageLine,
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Name(); got != tt.want {
				t.Errorf("Component.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComponent_Usage(t *testing.T) {
	tests := []struct {
		name string
		c    *Component
		want string
	}{
		{
			name: "Without Flags",
			c: &Component{
				UsageLine: UsageLine,
				Run:       Passthrough,
			},
			want: `Usage: test [-i input]
`,
		},
		{
			name: "With Flags",
			c: &Component{
				UsageLine: UsageLine,
				Run:       Passthrough,
			},
			want: `Usage: test [-i input]

The flags are:
  -i string
    	input of the test component
`,
		},
		{
			name: "Non Runnable",
			c: &Component{
				UsageLine: UsageLine,
			},
			want: `
The flags are:
  -i string
    	input of the test component
`,
		},
		{
			name: "Runnable",
			c: &Component{
				UsageLine: UsageLine,
				Run:       Passthrough,
			},
			want: `Usage: test [-i input]

The flags are:
  -i string
    	input of the test component
`,
		},
		{
			name: "With Flag and Long",
			c: &Component{
				UsageLine: UsageLine,
				Long:      Long,
				Run:       Passthrough,
			},
			want: `Usage: test [-i input]
Long usage line for the application designed to test formatting.

The flags are:
  -i string
    	input of the test component
`,
		},
		{
			name: "Subcomponent",
			c: &Component{
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
			},
			want: `Usage: test [-i input]
Long usage line for the application designed to test formatting.

The components are:
  subcomponent1 description of subcomponent 1
  subcomponent2 description of subcomponent 2

The flags are:
  -i string
    	input of the test component
`,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.c.SetOutput(&buf)
			if i != 0 {
				tt.c.FlagSet().String("i", "", "input of the test component")
			}
			tt.c.Usage()

			if got := buf.String(); got != tt.want {
				t.Errorf("Component.Usage() = %v, want %v", got, tt.want)
			}
		})
	}
}
