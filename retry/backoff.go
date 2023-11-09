package retry

import (
	"math"
	"math/rand"
	"time"
)

type FnBackoff func(attemptNum int, min, max time.Duration) time.Duration

type Backoff struct {
	min, max   time.Duration
	maxAttempt int
	attemptNum int
	backoff    FnBackoff
}

func NewBackoff(min, max time.Duration, maxAttempt int) *Backoff {
	return &Backoff{
		min:        min,
		max:        max,
		maxAttempt: maxAttempt,
		backoff:    ExponentialBackoff,
	}
}

const Stop time.Duration = -1

func (b *Backoff) Next() time.Duration {
	if b.attemptNum >= b.maxAttempt {
		return Stop
	}
	b.attemptNum++

	return b.backoff(b.attemptNum, b.min, b.max)
}

func (b *Backoff) Reset() {
	b.attemptNum = 0
}

func ExponentialBackoff(attemptNum int, min, max time.Duration) time.Duration {
	factor := 2.0
	delay := time.Duration(math.Pow(factor, float64(attemptNum)) * float64(min))
	jitter := time.Duration(rand.Float64() * float64(min) * float64(attemptNum))

	delay += jitter
	if delay > max {
		delay = max
	}

	return delay
}
