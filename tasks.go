package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	tasksBasePathV1 = "/v1/tasks"
)

const (
	tasksActive      = "active"
	tasksAcknowledge = "acknowledge"
)

// TasksService is an interface for managing Tasks with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/tasks
type TasksService interface {
	ListActive(context.Context) ([]Task, *Response, error)
	Acknowledge(context.Context, string) (*Task, *Response, error)
	AcknowledgeAll(context.Context, *TaskAcknowledgeAllOptions) (*Response, error)
	Get(context.Context, string) (*Task, *Response, error)
	List(context.Context, *TaskListOptions) ([]Task, *Response, error)
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

type TaskSorting string

const (
	TaskSortingAsc  TaskSorting = "asc"
	TaskSortingDesc TaskSorting = "desc"
)

// TaskResponse is an EdgecenterCloud response with list of created tasks.
type TaskResponse struct {
	Tasks []string `json:"tasks"`
}

type TaskAcknowledgeAllOptions struct {
	ProjectID int    `url:"project_id,omitempty" validate:"omitempty"`
	RegionID  string `url:"region_id,omitempty" validate:"omitempty"`
}

type TaskListOptions struct {
	ClientID          int         `url:"client_id,omitempty" validate:"omitempty"`
	ProjectID         int         `url:"project_id,omitempty" validate:"omitempty"`
	RegionID          int         `url:"region_id,omitempty" validate:"omitempty"`
	ScheduleID        string      `url:"schedule_id,omitempty" validate:"omitempty"`
	LifecyclePolicyID int         `url:"lifecycle_policy_id,omitempty" validate:"omitempty"`
	IsAcknowledged    bool        `url:"is_acknowledged,omitempty" validate:"omitempty"`
	FromTimestamp     string      `url:"from_timestamp,omitempty" validate:"omitempty"`
	ToTimestamp       string      `url:"to_timestamp,omitempty" validate:"omitempty"`
	Limit             int         `url:"limit,omitempty" validate:"omitempty"`
	Offset            int         `url:"offset,omitempty" validate:"omitempty"`
	TaskType          string      `url:"task_type,omitempty" validate:"omitempty"`
	State             TaskState   `url:"state,omitempty" validate:"omitempty"`
	Sorting           TaskSorting `url:"sorting,omitempty" validate:"omitempty"`
}

// tasksRoot represents a Tasks root.
type tasksRoot struct {
	Count int
	Tasks []Task `json:"results"`
}

// ListActive get active tasks.
func (s *TasksServiceOp) ListActive(ctx context.Context) ([]Task, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(tasksBasePathV1), tasksActive)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(tasksRoot)
	resp, err := s.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks.Tasks, resp, err
}

// Acknowledge one task on project scope.
func (s *TasksServiceOp) Acknowledge(ctx context.Context, taskID string) (*Task, *Response, error) {
	if resp, err := isValidUUID(taskID, "taskID"); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", tasksBasePathV1, taskID, tasksAcknowledge)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	task := new(Task)
	resp, err := s.client.Do(ctx, req, task)
	if err != nil {
		return nil, resp, err
	}

	return task, resp, err
}

// AcknowledgeAll client tasks in project or region.
func (s *TasksServiceOp) AcknowledgeAll(ctx context.Context, opts *TaskAcknowledgeAllOptions) (*Response, error) {
	path := fmt.Sprintf("%s/%s", tasksBasePathV1, tasksAcknowledge)

	path, err := addOptions(path, opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Get individual Task.
func (s *TasksServiceOp) Get(ctx context.Context, taskID string) (*Task, *Response, error) {
	if resp, err := isValidUUID(taskID, "taskID"); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", tasksBasePathV1, taskID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	task := new(Task)
	resp, err := s.client.Do(ctx, req, task)
	if err != nil {
		return nil, resp, err
	}

	return task, resp, err
}

// List gets tasks.
func (s *TasksServiceOp) List(ctx context.Context, opts *TaskListOptions) ([]Task, *Response, error) {
	path, err := addOptions(tasksBasePathV1, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(tasksRoot)
	resp, err := s.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks.Tasks, resp, err
}
