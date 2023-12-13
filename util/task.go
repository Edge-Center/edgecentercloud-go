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
	defaultTimeout         = time.Minute
)

var (
	errTaskWaitTimeout    = errors.New("a timeout occurred")
	errTaskWithErrorState = errors.New("task with error state")
	errTaskStateUnknown   = errors.New("unknown task state")
)

type TaskResult struct {
	DdosProfiles     []int    `json:"ddos_profiles"`
	FloatingIPs      []string `json:"floatingips"`
	HealthMonitors   []string `json:"healthmonitors"`
	Images           []string `json:"images"`
	Instances        []string `json:"instances"`
	L7Polices        []string `json:"l7polices"`
	L7Rules          []string `json:"l7rules"`
	Listeners        []string `json:"listeners"`
	Loadbalancers    []string `json:"loadbalancers"`
	Members          []string `json:"members"`
	Networks         []string `json:"networks"`
	Pools            []string `json:"pools"`
	Ports            []string `json:"ports"`
	Projects         []string `json:"projects"`
	ReservedFixedIPs []string `json:"reserved_fixed_ips"`
	Routers          []string `json:"routers"`
	Secrets          []string `json:"secrets"`
	ServerGroups     []string `json:"servergroups"`
	Snapshots        []string `json:"snapshots"`
	Subnets          []string `json:"subnets"`
	Volumes          []string `json:"volumes"`
}

type TaskAPIFunc[T any] func(ctx context.Context, opt T) (*edgecloud.TaskResponse, *edgecloud.Response, error)

func ExecuteAndExtractTaskResult[T any](ctx context.Context, apiFunc TaskAPIFunc[T], opt T, client *edgecloud.Client, timeouts ...time.Duration) (*TaskResult, error) {
	task, _, err := apiFunc(ctx, opt)
	if err != nil {
		return nil, err
	}

	taskInfo, err := WaitAndGetTaskInfo(ctx, client, task.Tasks[0], timeouts...)
	if err != nil {
		return nil, err
	}

	return ExtractTaskResultFromTask(taskInfo)
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
			if failCount <= taskFailure {
				failCount++
				time.Sleep(taskGetInfoRetrySecond * time.Second)
				continue
			}

			return nil, err
		}

		switch taskInfo.State {
		case edgecloud.TaskStateRunning, edgecloud.TaskStateNew:
			<-time.After(taskGetInfoRetrySecond * time.Second)
		case edgecloud.TaskStateError:
			return nil, edgecloud.NewArgError("taskID", errTaskWithErrorState.Error())
		case edgecloud.TaskStateFinished:
			return taskInfo, nil
		default:
			return nil, fmt.Errorf("%w: [%s]", errTaskStateUnknown, taskInfo.State)
		}
	}
}

func WaitForTaskComplete(ctx context.Context, client *edgecloud.Client, taskID string, timeouts ...time.Duration) error {
	timeout := defaultTimeout

	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	done := make(chan error)
	go func() {
		_, err := waitTask(ctx, client, taskID)
		done <- err
		close(done)
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return edgecloud.NewArgError("taskID", errTaskWaitTimeout.Error())
	}
}

func WaitAndGetTaskInfo(ctx context.Context, client *edgecloud.Client, taskID string, timeouts ...time.Duration) (*edgecloud.Task, error) {
	timeout := defaultTimeout

	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	done := make(chan struct{})
	var taskInfo *edgecloud.Task
	var taskErr error

	go func() {
		taskInfo, taskErr = waitTask(ctx, client, taskID)
		close(done)
	}()

	select {
	case <-done:
		return taskInfo, taskErr
	case <-time.After(timeout):
		return nil, edgecloud.NewArgError("taskID", errTaskWaitTimeout.Error())
	}
}
