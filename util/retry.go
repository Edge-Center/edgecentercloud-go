package util

import "github.com/avast/retry-go/v4"

const attempts = 10

func WithRetry(f retry.RetryableFunc) error {

	return retry.Do(
		f,
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(attempts),
	)
}
