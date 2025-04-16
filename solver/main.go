package solver

import (
	"errors"
	"math/rand"
)

type Range struct {
	Min int
	Max int
}

var (
	ErrRangesDoNotOverlap = errors.New("ranges do not overlap")
)

func Solve(employee, employer Range) (int, error) {
	if employer.Max < employee.Min {
		return 0, ErrRangesDoNotOverlap
	}

	if employee.Max < employer.Min {
		return employer.Min, nil
	}

	if employee.Max == employer.Min {
		return employer.Min, nil
	}

	if employee.Min == employer.Max {
		return employee.Min, nil
	}

	min := Max(employee.Min, employer.Min)
	max := Min(employee.Max, employer.Max)

	salary := rand.Intn(max-min+1) + min

	return salary, nil
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
