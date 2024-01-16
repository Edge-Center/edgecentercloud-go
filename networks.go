package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	networksBasePathV1          = "/v1/networks"
	availablenetworksBasePathV1 = "/v1/availablenetworks"
)

const (
	networksPorts = "ports"
)

// NetworksService is an interface for creating and managing Networks with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/networks
type NetworksService interface {
	List(context.Context, *NetworkListOptions) ([]Network, *Response, error)
	Get(context.Context, string) (*Network, *Response, error)
	Create(context.Context, *NetworkCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
	UpdateName(context.Context, string, *Name) (*Network, *Response, error)
	ListNetworksWithSubnets(context.Context, *NetworksWithSubnetsOptions) ([]NetworkSubnetwork, *Response, error)
	PortList(context.Context, string) ([]PortsInstance, *Response, error)

	NetworksMetadata
}

type NetworksMetadata interface {
	MetadataList(context.Context, string) ([]MetadataDetailed, *Response, error)
	MetadataCreate(context.Context, string, *Metadata) (*Response, error)
	MetadataUpdate(context.Context, string, *Metadata) (*Response, error)
	MetadataDeleteItem(context.Context, string, *MetadataItemOptions) (*Response, error)
	MetadataGetItem(context.Context, string, *MetadataItemOptions) (*MetadataDetailed, *Response, error)
}

// NetworksServiceOp handles communication with Networks methods of the EdgecenterCloud API.
type NetworksServiceOp struct {
	client *Client
}

var _ NetworksService = &NetworksServiceOp{}

// Network represents an EdgecenterCloud Network.
type Network struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	CreatedAt      string             `json:"created_at"`
	CreatorTaskID  string             `json:"creator_task_id"`
	Default        bool               `json:"default"`
	External       bool               `json:"external"`
	MTU            int                `json:"mtu"`
	Metadata       []MetadataDetailed `json:"metadata,omitempty"`
	ProjectID      int                `json:"project_id"`
	Region         string             `json:"region"`
	RegionID       int                `json:"region_id"`
	SegmentationID int                `json:"segmentation_id"`
	Shared         bool               `json:"shared"`
	Subnets        []string           `json:"subnets"`
	TaskID         string             `json:"task_id"`
	Type           string             `json:"type"`
	UpdatedAt      string             `json:"updated_at"`
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

// NetworkListOptions specifies the optional query parameters to List method.
type NetworkListOptions struct {
	OrderBy    string `url:"order_by,omitempty"  validate:"omitempty"`
	MetadataKV string `url:"metadata_kv,omitempty"  validate:"omitempty"`
	MetadataK  string `url:"metadata_k,omitempty"  validate:"omitempty"`
}

// networksRoot represents a Networks root.
type networksRoot struct {
	Count    int
	Networks []Network `json:"results"`
}

// NetworkSubnetwork represents an EdgecenterCloud Network with info about Subnets.
type NetworkSubnetwork struct {
	Metadata       []MetadataDetailed `json:"metadata,omitempty"`
	UpdatedAt      string             `json:"updated_at"`
	Name           string             `json:"name"`
	CreatedAt      string             `json:"created_at"`
	Type           string             `json:"type"`
	External       bool               `json:"external"`
	TaskID         string             `json:"task_id"`
	Default        bool               `json:"default"`
	RegionID       int                `json:"region_id"`
	Shared         bool               `json:"shared"`
	Region         string             `json:"region"`
	MTU            int                `json:"mtu"`
	SegmentationID int                `json:"segmentation_id"`
	CreatorTaskID  string             `json:"creator_task_id"`
	ID             string             `json:"id"`
	ProjectID      int                `json:"project_id"`
	Subnets        []Subnetwork       `json:"subnets"`
}

// PortsInstance represent instances ports.
type PortsInstance struct {
	InstanceID   string `json:"instance_id"`
	InstanceName string `json:"instance_name"`
	ID           string `json:"id"`
}

// portsInstanceRoot represents a Port Instance root.
type portsInstanceRoot struct {
	Count         int
	PortsInstance []PortsInstance `json:"results"`
}

// NetworksWithSubnetsOptions specifies the optional query parameters to ListNetworksWithSubnets method.
type NetworksWithSubnetsOptions struct {
	NetworkID   string `url:"network_id,omitempty"  validate:"omitempty"`
	NetworkType string `url:"network_type,omitempty"  validate:"omitempty"`
	OrderBy     string `url:"order_by,omitempty"  validate:"omitempty"`
	Shared      bool   `url:"shared,omitempty"  validate:"omitempty"`
	MetadataKV  string `url:"metadata_kv,omitempty"  validate:"omitempty"`
	MetadataK   string `url:"metadata_k,omitempty"  validate:"omitempty"`
}

// networksSubnetworksRoot represents a NetworkSubnetwork root.
type networksSubnetworksRoot struct {
	Count             int
	NetworkSubnetwork []NetworkSubnetwork `json:"results"`
}

// List get networks.
func (s *NetworksServiceOp) List(ctx context.Context, opts *NetworkListOptions) ([]Network, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(networksBasePathV1)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(networksRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Networks, resp, err
}

// Get individual Network.
func (s *NetworksServiceOp) Get(ctx context.Context, networkID string) (*Network, *Response, error) {
	if resp, err := isValidUUID(networkID, "networkID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(networksBasePathV1), networkID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	network := new(Network)
	resp, err := s.client.Do(ctx, req, network)
	if err != nil {
		return nil, resp, err
	}

	return network, resp, err
}

// Create a Network.
func (s *NetworksServiceOp) Create(ctx context.Context, reqBody *NetworkCreateRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(networksBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
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

// Delete the Network.
func (s *NetworksServiceOp) Delete(ctx context.Context, networkID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(networkID, "networkID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(networksBasePathV1), networkID)

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

// UpdateName of the network.
func (s *NetworksServiceOp) UpdateName(ctx context.Context, networkID string, reqBody *Name) (*Network, *Response, error) {
	if resp, err := isValidUUID(networkID, "networkID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(networksBasePathV1), networkID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	network := new(Network)
	resp, err := s.client.Do(ctx, req, network)
	if err != nil {
		return nil, resp, err
	}

	return network, resp, err
}

// ListNetworksWithSubnets get networks with details of subnets.
func (s *NetworksServiceOp) ListNetworksWithSubnets(ctx context.Context, opts *NetworksWithSubnetsOptions) ([]NetworkSubnetwork, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(availablenetworksBasePathV1)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(networksSubnetworksRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.NetworkSubnetwork, resp, err
}

// PortList get instance ports by network_id.
func (s *NetworksServiceOp) PortList(ctx context.Context, networkID string) ([]PortsInstance, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(networksBasePathV1), networkID)
	path = fmt.Sprintf("%s/%s", path, networksPorts)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(portsInstanceRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.PortsInstance, resp, err
}

// MetadataList network detailed metadata items.
func (s *NetworksServiceOp) MetadataList(ctx context.Context, networkID string) ([]MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(networkID, "networkID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataList(ctx, s.client, networkID, networksBasePathV1)
}

// MetadataCreate or update network metadata.
func (s *NetworksServiceOp) MetadataCreate(ctx context.Context, networkID string, reqBody *Metadata) (*Response, error) {
	if resp, err := isValidUUID(networkID, "networkID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataCreate(ctx, s.client, networkID, networksBasePathV1, reqBody)
}

// MetadataUpdate network metadata.
func (s *NetworksServiceOp) MetadataUpdate(ctx context.Context, networkID string, reqBody *Metadata) (*Response, error) {
	if resp, err := isValidUUID(networkID, "networkID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataUpdate(ctx, s.client, networkID, networksBasePathV1, reqBody)
}

// MetadataDeleteItem a network metadata item by key.
func (s *NetworksServiceOp) MetadataDeleteItem(ctx context.Context, networkID string, opts *MetadataItemOptions) (*Response, error) {
	if resp, err := isValidUUID(networkID, "networkID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataDeleteItem(ctx, s.client, networkID, networksBasePathV1, opts)
}

// MetadataGetItem network detailed metadata.
func (s *NetworksServiceOp) MetadataGetItem(ctx context.Context, networkID string, opts *MetadataItemOptions) (*MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(networkID, "networkID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataGetItem(ctx, s.client, networkID, networksBasePathV1, opts)
}
