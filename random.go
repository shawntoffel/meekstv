package meekstv

import (
	"math"
)

// Implementation based on:
// "Algorithm AS 183: An Efficient and Portable Pseudo-Random Number Generator" by B. A. Wichmann and I. D. Hill.
type WichmannHillRandom interface {
	Next() float64
	NextInt(n int) int
}

type wichmannHillRandom struct {
	ix int
	iy int
	iz int
}

// Each seed should be a value between 1 and 30000.
func NewWichmannHillRandom(s1, s2, s3 int) WichmannHillRandom {
	return &wichmannHillRandom{
		ix: bound(s1),
		iy: bound(s2),
		iz: bound(s3),
	}
}

func (r *wichmannHillRandom) Next() float64 {
	// Use the modification noted in the paper for simpler logic, as integer arithmetic up to 5212632 is supported.
	r.ix = (171 * r.ix) % 30269
	r.iy = (172 * r.iy) % 30307
	r.iz = (170 * r.iz) % 30323

	random := float64(r.ix)/30269.0 + float64(r.iy)/30307.0 + float64(r.iz)/30323.0

	// Truncate to a value between 0 and 1.
	// math.Trunc is marginally faster here than math.Mod 1.0.
	return random - math.Trunc(random)
}

// Returns a pseudo-random int in the half-open interval [0,n).
func (r *wichmannHillRandom) NextInt(n int) int {
	return int(r.Next() * float64(n))
}

func bound(seed int) int {
	if seed < 1 {
		return 1
	}

	if seed > 30000 {
		return 30000
	}

	return seed
}
