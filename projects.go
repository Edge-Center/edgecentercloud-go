package edgecloud

import (
	"context"
	"net/http"
	"path"
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

var _ ProjectsService = &ProjectsServiceOp{}

// Get retrieves a single project by its ID.
func (p *ProjectsServiceOp) Get(ctx context.Context, projectID string) (*Project, *Response, error) {
	return p.getHelper(ctx, projectID)
}

func (p *ProjectsServiceOp) getHelper(ctx context.Context, projectID string) (*Project, *Response, error) {
	path := path.Join(projectsBasePath, projectID)

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := p.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}
