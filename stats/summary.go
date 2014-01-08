package stats

import (
	"math"
)

type Summary struct {
	Mean, devsq, Min, Max float64
	Sum                   float64
	N                     int
}

func (s *Summary) Reset() {
	s.N = 0
	s.Max = 0.0
	s.Mean = 0.0
	s.Min = 0.
	s.Sum = 0.
	s.devsq = 0.0
}

// Add accumulates running statistics for calculating variance and
// standard deviation using the Welford method (1962)
func (s *Summary) Add(x float64) {
	if s.N > 0 {
		if x < s.Min {
			s.Min = x
		}
		if x > s.Max {
			s.Max = x
		}
	} else {
		s.Min, s.Max = x, x
	}

	s.N++
	s.Sum += x
	t := x - s.Mean
	s.Mean += t / float64(s.N)
	s.devsq += t * (x - s.Mean)

}

func (s *Summary) AddValues(x []float64) {
	for _, z := range x {
		s.Add(z)
	}
}

// Var returns the sample variance
func (s *Summary) Var() (v float64) {
	if s.N > 2 {
		v = s.devsq / (float64(s.N) - 1)
	}
	return
}

// Sd returns the sample standard deviation
func (s *Summary) Sd() (v float64) {
	if s.N > 2 {
		v = math.Sqrt(s.devsq / (float64(s.N) - 1))
	}
	return
}

// VarP returns the population variance
func (s *Summary) VarP() (v float64) {
	if s.N > 1 {
		v = s.devsq / float64(s.N)
	}
	return
}

// SdP returns the population standard deviation
func (s *Summary) SdP() (v float64) {
	if s.N > 1 {
		v = math.Sqrt(s.devsq / float64(s.N))
	}
	return
}

// Range returns the range of the data
func (s *Summary) Range() (r float64) {
	if s.N > 1 {
		r = s.Max - s.Min
	}
	return
}
