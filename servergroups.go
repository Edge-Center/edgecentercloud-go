package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	servergroupsBasePathV1 = "/v1/servergroups"
)

// ServerGroupsService is an interface for creating and managing Server Groups with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/servergroups
type ServerGroupsService interface {
	List(context.Context) ([]ServerGroup, *Response, error)
	Get(context.Context, string) (*ServerGroup, *Response, error)
	Create(context.Context, *ServerGroupCreateRequest) (*ServerGroup, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

// ServerGroupsServiceOp handles communication with Server Groups methods of the EdgecenterCloud API.
type ServerGroupsServiceOp struct {
	client *Client
}

var _ ServerGroupsService = &ServerGroupsServiceOp{}

// ServerGroup represents an EdgecenterCloud Server Group.
type ServerGroup struct {
	ID        string                `json:"servergroup_id"`
	Policy    ServerGroupPolicy     `json:"policy"`
	Name      string                `json:"name"`
	Instances []ServerGroupInstance `json:"instances"`
	ProjectID int                   `json:"project_id"`
	RegionID  int                   `json:"region_id"`
	Region    string                `json:"region"`
}

// ServerGroupInstance represent an instances in server group.
type ServerGroupInstance struct {
	InstanceID   string `json:"instance_id"`
	InstanceName string `json:"instance_name"`
}

type ServerGroupPolicy string

const (
	ServerGroupPolicyAffinity     ServerGroupPolicy = "affinity"
	ServerGroupPolicyAntiAffinity ServerGroupPolicy = "anti-affinity"
)

// ServerGroupCreateRequest represents a request to create a Server Group.
type ServerGroupCreateRequest struct {
	Name   string            `json:"name" required:"true"`
	Policy ServerGroupPolicy `json:"policy" required:"true" validate:"enum"`
}

// serverGroupsRoot represents a Server Group root.
type serverGroupsRoot struct {
	Count       int
	ServerGroup []ServerGroup `json:"results"`
}

// List get Server Groups.
func (s *ServerGroupsServiceOp) List(ctx context.Context) ([]ServerGroup, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(servergroupsBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(serverGroupsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.ServerGroup, resp, err
}

// Get individual Server Group.
func (s *ServerGroupsServiceOp) Get(ctx context.Context, serverGroupID string) (*ServerGroup, *Response, error) {
	if resp, err := isValidUUID(serverGroupID, "serverGroupID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(servergroupsBasePathV1), serverGroupID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	serverGroup := new(ServerGroup)
	resp, err := s.client.Do(ctx, req, serverGroup)
	if err != nil {
		return nil, resp, err
	}

	return serverGroup, resp, err
}

// Create a Server Group.
func (s *ServerGroupsServiceOp) Create(ctx context.Context, reqBody *ServerGroupCreateRequest) (*ServerGroup, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(servergroupsBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	serverGroup := new(ServerGroup)
	resp, err := s.client.Do(ctx, req, serverGroup)
	if err != nil {
		return nil, resp, err
	}

	return serverGroup, resp, err
}

// Delete the Server Group.
func (s *ServerGroupsServiceOp) Delete(ctx context.Context, serverGroupID string) (*Response, error) {
	if resp, err := isValidUUID(serverGroupID, "serverGroupID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(servergroupsBasePathV1), serverGroupID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
