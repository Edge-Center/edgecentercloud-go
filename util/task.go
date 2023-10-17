package util

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"

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

type TaskResult struct {
	DdosProfiles   []int    `json:"ddos_profiles"`
	FloatingIPs    []string `json:"floatingips"`
	HealthMonitors []string `json:"healthmonitors"`
	Images         []string `json:"images"`
	Instances      []string `json:"instances"`
	Listeners      []string `json:"listeners"`
	LoadBalancers  []string `json:"loadbalancers"`
	Members        []string `json:"members"`
	Networks       []string `json:"networks"`
	Pools          []string `json:"pools"`
	Ports          []string `json:"ports"`
	Routers        []string `json:"routers"`
	Secrets        []string `json:"secrets"`
	Snapshots      []string `json:"snapshots"`
	Subnets        []string `json:"subnets"`
	Volumes        []string `json:"volumes"`
}

func ExtractTaskResultFromTask(task *edgecloud.Task) (*TaskResult, error) {
	var result TaskResult
	if err := mapstructure.Decode(task.CreatedResources, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

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
