package util

import (
	"context"
	"errors"
	"fmt"
	"time"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

const (
	// taskFailure is the amount of times we can fail before deciding
	// the check for task result is a total failure.
	taskFailure            = 3
	taskGetInfoRetrySecond = 5
)

var (
	errTaskWaitTimeout    = errors.New("a timeout occurred")
	errTaskWithErrorState = errors.New("task with error state")
	errTaskStateUnknown   = errors.New("unknown task state")
)

// WaitForTaskComplete waits for the task to complete.
func WaitForTaskComplete(ctx context.Context, client *edgecloud.Client, taskID string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel() // О

	completed := false
	failCount := 0
	for !completed {
		taskInfo, _, err := client.Tasks.Get(ctx, taskID)
		if err != nil {
			select {
			case <-ctx.Done():
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					return errTaskWaitTimeout
				}
				return ctx.Err()
			default:
			}
			if failCount <= taskFailure {
				failCount++
				continue
			}

			return err
		}

		switch taskInfo.State {
		case edgecloud.TaskStateRunning, edgecloud.TaskStateNew:
			select {
			case <-time.After(taskGetInfoRetrySecond * time.Second):
			case <-ctx.Done():
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					return errTaskWaitTimeout
				}
				return ctx.Err()
			}
		case edgecloud.TaskStateError:
			return errTaskWithErrorState
		case edgecloud.TaskStateFinished:
			completed = true
		default:
			return fmt.Errorf("%w: [%s]", errTaskStateUnknown, taskInfo.State)
		}
	}

	return nil
}
