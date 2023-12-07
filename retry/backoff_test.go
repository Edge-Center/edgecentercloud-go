package retry

import (
	"math"
	"testing"
	"time"
)

func TestExponentialBackoff(t *testing.T) {
	tests := map[string]struct {
		attemptNum int
		min, max   time.Duration

		wantMin, wantMax time.Duration
	}{
		"1": {
			attemptNum: 1,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    time.Duration(math.Pow(2, float64(1))),
			wantMax:    10 * time.Second,
		},
		"2": {
			attemptNum: 2,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    time.Duration(math.Pow(2, float64(2))),
			wantMax:    10 * time.Second,
		},
		"3": {
			attemptNum: 3,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    time.Duration(math.Pow(2, float64(3))),
			wantMax:    10 * time.Second,
		},
		"4": {
			attemptNum: 4,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    time.Duration(math.Pow(2, float64(4))),
			wantMax:    10 * time.Second,
		},
		"over": {
			attemptNum: 11,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    10 * time.Second,
			wantMax:    10 * time.Second,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := ExponentialBackoff(tc.attemptNum, tc.min, tc.max)
			between(t, actual, tc.wantMin, tc.wantMax)
		})
	}
}

func between(t *testing.T, actual, low, high time.Duration) {
	t.Helper()
	if actual < low {
		t.Fatalf("Got %s, Expecting >= %s", actual, low)
	}
	if actual > high {
		t.Fatalf("Got %s, Expecting <= %s", actual, high)
	}
}
