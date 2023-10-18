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

// serverGroupRoot represents a Server Group root.
type serverGroupRoot struct {
	ServerGroup *ServerGroup `json:"server_group"`
}

// Get individual Server Group.
func (s *ServerGroupsServiceOp) Get(ctx context.Context, serverGroupID string) (*ServerGroup, *Response, error) {
	if err := isValidUUID(serverGroupID, "serverGroupID"); err != nil {
		return nil, nil, err
	}

	if err := s.client.Validate(); err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addServicePath(servergroupsBasePathV1), serverGroupID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(serverGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.ServerGroup, resp, err
}

// Create a Server Group.
func (s *ServerGroupsServiceOp) Create(ctx context.Context, createRequest *ServerGroupCreateRequest) (*ServerGroup, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if err := s.client.Validate(); err != nil {
		return nil, nil, err
	}

	path := s.client.addServicePath(servergroupsBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(serverGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.ServerGroup, resp, err
}

// Delete the Server Group.
func (s *ServerGroupsServiceOp) Delete(ctx context.Context, serverGroupID string) (*Response, error) {
	if err := isValidUUID(serverGroupID, "serverGroupID"); err != nil {
		return nil, err
	}

	if err := s.client.Validate(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addServicePath(servergroupsBasePathV1), serverGroupID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
