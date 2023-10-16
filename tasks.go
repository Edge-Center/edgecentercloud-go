package edgecloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	tasksBasePathV1 = "/v1/tasks"
)

// TasksService is an interface for managing Tasks with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/tasks
type TasksService interface {
	Get(context.Context, string) (*Task, *Response, error)
}

// TasksServiceOp handles communication with Tasks methods of the EdgecenterCloud API.
type TasksServiceOp struct {
	client *Client
}

var _ TasksService = &TasksServiceOp{}

// Task represents an EdgecenterCloud Task.
type Task struct {
	ID                string                  `json:"id"`
	TaskType          string                  `json:"task_type"`
	ProjectID         int                     `json:"project_id"`
	RegionID          int                     `json:"region_id"`
	ClientID          int                     `json:"client_id"`
	UserID            int                     `json:"user_id"`
	UserClientID      int                     `json:"user_client_id"`
	State             TaskState               `json:"state"`
	Data              *map[string]interface{} `json:"data"`
	CreatedResources  map[string]interface{}  `json:"created_resources"`
	RequestID         string                  `json:"request_id"`
	Error             *string                 `json:"error"`
	CreatedOn         string                  `json:"created_on"`
	UpdatedOn         *string                 `json:"updated_on,omitempty"`
	FinishedOn        *string                 `json:"finished_on,omitempty"`
	AcknowledgedAt    string                  `json:"acknowledged_at,omitempty"`
	AcknowledgedBy    int                     `json:"acknowledged_by,omitempty"`
	JobID             string                  `json:"job_id"`
	ScheduleID        string                  `json:"schedule_id"`
	LifecyclePolicyID int                     `json:"lifecycle_policy_id"`
}

type TaskState string

const (
	TaskStateNew      TaskState = "NEW"
	TaskStateRunning  TaskState = "RUNNING"
	TaskStateFinished TaskState = "FINISHED"
	TaskStateError    TaskState = "ERROR"
)

// TaskResponse is an EdgecenterCloud response with list of created tasks.
type TaskResponse struct {
	Tasks []string `json:"tasks"`
}

// taskRoot represents a Task root.
type taskRoot struct {
	Task *Task `json:"task"`
}

// Get individual Task.
func (s *TasksServiceOp) Get(ctx context.Context, taskID string) (*Task, *Response, error) {
	if _, err := uuid.Parse(taskID); err != nil {
		return nil, nil, NewArgError("taskID", "should be the correct UUID")
	}

	path := fmt.Sprintf("%s/%s", tasksBasePathV1, taskID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(taskRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Task, resp, err
}
