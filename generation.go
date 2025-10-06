package jcode

import (
	"math"
	"time"
)

type Curve interface {
	Evaluate(at time.Duration) (Waypoint, Speed)
	Duration() time.Duration
}

type CircleCurve struct {
	Center Waypoint
	Radius float64
	Speed  float64
}

func (c CircleCurve) Evaluate(at time.Duration) (Waypoint, Speed) {
	time := c.Duration().Seconds()
	t := at.Seconds() / time * math.Pi * 2
	return Waypoint{math.Cos(t)*c.Radius + c.Center.XPos, math.Sin(t)*c.Radius + c.Center.YPos}, Speed{c.Speed}
}

func (c CircleCurve) Duration() time.Duration {
	return time.Microsecond * time.Duration(math.Pi*2*c.Radius/c.Speed*1000000)
}

func ExportCurve(c Curve, pointPerSecond float64) []Instruction {
	code := make([]Instruction, 0)
	t := 0.0
	lastSpeed := Speed{-834530943289}
	step := 1.0 / pointPerSecond
	for t < c.Duration().Seconds()+step {
		tc := min(t, c.Duration().Seconds())
		wp, sp := c.Evaluate(time.Microsecond * time.Duration(tc*1000000))
		if sp != lastSpeed {
			code = append(code, sp, wp)
			lastSpeed = sp
		} else {
			code = append(code, wp)
		}
		t += step
	}
	return code
}
