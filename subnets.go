package edgecloud

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

const (
	subnetsBasePathV1 = "/v1/subnets"
)

// SubnetworksService is an interface for creating and managing Subnetworks with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/subnets
type SubnetworksService interface {
	List(context.Context, *SubnetworkListOptions) ([]Subnetwork, *Response, error)
	Get(context.Context, string) (*Subnetwork, *Response, error)
	Create(context.Context, *SubnetworkCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
}

// SubnetworksServiceOp handles communication with Subnetworks methods of the EdgecenterCloud API.
type SubnetworksServiceOp struct {
	client *Client
}

var _ SubnetworksService = &SubnetworksServiceOp{}

// Subnetwork represents an EdgecenterCloud Subnetwork.
type Subnetwork struct {
	ID                     string      `json:"id"`
	Name                   string      `json:"name"`
	NetworkID              string      `json:"network_id"`
	IPVersion              int         `json:"ip_version"`
	EnableDHCP             bool        `json:"enable_dhcp"`
	ConnectToNetworkRouter bool        `json:"connect_to_network_router"`
	CIDR                   string      `json:"cidr"` // TODO add cidr parsing.
	CreatedAt              string      `json:"created_at"`
	UpdatedAt              string      `json:"updated_at"`
	CreatorTaskID          string      `json:"creator_task_id"`
	TaskID                 string      `json:"task_id"`
	AvailableIps           int         `json:"available_ips"`
	TotalIps               int         `json:"total_ips"`
	HasRouter              bool        `json:"has_router"`
	DNSNameservers         []net.IP    `json:"dns_nameservers"`
	HostRoutes             []HostRoute `json:"host_routes"`
	GatewayIP              net.IP      `json:"gateway_ip"`
	Metadata               []Metadata  `json:"metadata,omitempty"`
	Region                 string      `json:"region"`
	ProjectID              int         `json:"project_id"`
	RegionID               int         `json:"region_id"`
}

// SubnetworkCreateRequest represents a request to create a Subnetwork.
type SubnetworkCreateRequest struct {
	Name                   string      `json:"name" required:"true"`
	NetworkID              string      `json:"network_id" required:"true"`
	EnableDHCP             bool        `json:"enable_dhcp,omitempty"`
	CIDR                   string      `json:"cidr" required:"true"`
	ConnectToNetworkRouter bool        `json:"connect_to_network_router"`
	DNSNameservers         []net.IP    `json:"dns_nameservers,omitempty"`
	GatewayIP              *net.IP     `json:"gateway_ip"`
	Metadata               Metadata    `json:"metadata"`
	HostRoutes             []HostRoute `json:"host_routes,omitempty"`
}

type CIDR struct {
	net.IPNet
}

// HostRoute represents a route that should be used by devices with IPs from
// a subnet (not including local subnet route).
type HostRoute struct {
	Destination net.IPNet `json:"destination"`
	NextHop     net.IP    `json:"nexthop"`
}

// SubnetworkListOptions specifies the optional query parameters to List method.
type SubnetworkListOptions struct {
	NetworkID  bool   `url:"network_id,omitempty"  validate:"omitempty"`
	MetadataKV string `url:"metadata_kv,omitempty"  validate:"omitempty"`
	MetadataK  string `url:"metadata_k,omitempty"  validate:"omitempty"`
}

// subnetworkRoot represents a Subnetworks root.
type subnetworkRoot struct {
	Count       int
	Subnetworks []Subnetwork `json:"results"`
}

// List get subnetworks.
func (s *SubnetworksServiceOp) List(ctx context.Context, opts *SubnetworkListOptions) ([]Subnetwork, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addServicePath(subnetsBasePathV1)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(subnetworkRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Subnetworks, resp, err
}

// Get individual Subnetwork.
func (s *SubnetworksServiceOp) Get(ctx context.Context, subnetworkID string) (*Subnetwork, *Response, error) {
	if resp, err := isValidUUID(subnetworkID, "subnetworkID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addServicePath(subnetsBasePathV1), subnetworkID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	subnetwork := new(Subnetwork)
	resp, err := s.client.Do(ctx, req, subnetwork)
	if err != nil {
		return nil, resp, err
	}

	return subnetwork, resp, err
}

// Create a Subnetwork.
func (s *SubnetworksServiceOp) Create(ctx context.Context, createRequest *SubnetworkCreateRequest) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addServicePath(subnetsBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
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

// Delete the Subnetwork.
func (s *SubnetworksServiceOp) Delete(ctx context.Context, subnetworkID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(subnetworkID, "subnetworkID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addServicePath(subnetsBasePathV1), subnetworkID)

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
