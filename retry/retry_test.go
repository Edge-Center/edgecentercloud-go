package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetryer(t *testing.T) {
	t.Run("5 fails", func(t *testing.T) {
		maxAttempt := 5
		backoff := NewBackoff(100*time.Millisecond, time.Second, maxAttempt)
		r := NewRetryer(backoff, nil)

		var attempts int
		retryerErr := errors.New("error")
		var elapsed []time.Duration

		start := time.Now()
		err := r.Run(context.Background(), func(ctx context.Context) error {
			attempts++
			elapsed = append(elapsed, time.Since(start))
			return retryerErr
		})

		assert.Equal(t, maxAttempt+1, attempts)
		assert.Equal(t, retryerErr, err)
		assert.Len(t, elapsed, 6)
	})

	t.Run("2 fails and succeed", func(t *testing.T) {
		maxAttempt := 5
		backoff := NewBackoff(10*time.Millisecond, time.Second, maxAttempt)
		r := NewRetryer(backoff, nil)

		var attempts int
		retryerErr := errors.New("error")
		var elapsed []time.Duration

		start := time.Now()
		err := r.Run(context.Background(), func(ctx context.Context) error {
			attempts++
			elapsed = append(elapsed, time.Since(start))

			if attempts == 2 {
				return nil
			}

			return retryerErr
		})

		assert.Equal(t, 2, attempts)
		assert.Nil(t, err)
		assert.Len(t, elapsed, 2)
	})
}
