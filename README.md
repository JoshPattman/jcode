# `jcode` - A Simple 2D Robot Control Protocol
- JCODE is my own take on a 2D CNC language (similar to GCODE for 3D printers).
- It is designed to be extremely simple such that it can be easily re-implemented on many types of robot.
    - This simplicity is supposed to move as much of the maths of calculating paths to the powerful computer and away from the microcontroller.
- I have designed this primarily for a drawing robot I am working on.

## Instructions
- All instructions are followed by a semicolon.
- Whitespace before/after instructions should not be important (you can either put all instructions on one line or have a newline between each).
### Controller -> Robot
- `W <x> <y>;`: Set a waypoint to position `x,y`.
- `S <s>;`: Set the speed to move between waypoints to `s`.
- `D <d>;`: Stop where we are for `d` microseconds.
- `P <p>;`: Set the pen to either up `U` or down `D`.
- `H;`: Auto home the robot (and wait for homing to complete).
### Robot -> Controller
- `C;`: Sent to signal robot has just consumed an instruction.
- `L <msg>;`: Sent to indicate a log message.

## Example
The below example moves to a start position, waits for a second, then draws a line, finally returning to a start position.
```
S 2;
W 0 5;
D 1000;
P D;
W 5.5 5;
P U;
W 0 5;
```
