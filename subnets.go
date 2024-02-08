package edgecloud

import (
	"context"
	"encoding/json"
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
	Update(context.Context, string, *SubnetworkUpdateRequest) (*Subnetwork, *Response, error)

	SubnetworksMetadata
}

type SubnetworksMetadata interface {
	MetadataList(context.Context, string) ([]MetadataDetailed, *Response, error)
	MetadataCreate(context.Context, string, *Metadata) (*Response, error)
	MetadataUpdate(context.Context, string, *Metadata) (*Response, error)
	MetadataDeleteItem(context.Context, string, *MetadataItemOptions) (*Response, error)
	MetadataGetItem(context.Context, string, *MetadataItemOptions) (*MetadataDetailed, *Response, error)
}

// SubnetworksServiceOp handles communication with Subnetworks methods of the EdgecenterCloud API.
type SubnetworksServiceOp struct {
	client *Client
}

var _ SubnetworksService = &SubnetworksServiceOp{}

// Subnetwork represents an EdgecenterCloud Subnetwork.
type Subnetwork struct {
	ID                     string             `json:"id"`
	Name                   string             `json:"name"`
	NetworkID              string             `json:"network_id"`
	IPVersion              int                `json:"ip_version"`
	EnableDHCP             bool               `json:"enable_dhcp"`
	ConnectToNetworkRouter bool               `json:"connect_to_network_router"`
	CIDR                   string             `json:"cidr"` // TODO add cidr parsing.
	CreatedAt              string             `json:"created_at"`
	UpdatedAt              string             `json:"updated_at"`
	CreatorTaskID          string             `json:"creator_task_id"`
	TaskID                 string             `json:"task_id"`
	AvailableIps           int                `json:"available_ips"`
	TotalIps               int                `json:"total_ips"`
	HasRouter              bool               `json:"has_router"`
	DNSNameservers         []net.IP           `json:"dns_nameservers"`
	HostRoutes             []HostRoute        `json:"host_routes"`
	GatewayIP              net.IP             `json:"gateway_ip"`
	Metadata               []MetadataDetailed `json:"metadata,omitempty"`
	Region                 string             `json:"region"`
	ProjectID              int                `json:"project_id"`
	RegionID               int                `json:"region_id"`
}

// SubnetworkCreateRequest represents a request to create a Subnetwork.
type SubnetworkCreateRequest struct {
	Name                   string      `json:"name" required:"true"`
	NetworkID              string      `json:"network_id" required:"true"`
	EnableDHCP             bool        `json:"enable_dhcp,omitempty"`
	CIDR                   string      `json:"cidr" required:"true"`
	ConnectToNetworkRouter bool        `json:"connect_to_network_router"`
	DNSNameservers         []net.IP    `json:"dns_nameservers,omitempty"`
	GatewayIP              *net.IP     `json:"gateway_ip,omitempty"`
	Metadata               Metadata    `json:"metadata"`
	HostRoutes             []HostRoute `json:"host_routes,omitempty"`
}

// SubnetworkUpdateRequest represents a request to update a Subnetwork properties.
type SubnetworkUpdateRequest struct {
	Name           string      `json:"name" required:"true"`
	DNSNameservers []net.IP    `json:"dns_nameservers"`
	EnableDHCP     bool        `json:"enable_dhcp"`
	HostRoutes     []HostRoute `json:"host_routes"`
	GatewayIP      *net.IP     `json:"gateway_ip"`
}

type CIDR struct {
	net.IPNet
}

// UnmarshalJSON - implements Unmarshaler interface for CIDR.
func (c *CIDR) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	cd, err := ParseCIDRString(s)
	if err != nil {
		return err
	}
	*c = *cd

	return nil
}

// MarshalJSON - implements Marshaler interface for CIDR.
func (c CIDR) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// String - implements Stringer.
func (c CIDR) String() string {
	return c.IPNet.String()
}

func ParseCIDRString(s string) (*CIDR, error) {
	_, nt, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	return &CIDR{IPNet: *nt}, nil
}

// HostRoute represents a route that should be used by devices with IPs from
// a subnet (not including local subnet route).
type HostRoute struct {
	Destination CIDR   `json:"destination"`
	NextHop     net.IP `json:"nexthop"`
}

// SubnetworkListOptions specifies the optional query parameters to List method.
type SubnetworkListOptions struct {
	NetworkID  string `url:"network_id,omitempty"  validate:"omitempty"`
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

	path := s.client.addProjectRegionPath(subnetsBasePathV1)
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

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(subnetsBasePathV1), subnetworkID)

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
func (s *SubnetworksServiceOp) Create(ctx context.Context, reqBody *SubnetworkCreateRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(subnetsBasePathV1)

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

// Delete the Subnetwork.
func (s *SubnetworksServiceOp) Delete(ctx context.Context, subnetworkID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(subnetworkID, "subnetworkID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(subnetsBasePathV1), subnetworkID)

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

// Update the Subnetwork properties.
func (s *SubnetworksServiceOp) Update(ctx context.Context, subnetworkID string, reqBody *SubnetworkUpdateRequest) (*Subnetwork, *Response, error) {
	if resp, err := isValidUUID(subnetworkID, "subnetworkID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(subnetsBasePathV1), subnetworkID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
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

// MetadataList subnetwork detailed metadata items.
func (s *SubnetworksServiceOp) MetadataList(ctx context.Context, subnetworkID string) ([]MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(subnetworkID, "subnetworkID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataList(ctx, s.client, subnetworkID, subnetsBasePathV1)
}

// MetadataCreate or update subnetwork metadata.
func (s *SubnetworksServiceOp) MetadataCreate(ctx context.Context, subnetworkID string, reqBody *Metadata) (*Response, error) {
	if resp, err := isValidUUID(subnetworkID, "subnetworkID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataCreate(ctx, s.client, subnetworkID, subnetsBasePathV1, reqBody)
}

// MetadataUpdate subnetwork metadata.
func (s *SubnetworksServiceOp) MetadataUpdate(ctx context.Context, subnetworkID string, reqBody *Metadata) (*Response, error) {
	if resp, err := isValidUUID(subnetworkID, "subnetworkID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataUpdate(ctx, s.client, subnetworkID, subnetsBasePathV1, reqBody)
}

// MetadataDeleteItem a subnetwork metadata item by key.
func (s *SubnetworksServiceOp) MetadataDeleteItem(ctx context.Context, subnetworkID string, opts *MetadataItemOptions) (*Response, error) {
	if resp, err := isValidUUID(subnetworkID, "subnetworkID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataDeleteItem(ctx, s.client, subnetworkID, subnetsBasePathV1, opts)
}

// MetadataGetItem subnetwork detailed metadata.
func (s *SubnetworksServiceOp) MetadataGetItem(ctx context.Context, subnetworkID string, opts *MetadataItemOptions) (*MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(subnetworkID, "subnetworkID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataGetItem(ctx, s.client, subnetworkID, subnetsBasePathV1, opts)
}
