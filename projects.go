package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	projectsBasePath = "/v1/projects"
)

// ProjectsService is an interface for creating and managing Projects with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/projects
type ProjectsService interface {
	Get(context.Context, string) (*Project, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
	Update(context.Context, string, *ProjectUpdateRequest) (*Project, *Response, error)
	List(context.Context, *ProjectListOptions) ([]Project, *Response, error)
	Create(context.Context, *ProjectCreateRequest) (*Project, *Response, error)
}

// ProjectsServiceOp handles communication with Projects methods of the EdgecenterCloud API.
type ProjectsServiceOp struct {
	client *Client
}

var _ ProjectsService = &ProjectsServiceOp{}

// ProjectState the model 'ProjectState'.
type ProjectState string

// List of ProjectState.
const (
	ProjectStateActive   ProjectState = "ACTIVE"
	ProjectStateDeleted  ProjectState = "DELETED"
	ProjectStateDeleting ProjectState = "DELETING"
)

// Project represents a EdgecenterCloud Project configuration.
type Project struct {
	ID          int          `json:"id"`
	ClientID    int          `json:"client_id"`
	CreatedAt   string       `json:"created_at"`
	Description string       `json:"description"`
	IsDefault   bool         `json:"is_default"`
	Name        string       `json:"name"`
	State       ProjectState `json:"state"`
	TaskID      *string      `json:"task_id"`
}

// ProjectUpdateRequest represents a request to update a Project.
type ProjectUpdateRequest struct {
	Description string `json:"description" required:"true"`
	Name        string `json:"name" required:"true"`
}

// ProjectCreateRequest represents a request to create a Project.
type ProjectCreateRequest struct {
	Name        string `json:"name" required:"true"`
	ClientID    string `json:"client_id"`
	State       string `json:"state"`
	Description string `json:"description"`
}

type ProjectListOptions struct {
	ClientID       string `url:"key,omitempty" validate:"omitempty"`
	OrderBy        string `url:"order_by,omitempty" validate:"omitempty"`
	Name           string `url:"name,omitempty" validate:"omitempty"`
	IncludeDeleted bool   `url:"include_deleted,omitempty" validate:"omitempty"`
}

// projectsRoot represents Projects root.
type projectsRoot struct {
	Count    int
	Projects []Project `json:"results"`
}

// Get retrieves a single project by its ID.
func (s *ProjectsServiceOp) Get(ctx context.Context, projectID string) (*Project, *Response, error) {
	path := fmt.Sprintf("%s/%s", projectsBasePath, projectID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}

// Delete a project.
func (s *ProjectsServiceOp) Delete(ctx context.Context, projectID string) (*TaskResponse, *Response, error) {
	path := fmt.Sprintf("%s/%s", projectsBasePath, projectID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(TaskResponse)
	resp, err := s.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}

// Update a project.
func (s *ProjectsServiceOp) Update(ctx context.Context, projectID string, reqBody *ProjectUpdateRequest) (*Project, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s", projectsBasePath, projectID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}

// List gets projects.
func (s *ProjectsServiceOp) List(ctx context.Context, opts *ProjectListOptions) ([]Project, *Response, error) {
	path, err := addOptions(projectsBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	projectRoot := new(projectsRoot)
	resp, err := s.client.Do(ctx, req, projectRoot)
	if err != nil {
		return nil, resp, err
	}

	return projectRoot.Projects, resp, err
}

// Create a project.
func (s *ProjectsServiceOp) Create(ctx context.Context, reqBody *ProjectCreateRequest) (*Project, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, projectsBasePath, reqBody)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}
