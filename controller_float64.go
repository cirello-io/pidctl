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
	"math/big"
	"time"
)

// ControllerFloat64 implements a PID controller using float64 for inputs and
// outputs.
type ControllerFloat64 struct {
	controller Controller
}

// NewControllerFloat64 creates a new PID controller using float64 for inputs
// and outputs.
func NewControllerFloat64(p, i, d, setpoint float64) *ControllerFloat64 {
	c := &ControllerFloat64{
		controller: Controller{
			P:        ratFloat64(p),
			I:        ratFloat64(i),
			D:        ratFloat64(d),
			Setpoint: ratFloat64(setpoint),
		},
	}
	c.controller.init()
	return c
}

// SetSetpoint changes the desired setpoint of the controller.
func (c *ControllerFloat64) SetSetpoint(setpoint float64) *ControllerFloat64 {
	c.controller.Setpoint = ratFloat64(setpoint)
	return c
}

// SetMin changes the minimum output value of the controller.
func (c *ControllerFloat64) SetMin(min float64) *ControllerFloat64 {
	c.controller.Min = ratFloat64(min)
	return c
}

// SetMax changes the maxium output value of the controller.
func (c *ControllerFloat64) SetMax(max float64) *ControllerFloat64 {
	c.controller.Max = ratFloat64(max)
	return c
}

// Compute updates the controller with the given process value since the last
// update. It returns the new output that should be used by the device to reach
// the desired set point. Internally it assumes the duration between calls is
// constant.
func (c *ControllerFloat64) Compute(pv float64) float64 {
	v, _ := c.controller.Compute(ratFloat64(pv)).Float64()
	return v
}

// Accumulate updates the controller with the given process value and duration
// since the last update. It returns the new output that should be used by the
// device to reach the desired set point.
func (c *ControllerFloat64) Accumulate(pv float64, deltaTime time.Duration) float64 {
	v, _ := c.controller.Accumulate(ratFloat64(pv), deltaTime).Float64()
	return v
}

func ratFloat64(i float64) *big.Rat {
	return new(big.Rat).SetFloat64(i)
}
