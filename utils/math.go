package utils

import (
	"sort"
)

// Min returns the lowest value from the provided parameters.
func Min[T Number](values ...T) T {
	var acc T = values[0]

	for _, v := range values {
		if v < acc {
			acc = v
		}
	}
	return acc
}

// Max returns the biggest value from the provided parameters.
func Max[T Number](values ...T) T {
	var acc T = values[0]

	for _, v := range values {
		if v > acc {
			acc = v
		}
	}
	return acc
}

// Abs returns the absolut value of x.
func Abs[T Number](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Clamp returns a range-limited number between min and max.
func Clamp[T Number](num, min, max T) T {
	if num <= min {
		return min
	} else if num >= max {
		return max
	}
	return num
}

// InRange checks if a number is inside a range.
func InRange[T Number](num, lo, up T) bool {
	if num >= lo && num <= up {
		return true
	}
	return false
}

func Avg[T Number](n ...T) (avg T) {
	return Sum(n...) / T(len(n))
}

// AvgF64 calculates average with final division using float64 type
// int overflow can still happen
func AvgF64[T Number](n ...T) (avg float64) {
	return float64(Sum(n...)) / float64(len(n))
}

// AvgF64F64 calculates average after converting any input to float64
// to avoid integer overflows
func AvgF64F64[T Number](n ...T) (avg float64) {
	for _, v := range n {
		avg = avg + float64(v)
	}
	avg = avg / float64(len(n))
	return avg
}

func Median[T Number](n ...T) (median T) {
	// TODO probably want math-specific sort here, not to call function every time.
	sort.Slice(n, func(x, y int) bool { return n[x] < n[y] })
	if (len(n) % 2) == 0 {
		return ((n[len(n)/2-1]) + (n[(len(n) / 2)])) / 2
	} else {
		return n[len(n)/2]
	}
}

// MedianF64 calculates median with final division using float64 type
func MedianF64[T Number](n ...T) (median float64) {
	// TODO probably want math-specific sort here, not to call function every time.
	sort.Slice(n, func(x, y int) bool { return n[x] < n[y] })
	if (len(n) % 2) == 0 {
		return (float64(n[len(n)/2-1]) + float64(n[(len(n)/2)])) / 2.0
	} else {
		return float64(n[len(n)/2])
	}
}
