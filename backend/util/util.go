package util

import (
	"context"
	"math/rand"
	"time"
)

// ExponentialBackoff: the more times in a row something fails, the longer we wait.
type ExponentialBackoff struct {
	Factor float64
	curr   *time.Duration
	Min    time.Duration
	Max    time.Duration
	Jitter bool
}

func (e *ExponentialBackoff) Reset() {
	e.curr = &e.Min
}

/* #nosec */
func (e *ExponentialBackoff) DelayOnFail(ctx context.Context) {

	if e.curr == nil {
		e.curr = &e.Min
	} else {
		newDurationVal := time.Duration(float64(e.Factor) * float64(*e.curr))

		if e.Jitter {
			// Randomly increase the result by up to 10%, to introduce a small amount of jitter
			newDurationVal = time.Duration(float64(newDurationVal) * (1 + (rand.Float64() * float64(0.10))))
		}

		e.curr = &newDurationVal
	}

	if *e.curr < e.Min {
		e.curr = &e.Min
	}

	if *e.curr > e.Max {
		e.curr = &e.Max
	}
	select {
	case <-ctx.Done():
	case <-time.After(*e.curr):
	}

}