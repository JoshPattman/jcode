package jcode

import (
	"math"
	"time"
)

// Calculates the distance between two waypoints
func Dist(a, b Waypoint) float64 {
	return math.Sqrt(math.Pow(a.XPos-b.XPos, 2) + math.Pow(a.YPos-b.YPos, 2))
}

// Calculates the time to travel between two waypoints at a given speed
func Time(a, b Waypoint, s Speed) time.Duration {
	return time.Microsecond * time.Duration(1000000*Dist(a, b)/s.Speed)
}

// Calculates if it is possible to move from a to b at the given speed
// (i.e. if the two waypoints are in different positions and the speed is 0, its impossible)
func Possible(a, b Waypoint, s Speed) bool {
	if a == b {
		return true
	}
	if s.Speed > 0 {
		return true
	}
	return false
}
