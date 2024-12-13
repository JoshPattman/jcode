package jcode

import (
	"bytes"
	"math"
	"os"
	"reflect"
	"testing"
	"time"
)

// Not a real test, just used for debugging
func TestJCodeStdout(t *testing.T) {
	enc := NewEncoder(os.Stdout)

	ins := []Instruction{
		Pen{PenUp},
		Speed{10},
		Waypoint{0, 0},
		Pen{PenDown},
		Waypoint{0, 10},
		Delay{time.Second},
	}

	if err := enc.Write(ins...); err != nil {
		panic(err)
	}
}

// Not a real test, just used for debugging
func TestJCodeGenStdout(t *testing.T) {
	enc := NewEncoder(os.Stdout)

	curve := CircleCurve{
		Center: Waypoint{YPos: 5},
		Radius: 1,
		Speed:  1,
	}

	if err := enc.Write(ExportCurve(curve, 10)...); err != nil {
		panic(err)
	}
}

// Test encoding and decoding instructions
func TestJCodeEncodeDecode(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	enc := NewEncoder(buf)
	dec := NewDecoder(buf)

	curve := CircleCurve{
		Radius: 1,
		Speed:  1,
	}

	input := ExportCurve(curve, 10)

	if err := enc.Write(input...); err != nil {
		panic(err)
	}

	for i := range input {
		ins, err := dec.Read()
		if err != nil {
			t.Fatalf("%v", err)
		}
		if reflect.TypeOf(input[i]) != reflect.TypeOf(ins) || !equal(ins, input[i]) {
			t.Fatalf("input %v did not match output %v (type %T and %T)", input[i], ins, input[i], ins)
		}
	}
}

func equal(i1, i2 Instruction) bool {
	delta := 0.001
	switch i1 := i1.(type) {
	case Waypoint:
		i2 := i2.(Waypoint)
		return floatDelta(i1.XPos, i2.XPos, delta) && floatDelta(i1.YPos, i2.YPos, delta)
	case Speed:
		i2 := i2.(Speed)
		return floatDelta(i1.Speed, i2.Speed, delta)
	case Delay, Pen:
		return i1 == i2
	default:
		panic("no")
	}
}

func floatDelta(x, y, d float64) bool {
	return math.Abs(x-y) < d
}
