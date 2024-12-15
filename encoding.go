package jcode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

// An encoder writes jcode instructions to an io.Writer.
type Encoder struct {
	w io.Writer
}

// Creates a new [Encoder].
func NewEncoder(w io.Writer) *Encoder {
	if w == nil {
		panic("cannot create jcode writer on nil value")
	}
	return &Encoder{
		w: w,
	}
}

// Writes a set of instructions to the writer.
func (j *Encoder) Write(instructions ...Instruction) error {
	for _, ins := range instructions {
		s := ""
		switch ins := ins.(type) {
		case Waypoint:
			s = fmt.Sprintf("W %.3f %.3f;", ins.XPos, ins.YPos)
		case Speed:
			s = fmt.Sprintf("S %.3f;", ins.Speed)
		case Delay:
			s = fmt.Sprintf("D %d;", ins.Duration.Milliseconds())
		case Pen:
			s = fmt.Sprintf("P %s;", map[PenMode]string{PenDown: "D", PenUp: "U"}[ins.Mode])
		case Consumed:
			s = "C;"
		case Log:
			s = fmt.Sprintf("L %s;", ins.Message)
		case AutoHome:
			s = "H;"
		default:
			return fmt.Errorf("invalid instruction of type %T", ins)
		}
		_, err := j.w.Write([]byte(s))
		if err != nil {
			return err
		}
	}
	return nil
}

// A dencoder reads jcode instructions from a io.Reader.
// This will read the buffer even when you are not actively using the Decoder, so be careful!
type Decoder struct {
	r *bufio.Reader
}

// Creates a new [Decoder].
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: bufio.NewReader(r),
	}
}

// Reads an instruction from the reader. This is blocking.
func (j *Decoder) Read() (Instruction, error) {
	s, err := j.r.ReadString(';')
	s = strings.Trim(s, " \n\r\t;")
	if err != nil {
		return nil, err
	}
	if len(s) == 0 {
		return nil, fmt.Errorf("empty instruction")
	}
	ins := s[0]
	s = strings.Trim(s[1:], " ")
	switch ins {
	case 'W':
		w := Waypoint{}
		if _, err := fmt.Sscanf(s, "%f %f", &w.XPos, &w.YPos); err != nil {
			return nil, err
		}
		return w, nil
	case 'S':
		w := Speed{}
		if _, err := fmt.Sscanf(s, "%f", &w.Speed); err != nil {
			return nil, err
		}
		return w, nil
	case 'D':
		millis := time.Duration(0)
		if _, err := fmt.Sscanf(s, "%d", &millis); err != nil {
			return nil, err
		}
		return Delay{time.Millisecond * millis}, nil
	case 'P':
		mode := ""
		if _, err := fmt.Sscanf(s, "%s", &mode); err != nil {
			return nil, err
		}
		p := Pen{}
		switch mode {
		case "U":
			p.Mode = PenUp
		case "D":
			p.Mode = PenDown
		default:
			return nil, fmt.Errorf("invalid pen mode '%s'", mode)
		}
		return p, nil
	case 'C':
		return Consumed{}, nil
	case 'L':
		return Log{s}, nil
	case 'H':
		return AutoHome{}, nil
	default:
		return nil, fmt.Errorf("invalid instruction '%v'", ins)
	}
}

// Helper function to start a seperate goroutine that reads instructions from the reader.
func BeginInstructionProcessing(r io.Reader, bufSize int) chan Instruction {
	instructionReader := NewDecoder(r)
	result := make(chan Instruction, bufSize)
	go func() {
		for {
			ins, err := instructionReader.Read()
			if errors.Is(err, io.EOF) {
				close(result)
				return
			} else if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			result <- ins
		}
	}()
	return result
}
