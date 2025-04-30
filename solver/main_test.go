package solver

import (
	"testing"
)

func TestEvaluator(t *testing.T) {
	t.Run("Returns value from a shared range", func(t *testing.T) {
		tests := []struct {
			employee Range
			employer Range
		}{
			{
				employee: Range{0, 100},
				employer: Range{50, 150},
			},
			{
				employee: Range{0, 100},
				employer: Range{50, 80},
			},
			{
				employee: Range{50, 100},
				employer: Range{80, 80},
			},
		}

		for _, test := range tests {
			salary, err := Solve(test.employee, test.employer)

			assertError(t, err, nil)
			assertCorrectChoice(t, test.employee, test.employer, salary)
		}
	})

	t.Run("Returns ErrRangesDoNotOverlap error when salary is lower than requirement", func(t *testing.T) {
		employee := Range{200, 300}
		employer := Range{120, 150}

		_, err := Solve(employee, employer)

		assertError(t, err, ErrRangesDoNotOverlap)
	})

	t.Run("Returns the lowest salary if Max requirement is lower than Min salary", func(t *testing.T) {
		employee := Range{50, 100}
		employer := Range{200, 300}

		salary, err := Solve(employee, employer)

		assertCorrectChoice(t, employee, employer, salary)
		assertExactValue(t, salary, employer.Min)
		assertError(t, err, nil)
	})

	t.Run("Returns exact value if ranges 'touch' each other on employee Max", func(t *testing.T) {
		employee := Range{50, 100}
		employer := Range{100, 300}

		salary, err := Solve(employee, employer)

		assertCorrectChoice(t, employee, employer, salary)
		assertExactValue(t, salary, employer.Min)
		assertError(t, err, nil)
	})

	t.Run("Returns exact value if ranges 'touch' each other on employer Max", func(t *testing.T) {
		employee := Range{50, 100}
		employer := Range{40, 50}

		salary, err := Solve(employee, employer)

		assertCorrectChoice(t, employee, employer, salary)
		assertExactValue(t, salary, employee.Min)
		assertError(t, err, nil)
	})
}

func assertCorrectChoice(t testing.TB, employee, employer Range, salary int64) {
	if salary < employee.Min {
		t.Errorf("salary of %d is less than employee Min of %d", salary, employee.Min)
	}
	if salary > employer.Max {
		t.Errorf("salary of %d is more than employer Max of %d", salary, employer.Max)
	}
}

func assertExactValue(t testing.TB, got, want int64) {
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func assertError(t testing.TB, got, want error) {
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}
