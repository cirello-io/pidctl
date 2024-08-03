/*
Copyright 2024 github.com/ucirello and cirello.io

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pidctl

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
	"text/tabwriter"
)

func TestControllerFloat64(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				r := recover()
				switch {
				case r == nil:
				case test.expectPanic && fmt.Sprintf("%T", r) == "*pidctl.MinMaxError":
					t.Log("trapped:", r)
				default:
					panic(r)
				}
			}()
			c := NewControllerFloat64(test.p, test.i, test.d, 0)

			if test.min != 0 || test.max != 0 {
				c.SetMin(test.min)
				c.SetMax(test.max)
			}

			var buf bytes.Buffer
			log := tabwriter.NewWriter(&buf, 8, 0, 1, ' ', 0)
			fmt.Fprint(log, "\tcycle\tgot\texpected\tsetpoint\tinput\toutput\n")
			for i, u := range test.steps {
				if u.setpoint != 0 {
					c.SetSetpoint(u.setpoint)
				}
				got := c.Accumulate(u.input, test.stepDuration)
				roundedGot, roundedExpected := fmt.Sprintf("%0.3f", got), fmt.Sprintf("%0.3f", u.output)
				msg := ""
				if roundedGot != roundedExpected {
					msg = "error"
					t.Fail()
				}
				fmt.Fprintf(log, "%s\t%d\t%v\t%v\t%v\t%v\t%v\n", msg, i, roundedGot, roundedExpected, u.setpoint, u.input, u.output)
			}
			log.Flush()
			scanner := bufio.NewScanner(&buf)
			for scanner.Scan() {
				t.Log(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				panic(fmt.Sprint("reading table output:", err))
			}
		})
	}
}
