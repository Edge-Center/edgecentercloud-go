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

// waitTask waits for the task to complete.
func waitTask(ctx context.Context, client *edgecloud.Client, taskID string) (*edgecloud.Task, error) {
	failCount := 0
	for {
		taskInfo, _, err := client.Tasks.Get(ctx, taskID)
		if err != nil {
			select {
			case <-ctx.Done():
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					return nil, errTaskWaitTimeout
				}
				return nil, ctx.Err()
			default:
			}
			if failCount <= taskFailure {
				failCount++
				continue
			}

			return nil, err
		}

		switch taskInfo.State {
		case edgecloud.TaskStateRunning, edgecloud.TaskStateNew:
			select {
			case <-time.After(taskGetInfoRetrySecond * time.Second):
			case <-ctx.Done():
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					return nil, errTaskWaitTimeout
				}
				return nil, ctx.Err()
			}
		case edgecloud.TaskStateError:
			return nil, errTaskWithErrorState
		case edgecloud.TaskStateFinished:
			return taskInfo, nil
		default:
			return nil, fmt.Errorf("%w: [%s]", errTaskStateUnknown, taskInfo.State)
		}
	}
}

func WaitForTaskComplete(ctx context.Context, client *edgecloud.Client, taskID string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := waitTask(ctx, client, taskID)
	return err
}

func WaitAndGetTaskInfo(ctx context.Context, client *edgecloud.Client, taskID string, timeout time.Duration) (*edgecloud.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return waitTask(ctx, client, taskID)
}
