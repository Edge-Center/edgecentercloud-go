package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	networksBasePathV1 = "/v1/networks"
)

// NetworksService is an interface for creating and managing Networks with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/networks
type NetworksService interface {
	Get(context.Context, string) (*Network, *Response, error)
	Create(context.Context, *NetworkCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
}

// NetworksServiceOp handles communication with Networks methods of the EdgecenterCloud API.
type NetworksServiceOp struct {
	client *Client
}

var _ NetworksService = &NetworksServiceOp{}

// Network represents an EdgecenterCloud Network.
type Network struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	CreatedAt      string   `json:"created_at"`
	CreatorTaskID  string   `json:"creator_task_id"`
	Default        bool     `json:"default"`
	External       bool     `json:"external"`
	MTU            int      `json:"mtu"`
	Metadata       Metadata `json:"metadata"`
	ProjectID      int      `json:"project_id"`
	Region         string   `json:"region"`
	RegionID       int      `json:"region_id"`
	SegmentationID int      `json:"segmentation_id"`
	Shared         bool     `json:"shared"`
	Subnets        []string `json:"subnets"`
	TaskID         string   `json:"task_id"`
	Type           string   `json:"type"`
	UpdatedAt      string   `json:"updated_at"`
}

// NetworkCreateRequest represents a request to create a Network.
type NetworkCreateRequest struct {
	Name         string      `json:"name" required:"true" validate:"required"`
	CreateRouter bool        `json:"create_router"`
	Type         NetworkType `json:"type,omitempty" validate:"omitempty"`
	Metadata     Metadata    `json:"metadata,omitempty" validate:"omitempty,dive"`
}

type NetworkType string

const (
	VLAN  NetworkType = "vlan"
	VXLAN NetworkType = "vxlan"
)

// networkRoot represents a Network root.
type networkRoot struct {
	Network *Network      `json:"network"`
	Tasks   *TaskResponse `json:"tasks"`
}

// Get individual Network.
func (s *NetworksServiceOp) Get(ctx context.Context, networkID string) (*Network, *Response, error) {
	if err := isValidUUID(networkID, "networkID"); err != nil {
		return nil, nil, err
	}

	if err := s.client.Validate(); err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addServicePath(networksBasePathV1), networkID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(networkRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Network, resp, err
}

// Create a Network.
func (s *NetworksServiceOp) Create(ctx context.Context, createRequest *NetworkCreateRequest) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if err := s.client.Validate(); err != nil {
		return nil, nil, err
	}

	path := s.client.addServicePath(networksBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(networkRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}

// Delete the Network.
func (s *NetworksServiceOp) Delete(ctx context.Context, networkID string) (*TaskResponse, *Response, error) {
	if err := isValidUUID(networkID, "networkID"); err != nil {
		return nil, nil, err
	}

	if err := s.client.Validate(); err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addServicePath(networksBasePathV1), networkID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(networkRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}
