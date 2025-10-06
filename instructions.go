package jcode

import "time"

type Instruction interface {
	jcode()
}

func (w Waypoint) jcode() {}
func (w Speed) jcode()    {}
func (w Delay) jcode()    {}
func (w Pen) jcode()      {}
func (c Consumed) jcode() {}
func (e Log) jcode()      {}
func (a AutoHome) jcode() {}

type Waypoint struct {
	XPos, YPos float64
}

type Speed struct {
	Speed float64
}

type Delay struct {
	Duration time.Duration
}

type PenMode bool

const (
	PenUp   PenMode = true
	PenDown PenMode = false
)

type Pen struct {
	Mode PenMode
}

type Consumed struct{}

type Log struct {
	Message string
}

type AutoHome struct{}
