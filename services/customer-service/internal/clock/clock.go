// Package clock provides time-related functionality with interfaces and implementations
// for managing time in applications. It includes both real-time and fixed-time clock
// implementations, allowing for flexible time handling in both production and testing
// environments.

package clock

import "time"

// Clock defines an interface for obtaining the current time.
type Clock interface {
	Now() time.Time
}

// RealClock is a concrete implementation of the Clock interface that provides the current system time.
type RealClock struct{}

// Now returns the current local time provided by the system clock.
func (RealClock) Now() time.Time {
	return time.Now()
}

// FixedClock represents a clock that always returns a fixed, predefined time.
type FixedClock struct {
	FixedTime time.Time
}

// Now returns the fixed predefined time set in the FixedClock instance.
func (c FixedClock) Now() time.Time {
	return c.FixedTime
}
