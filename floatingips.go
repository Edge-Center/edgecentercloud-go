package edgecloud

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

const (
	floatingipsBasePathV1      = "/v1/floatingips"
	availableFloatingipsPathV1 = "/v1/availablefloatingips"
)

const (
	floatingipsAssign   = "assign"
	floatingipsUnAssign = "unassign"
)

// FloatingIPsService is an interface for creating and managing FloatingIPs with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/floatingips
type FloatingIPsService interface {
	List(context.Context) ([]FloatingIP, *Response, error)
	Get(context.Context, string) (*FloatingIP, *Response, error)
	Create(context.Context, *FloatingIPCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
	Assign(context.Context, string, *AssignFloatingIPRequest) (*FloatingIP, *Response, error)
	UnAssign(context.Context, string) (*FloatingIP, *Response, error)
	ListAvailable(context.Context) ([]FloatingIP, *Response, error)

	FloatingIPMetadata
}

type FloatingIPMetadata interface {
	MetadataList(context.Context, string) ([]MetadataDetailed, *Response, error)
	MetadataCreate(context.Context, string, *MetadataCreateRequest) (*Response, error)
	MetadataUpdate(context.Context, string, *MetadataCreateRequest) (*Response, error)
	MetadataDeleteItem(context.Context, string, *MetadataItemOptions) (*Response, error)
	MetadataGetItem(context.Context, string, *MetadataItemOptions) (*MetadataDetailed, *Response, error)
}

// FloatingipsServiceOp handles communication with FloatingIPs methods of the EdgecenterCloud API.
type FloatingipsServiceOp struct {
	client *Client
}

var _ FloatingIPsService = &FloatingipsServiceOp{}

// FloatingIP represents an EdgecenterCloud FloatingIP.
type FloatingIP struct {
	ID                string       `json:"id"`
	CreatedAt         string       `json:"created_at"`
	UpdatedAt         string       `json:"updated_at"`
	Status            string       `json:"status"`
	FixedIPAddress    net.IP       `json:"fixed_ip_address,omitempty"`
	FloatingIPAddress string       `json:"floating_ip_address,omitempty"`
	DNSDomain         string       `json:"dns_domain"`
	DNSName           string       `json:"dns_name"`
	RouterID          string       `json:"router_id"`
	SubnetID          string       `json:"subnet_id"`
	CreatorTaskID     string       `json:"creator_task_id"`
	Metadata          []Metadata   `json:"metadata,omitempty"`
	TaskID            string       `json:"task_id"`
	PortID            string       `json:"port_id,omitempty"`
	ProjectID         int          `json:"project_id"`
	RegionID          int          `json:"region_id"`
	Region            string       `json:"region"`
	Instance          Instance     `json:"instance,omitempty"`
	Loadbalancer      Loadbalancer `json:"loadbalancer,omitempty"`
}

type FloatingIPSource string

const (
	NewFloatingIP      FloatingIPSource = "new"
	ExistingFloatingIP FloatingIPSource = "existing"
)

type InterfaceFloatingIP struct {
	Source             FloatingIPSource `json:"source" validate:"required,enum"`
	ExistingFloatingID string           `json:"existing_floating_id" validate:"rfe=Source:existing,sfe=Source:new,omitempty,UUID"`
}

// FloatingIPCreateRequest represents a request to create a Floating IP.
type FloatingIPCreateRequest struct {
	PortID         string   `json:"port_id,omitempty"`
	FixedIPAddress net.IP   `json:"fixed_ip_address,omitempty"`
	Metadata       Metadata `json:"metadata,omitempty"`
}

// AssignFloatingIPRequest represents a request to assign a Floating IP to an instance or a load balancer.
type AssignFloatingIPRequest struct {
	PortID         string   `json:"port_id" validate:"required"`
	FixedIPAddress net.IP   `json:"fixed_ip_address,omitempty"`
	Metadata       Metadata `json:"metadata,omitempty"`
}

// floatingipsRoot represents a FloatingIPs root.
type floatingipsRoot struct {
	Count       int
	FloatingIPs []FloatingIP `json:"results"`
}

// List get floating IPs.
func (s *FloatingipsServiceOp) List(ctx context.Context) ([]FloatingIP, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(floatingipsBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(floatingipsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.FloatingIPs, resp, err
}

// Get a Floating IP.
func (s *FloatingipsServiceOp) Get(ctx context.Context, fipID string) (*FloatingIP, *Response, error) {
	if resp, err := isValidUUID(fipID, "fipID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(floatingipsBasePathV1), fipID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	fip := new(FloatingIP)
	resp, err := s.client.Do(ctx, req, fip)
	if err != nil {
		return nil, resp, err
	}

	return fip, resp, err
}

// Create a Floating IP.
func (s *FloatingipsServiceOp) Create(ctx context.Context, reqBody *FloatingIPCreateRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(floatingipsBasePathV1)

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

// Delete the Floating IP.
func (s *FloatingipsServiceOp) Delete(ctx context.Context, fipID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(fipID, "fipID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(floatingipsBasePathV1), fipID)

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

// Assign a floating IP to an instance or a load balancer.
func (s *FloatingipsServiceOp) Assign(ctx context.Context, fipID string, reqBody *AssignFloatingIPRequest) (*FloatingIP, *Response, error) {
	if resp, err := isValidUUID(fipID, "fipID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(floatingipsBasePathV1), fipID, floatingipsAssign)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	fip := new(FloatingIP)
	resp, err := s.client.Do(ctx, req, fip)
	if err != nil {
		return nil, resp, err
	}

	return fip, resp, err
}

// UnAssign a floating IP from an instance or a load balancer.
func (s *FloatingipsServiceOp) UnAssign(ctx context.Context, fipID string) (*FloatingIP, *Response, error) {
	if resp, err := isValidUUID(fipID, "fipID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(floatingipsBasePathV1), fipID, floatingipsUnAssign)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	fip := new(FloatingIP)
	resp, err := s.client.Do(ctx, req, fip)
	if err != nil {
		return nil, resp, err
	}

	return fip, resp, err
}

// ListAvailable floating IPs.
func (s *FloatingipsServiceOp) ListAvailable(ctx context.Context) ([]FloatingIP, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(availableFloatingipsPathV1)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(floatingipsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.FloatingIPs, resp, err
}

// MetadataList floating IP detailed metadata items.
func (s *FloatingipsServiceOp) MetadataList(ctx context.Context, fipID string) ([]MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(fipID, "fipID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataList(ctx, s.client, fipID, floatingipsBasePathV1)
}

// MetadataCreate or update floating IP metadata.
func (s *FloatingipsServiceOp) MetadataCreate(ctx context.Context, fipID string, reqBody *MetadataCreateRequest) (*Response, error) {
	if resp, err := isValidUUID(fipID, "fipID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataCreate(ctx, s.client, fipID, floatingipsBasePathV1, reqBody)
}

// MetadataUpdate floating IP metadata.
func (s *FloatingipsServiceOp) MetadataUpdate(ctx context.Context, fipID string, reqBody *MetadataCreateRequest) (*Response, error) {
	if resp, err := isValidUUID(fipID, "fipID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataUpdate(ctx, s.client, fipID, floatingipsBasePathV1, reqBody)
}

// MetadataDeleteItem a floating IP metadata item by key.
func (s *FloatingipsServiceOp) MetadataDeleteItem(ctx context.Context, fipID string, opts *MetadataItemOptions) (*Response, error) {
	if resp, err := isValidUUID(fipID, "fipID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataDeleteItem(ctx, s.client, fipID, floatingipsBasePathV1, opts)
}

// MetadataGetItem floating IP detailed metadata.
func (s *FloatingipsServiceOp) MetadataGetItem(ctx context.Context, fipID string, opts *MetadataItemOptions) (*MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(fipID, "fipID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataGetItem(ctx, s.client, fipID, floatingipsBasePathV1, opts)
}
