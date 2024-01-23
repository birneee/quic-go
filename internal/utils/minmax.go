package utils

import (
	"golang.org/x/exp/constraints"
	"math"
	"time"
)

// InfDuration is a duration of infinite length
const InfDuration = time.Duration(math.MaxInt64)

// MinNonZeroDuration return the minimum duration that's not zero.
func MinNonZeroDuration(a, b time.Duration) time.Duration {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	return min(a, b)
}

// MinTime returns the earlier time
func MinTime(a, b time.Time) time.Time {
	if a.After(b) {
		return b
	}
	return a
}

// MinNonZeroTime returns the earlist time that is not time.Time{}
// If both a and b are time.Time{}, it returns time.Time{}
func MinNonZeroTime(a, b time.Time) time.Time {
	if a.IsZero() {
		return b
	}
	if b.IsZero() {
		return a
	}
	return MinTime(a, b)
}

// MaxTime returns the later time
func MaxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

// MaxV is a variadic version of Max
// returns the maximum
// panic if no arguments
func MaxV[T constraints.Ordered](nums ...T) T {
	max := nums[0]
	for _, num := range nums[1:] {
		if num > max {
			max = num
		}
	}
	return max
}
