package edgecloud

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

const (
	reservedFixedIPsBasePathV1 = "/v1/reserved_fixed_ips"
)

const (
	reservedFixedIPsConnectedDevices = "connected_devices"
	reservedFixedIPsAvailableDevices = "available_devices"
)

// ReservedFixedIPsService is an interface for creating and managing Reserved Fixed IPs with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/reserved_fixed_ips
type ReservedFixedIPsService interface {
	List(context.Context, *ReservedFixedIPListOptions) ([]ReservedFixedIP, *Response, error)
	Create(context.Context, *ReservedFixedIPCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
	Get(context.Context, string) (*ReservedFixedIP, *Response, error)
	SwitchVIPStatus(context.Context, string, *SwitchVIPStatusRequest) (*ReservedFixedIP, *Response, error)
	ListInstancePorts(context.Context, string) ([]ReservedFixedIPInstancePort, *Response, error)
	AddInstancePorts(context.Context, string, *AddInstancePortsRequest) ([]ReservedFixedIPInstancePort, *Response, error)
	ReplaceInstancePorts(context.Context, string, *AddInstancePortsRequest) ([]ReservedFixedIPInstancePort, *Response, error)
	ListInstancePortsAvailable(context.Context, string) ([]ReservedFixedIPInstancePort, *Response, error)
}

// ReservedFixedIPsServiceOp handles communication with Reserved Fixed IPs methods of the EdgecenterCloud API.
type ReservedFixedIPsServiceOp struct {
	client *Client
}

var _ ReservedFixedIPsService = &ReservedFixedIPsServiceOp{}

// ReservedFixedIP represents an EdgecenterCloud ReservedFixedIP.
type ReservedFixedIP struct {
	Region              string                `json:"region"`
	CreatedAt           string                `json:"created_at"`
	UpdatedAt           string                `json:"updated_at"`
	Name                string                `json:"name"`
	RegionID            int                   `json:"region_id"`
	PortID              string                `json:"port_id,omitempty"`
	FixedIPAddress      net.IP                `json:"fixed_ip_address,omitempty"`
	TaskID              string                `json:"task_id"`
	IsVIP               bool                  `json:"is_vip"`
	IsExternal          bool                  `json:"is_external"`
	ProjectID           int                   `json:"project_id"`
	NetworkID           string                `json:"network_id"`
	CreatorTaskID       string                `json:"creator_task_id"`
	Status              string                `json:"status"`
	SubnetID            string                `json:"subnet_id"`
	AllowedAddressPairs []AllowedAddressPairs `json:"allowed_address_pairs"`
	Network             Network               `json:"network"`
	Reservation         Reservation           `json:"reservation"`
}

type Reservation struct {
	ResourceType string `json:"resource_type"`
	Status       string `json:"status"`
	ResourceID   string `json:"resource_id"`
}

// ReservedFixedIPListOptions specifies the optional query parameters to get ReservedFixedIP List method.
type ReservedFixedIPListOptions struct {
	ExternalOnly   bool   `url:"external_only,omitempty"  validate:"omitempty"`
	InternalOnly   bool   `url:"internal_only,omitempty"  validate:"omitempty"`
	AvailableOnly  bool   `url:"available_only,omitempty"  validate:"omitempty"`
	VIPOnly        bool   `url:"vip_only,omitempty"  validate:"omitempty"`
	SearchPrefixIP string `url:"search_prefix_ip,omitempty"  validate:"omitempty"`
	DeviceID       string `url:"device_id,omitempty"  validate:"omitempty"`
	Limit          int    `url:"limit,omitempty"  validate:"omitempty"`
	Offset         int    `url:"offset,omitempty"  validate:"omitempty"`
}

type ReservedFixedIPCreateRequest struct {
	IsVIP     bool                `json:"is_vip"`
	Type      ReservedFixedIPType `json:"type" required:"true" validate:"required,enum"`
	NetworkID string              `json:"network_id,omitempty" validate:"rfe=Type:ip_address;any_subnet,omitempty,uuid4"`
	SubnetID  string              `json:"subnet_id,omitempty" validate:"rfe=Type:subnet,omitempty,uuid4"`
	IPAddress string              `json:"ip_address,omitempty" validate:"rfe=Type:ip_address,omitempty"`
}

type ReservedFixedIPType string

const (
	ReservedFixedIPTypeExternal  = "external"
	ReservedFixedIPTypeSubnet    = "subnet"
	ReservedFixedIPTypeAnySubnet = "any_subnet"
	ReservedFixedIPTypeIPAddress = "ip_address"
)

type SwitchVIPStatusRequest struct {
	IsVIP bool `json:"is_vip"`
}

type ReservedFixedIPInstancePort struct {
	PortID        string   `json:"port_id,omitempty"`
	IPAssignments []PortIP `json:"ip_assignments"`
	InstanceID    string   `json:"instance_id,omitempty"`
	InstanceName  string   `json:"instance_name,omitempty"`
	Network       Network  `json:"network"`
}

type AddInstancePortsRequest struct {
	PortIDs []string `json:"port_ids"`
}

// reservedFixedIPRoot represents a ReservedFixedIPs root.
type reservedFixedIPRoot struct {
	Count            int
	ReservedFixedIPs []ReservedFixedIP `json:"results"`
}

// instancePortRoot represents a ReservedFixedIPInstancePort root.
type instancePortRoot struct {
	Count         int
	InstancePorts []ReservedFixedIPInstancePort `json:"results"`
}

// List get Reserved Fixed IPs.
func (s *ReservedFixedIPsServiceOp) List(ctx context.Context, opts *ReservedFixedIPListOptions) ([]ReservedFixedIP, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(reservedFixedIPsBasePathV1)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(reservedFixedIPRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.ReservedFixedIPs, resp, err
}

// Create a Reserved Fixed IP.
func (s *ReservedFixedIPsServiceOp) Create(ctx context.Context, reqBody *ReservedFixedIPCreateRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(reservedFixedIPsBasePathV1)

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

// Delete a Reserved Fixed IP.
func (s *ReservedFixedIPsServiceOp) Delete(ctx context.Context, reservedFixedIPID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(reservedFixedIPID, "reservedFixedIPID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(reservedFixedIPsBasePathV1), reservedFixedIPID)

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

// Get a Reserved Fixed IP.
func (s *ReservedFixedIPsServiceOp) Get(ctx context.Context, reservedFixedIPID string) (*ReservedFixedIP, *Response, error) {
	if resp, err := isValidUUID(reservedFixedIPID, "reservedFixedIPID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(reservedFixedIPsBasePathV1), reservedFixedIPID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	reservedFixedIP := new(ReservedFixedIP)
	resp, err := s.client.Do(ctx, req, reservedFixedIP)
	if err != nil {
		return nil, resp, err
	}

	return reservedFixedIP, resp, err
}

// SwitchVIPStatus of a Reserved Fixed IP.
func (s *ReservedFixedIPsServiceOp) SwitchVIPStatus(ctx context.Context, reservedFixedIPID string, reqBody *SwitchVIPStatusRequest) (*ReservedFixedIP, *Response, error) {
	if resp, err := isValidUUID(reservedFixedIPID, "reservedFixedIPID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(reservedFixedIPsBasePathV1), reservedFixedIPID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	reservedFixedIP := new(ReservedFixedIP)
	resp, err := s.client.Do(ctx, req, reservedFixedIP)
	if err != nil {
		return nil, resp, err
	}

	return reservedFixedIP, resp, err
}

// ListInstancePorts that share a VIP.
func (s *ReservedFixedIPsServiceOp) ListInstancePorts(ctx context.Context, reservedFixedIPID string) ([]ReservedFixedIPInstancePort, *Response, error) {
	if resp, err := isValidUUID(reservedFixedIPID, "reservedFixedIPID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(reservedFixedIPsBasePathV1), reservedFixedIPID, reservedFixedIPsConnectedDevices)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(instancePortRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.InstancePorts, resp, err
}

// AddInstancePorts that share a VIP.
func (s *ReservedFixedIPsServiceOp) AddInstancePorts(ctx context.Context, reservedFixedIPID string, reqBody *AddInstancePortsRequest) ([]ReservedFixedIPInstancePort, *Response, error) {
	if resp, err := isValidUUID(reservedFixedIPID, "reservedFixedIPID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(reservedFixedIPsBasePathV1), reservedFixedIPID, reservedFixedIPsConnectedDevices)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	root := new(instancePortRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.InstancePorts, resp, err
}

// ReplaceInstancePorts that share a VIP.
func (s *ReservedFixedIPsServiceOp) ReplaceInstancePorts(ctx context.Context, reservedFixedIPID string, reqBody *AddInstancePortsRequest) ([]ReservedFixedIPInstancePort, *Response, error) {
	if resp, err := isValidUUID(reservedFixedIPID, "reservedFixedIPID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(reservedFixedIPsBasePathV1), reservedFixedIPID, reservedFixedIPsConnectedDevices)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	root := new(instancePortRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.InstancePorts, resp, err
}

// ListInstancePortsAvailable for connecting to a VIP.
func (s *ReservedFixedIPsServiceOp) ListInstancePortsAvailable(ctx context.Context, reservedFixedIPID string) ([]ReservedFixedIPInstancePort, *Response, error) {
	if resp, err := isValidUUID(reservedFixedIPID, "reservedFixedIPID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(reservedFixedIPsBasePathV1), reservedFixedIPID, reservedFixedIPsAvailableDevices)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(instancePortRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.InstancePorts, resp, err
}
