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
