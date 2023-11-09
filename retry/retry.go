package retry

import (
	"context"
	"time"
)

type Action int

const (
	Succeed Action = iota
	Fail
	Retry
)

type Worker func(ctx context.Context) error

type Policy func(err error) Action

type Retryer struct {
	backoff     *Backoff
	retryPolicy Policy
}

func NewRetryer(backoff *Backoff, retryPolicy Policy) Retryer {
	if retryPolicy == nil {
		retryPolicy = DefaultRetryPolicy
	}

	return Retryer{
		backoff:     backoff,
		retryPolicy: retryPolicy,
	}
}

func (r *Retryer) Run(ctx context.Context, work Worker) error {
	defer r.backoff.Reset()
	for {
		err := work(ctx)

		switch r.retryPolicy(err) {
		case Succeed, Fail:
			return err
		case Retry:
			var delay time.Duration
			if delay = r.backoff.Next(); delay == Stop {
				return err
			}
			timeout := time.After(delay)
			if err := r.sleep(ctx, timeout); err != nil {
				return err
			}
		}
	}
}

func (r *Retryer) sleep(ctx context.Context, t <-chan time.Time) error {
	select {
	case <-t:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func DefaultRetryPolicy(err error) Action {
	if err == nil {
		return Succeed
	}
	return Retry
}
