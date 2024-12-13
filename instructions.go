package jcode

import "time"

type Instruction interface {
	JCode()
}

type Waypoint struct {
	XPos, YPos float64
}

func (w Waypoint) JCode() {}

type Speed struct {
	Speed float64
}

func (w Speed) JCode() {}

type Delay struct {
	Duration time.Duration
}

func (w Delay) JCode() {}

type PenMode bool

const (
	PenUp   PenMode = true
	PenDown PenMode = false
)

type Pen struct {
	Mode PenMode
}

func (w Pen) JCode() {}
