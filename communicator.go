package jcode

import (
	"errors"
	"io"
	"sync"
)

type Communicator struct {
	encoder   *Encoder
	decoder   *Decoder
	running   *sync.Mutex
	toRobot   chan Instruction
	fromRobot chan Instruction
	err       chan error
	maxBuffer int
}

func NewCommunicator(send io.Writer, recv io.Reader, maxBuffer int) *Communicator {
	return &Communicator{
		NewEncoder(send),
		NewDecoder(recv),
		&sync.Mutex{},
		make(chan Instruction),
		make(chan Instruction),
		make(chan error),
		maxBuffer,
	}
}

// Send instructions here to get sent to the robot,
// may block if the robot is currently at maximum instruction capacity.
func (com *Communicator) ToRobot() chan Instruction {
	return com.toRobot
}

// Recive instructions back from the robot.
// This also includes the "consumed" messages.
func (com *Communicator) FromRobot() chan Instruction {
	return com.fromRobot
}

// The channel of which errors will be returned.
// Once an error is returned here, the communicator will stop running.
func (com *Communicator) Error() chan error {
	return com.err
}

func (com *Communicator) Start() {
	// Manage the running state (can only run once and never again)
	if !com.running.TryLock() {
		com.err <- errors.New("communicator has already been run (can only use once)")
		return
	}

	// Check parameters
	if com.maxBuffer <= 0 {
		com.err <- errors.New("max buffer must be a positive non-zero int")
		return
	}

	// Setup intermediate channel from robot for this run
	intFromRobot := make(chan Instruction)

	// Start async reading instructions forever (until there is an error)
	go func() {
		for {
			ins, err := com.decoder.Read()
			if err != nil {
				com.err <- err
				return
			}
			intFromRobot <- ins
		}
	}()

	// Start the main communicator goroutine
	go func() {
		// This is how many instructions we have sent to the robot without it consuming
		robotBufCapacity := 0
		for {
			if robotBufCapacity < com.maxBuffer {
				// We can send stuff to the robot as it has free space
				select {
				case ins := <-com.toRobot:
					err := com.encoder.Write(ins)
					if err != nil {
						com.err <- err
						return
					}
					robotBufCapacity += 1
				case ins := <-com.fromRobot:
					if _, ok := ins.(Consumed); ok {
						// TODO: Repeated code
						// If robot says it consumed an instruction, remeber that the robot has one less
						robotBufCapacity -= 1
					}
					com.fromRobot <- ins
				}
			} else {
				// The robot is currently full, process instructions
				ins := <-intFromRobot
				if _, ok := ins.(Consumed); ok {
					// If robot says it consumed an instruction, remeber that the robot has one less
					robotBufCapacity -= 1
				}
				com.fromRobot <- ins
			}
		}
	}()
}
