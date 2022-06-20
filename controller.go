/*
Copyright 2019 github.com/ucirello and cirello.io

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
	"sync"
	"time"
)

var ratZero = new(big.Rat).SetInt64(0)

// Controller implements a PID controller.
type Controller struct {
	// P is the proportional gain
	P *big.Rat
	// I is integral reset
	I *big.Rat
	// D is the derivative term
	D *big.Rat
	// current setpoint
	Setpoint *big.Rat
	// Min the lowest value acceptable for the Output
	Min *big.Rat
	// Max the highest value acceptable for the Output
	Max *big.Rat

	prevProcessVariable *big.Rat
	accumulatedIntegral *big.Rat
	initOnce            sync.Once
}

func (p *Controller) init() {
	p.initOnce.Do(func() {
		if p.P == nil {
			p.P = new(big.Rat).SetInt64(0)
		}
		if p.I == nil {
			p.I = new(big.Rat).SetInt64(0)
		}
		if p.D == nil {
			p.D = new(big.Rat).SetInt64(0)
		}
		if p.Setpoint == nil {
			p.Setpoint = new(big.Rat).SetInt64(0)
		}
		if p.prevProcessVariable == nil {
			p.prevProcessVariable = new(big.Rat).SetInt64(0)
		}
		if p.accumulatedIntegral == nil {
			p.accumulatedIntegral = new(big.Rat).SetInt64(0)
		}
	})
}

// Accumulate updates the controller with the given value and duration since the
// last update. It returns the new output that should be used by the device to
// reach the desired set point.
func (p *Controller) Accumulate(v *big.Rat, duration time.Duration) *big.Rat {
	p.init()
	var (
		processVariable     = v.Set(v)
		dt                  = new(big.Rat)
		err                 = new(big.Rat)
		proportional        = new(big.Rat)
		integral            = new(big.Rat)
		derivative          = new(big.Rat)
		output              = new(big.Rat)
		prevProcessVariable *big.Rat
	)
	prevProcessVariable, p.prevProcessVariable = p.prevProcessVariable, processVariable
	dt.SetInt64(int64(duration / time.Second))
	err.Sub(p.Setpoint, processVariable)

	// Proportional Gain
	proportional.Mul(p.P, err)
	output.Add(output, proportional)

	// Integral Reset
	integral.
		Mul(err, dt).
		Mul(integral, p.I).
		Add(integral, p.accumulatedIntegral)
	p.accumulatedIntegral = p.enforceRange(integral) // avoid integral windup
	output.Add(output, p.accumulatedIntegral)

	// Derivative Term
	if dt.Cmp(ratZero) > 0 {
		derivative.Sub(processVariable, prevProcessVariable)
		derivative.Quo(derivative, dt)
		derivative.Neg(derivative)
	}
	derivative.Mul(derivative, p.D)
	output.Add(output, derivative)

	return p.enforceRange(output)
}

func (p *Controller) enforceRange(v *big.Rat) *big.Rat {
	if p.Max != nil && v.Cmp(p.Max) > 0 {
		return p.Max
	} else if p.Min != nil && v.Cmp(p.Min) < 0 {
		return p.Min
	}
	return v
}
