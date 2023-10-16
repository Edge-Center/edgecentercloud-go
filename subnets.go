package edgecloud

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/google/uuid"
)

const (
	subnetsBasePathV1 = "/v1/subnets"
)

// SubnetworksService is an interface for creating and managing Subnetworks with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/subnets
type SubnetworksService interface {
	Get(context.Context, string, *ServicePath) (*Subnetwork, *Response, error)
	Create(context.Context, *SubnetworkCreateRequest, *ServicePath) (*TaskResponse, *Response, error)
	Delete(context.Context, string, *ServicePath) (*TaskResponse, *Response, error)
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
	CIDR                   net.IPNet   `json:"cidr"`
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
	Metadata               Metadata    `json:"metadata"`
	Region                 string      `json:"region"`
	ProjectID              int         `json:"project_id"`
	RegionID               int         `json:"region_id"`
}

// SubnetworkCreateRequest represents a request to create a Subnetwork.
type SubnetworkCreateRequest struct {
	Name                   string      `json:"name" required:"true"`
	NetworkID              string      `json:"network_id" required:"true"`
	EnableDHCP             bool        `json:"enable_dhcp,omitempty"`
	CIDR                   net.IPNet   `json:"cidr" required:"true"`
	ConnectToNetworkRouter bool        `json:"connect_to_network_router"`
	DNSNameservers         []net.IP    `json:"dns_nameservers,omitempty"`
	GatewayIP              *net.IP     `json:"gateway_ip"`
	Metadata               Metadata    `json:"metadata"`
	HostRoutes             []HostRoute `json:"host_routes,omitempty"`
}

// HostRoute represents a route that should be used by devices with IPs from
// a subnet (not including local subnet route).
type HostRoute struct {
	Destination net.IPNet `json:"destination"`
	NextHop     net.IP    `json:"nexthop"`
}

// subnetworkRoot represents a Subnetwork root.
type subnetworkRoot struct {
	Subnetwork *Subnetwork   `json:"subnetwork"`
	Tasks      *TaskResponse `json:"tasks"`
}

// Get individual Network.
func (s *SubnetworksServiceOp) Get(ctx context.Context, subnetworkID string, p *ServicePath) (*Subnetwork, *Response, error) {
	if _, err := uuid.Parse(subnetworkID); err != nil {
		return nil, nil, NewArgError("subnetworkID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(subnetsBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, subnetworkID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(subnetworkRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Subnetwork, resp, err
}

// Create a Network.
func (s *SubnetworksServiceOp) Create(ctx context.Context, createRequest *SubnetworkCreateRequest, p *ServicePath) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(subnetsBasePathV1, p)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(subnetworkRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}

// Delete the Network.
func (s *SubnetworksServiceOp) Delete(ctx context.Context, subnetworkID string, p *ServicePath) (*TaskResponse, *Response, error) {
	if _, err := uuid.Parse(subnetworkID); err != nil {
		return nil, nil, NewArgError("subnetworkID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(subnetsBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, subnetworkID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(subnetworkRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}
