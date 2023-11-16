package util

import "github.com/avast/retry-go/v4"

const Attempts uint = 10

func WithRetry(f retry.RetryableFunc, attempts *uint) error {
	if attempts == nil {
		temp := Attempts
		attempts = &temp
	}

	return retry.Do(
		f,
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(*attempts),
	)
}
