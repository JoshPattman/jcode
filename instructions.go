package jcode

import "time"

type Instruction interface {
	jcode()
}

type Waypoint struct {
	XPos, YPos float64
}

func (w Waypoint) jcode() {}

type Speed struct {
	Speed float64
}

func (w Speed) jcode() {}

type Delay struct {
	Duration time.Duration
}

func (w Delay) jcode() {}

type PenMode bool

const (
	PenUp   PenMode = true
	PenDown PenMode = false
)

type Pen struct {
	Mode PenMode
}

func (w Pen) jcode() {}

type Consumed struct{}

func (c Consumed) jcode() {}

type Log struct {
	Message string
}

func (e Log) jcode() {}

type AutoHome struct{}

func (a AutoHome) jcode() {}
