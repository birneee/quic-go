package indi_utils

import "golang.org/x/exp/constraints"

func Base10Digits[T constraints.Integer](i T) int {
	if i == 0 {
		return 1
	}
	count := 0
	for i != 0 {
		i /= 10
		count++
	}
	return count
}
