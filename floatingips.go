package edgecloud

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

const (
	floatingipsBasePathV1 = "/v1/floatingips"
)

// FloatingIPsService is an interface for creating and managing FloatingIPs with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/floatingips
type FloatingIPsService interface {
	Get(context.Context, string) (*FloatingIP, *Response, error)
	Create(context.Context, *FloatingIPCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
}

// FloatingipsServiceOp handles communication with FloatingIPs methods of the EdgecenterCloud API.
type FloatingipsServiceOp struct {
	client *Client
}

var _ FloatingIPsService = &FloatingipsServiceOp{}

// FloatingIP represents an EdgecenterCloud FloatingIP.
type FloatingIP struct {
	ID                string   `json:"id"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
	Status            string   `json:"status"`
	FixedIPAddress    net.IP   `json:"fixed_ip_address,omitempty"`
	FloatingIPAddress string   `json:"floating_ip_address,omitempty"`
	DNSDomain         string   `json:"dns_domain"`
	DNSName           string   `json:"dns_name"`
	RouterID          string   `json:"router_id"`
	SubnetID          string   `json:"subnet_id"`
	CreatorTaskID     string   `json:"creator_task_id"`
	Metadata          Metadata `json:"metadata,omitempty"`
	TaskID            string   `json:"task_id"`
	PortID            string   `json:"port_id,omitempty"`
	ProjectID         int      `json:"project_id"`
	RegionID          int      `json:"region_id"`
	Region            string   `json:"region"`
	Instance          Instance `json:"instance,omitempty"`
}

type FloatingIPSource string

const (
	NewFloatingIP      FloatingIPSource = "new"
	ExistingFloatingIP FloatingIPSource = "existing"
)

type InterfaceFloatingIP struct {
	Source             FloatingIPSource `json:"source" validate:"required,enum"`
	ExistingFloatingID string           `json:"existing_floating_id" validate:"rfe=Source:existing,sfe=Source:new,omitempty,uuid"`
}

// FloatingIPCreateRequest represents a request to create a FloatingIP.
type FloatingIPCreateRequest struct {
	PortID         string   `json:"port_id,omitempty"`
	FixedIPAddress net.IP   `json:"fixed_ip_address,omitempty"`
	Metadata       Metadata `json:"metadata,omitempty"`
}

// Get individual FloatingIP.
func (s *FloatingipsServiceOp) Get(ctx context.Context, fipID string) (*FloatingIP, *Response, error) {
	if err := isValidUUID(fipID, "fipID"); err != nil {
		return nil, nil, err
	}

	if err := s.client.Validate(); err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addServicePath(floatingipsBasePathV1), fipID)

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

// Create a FloatingIP.
func (s *FloatingipsServiceOp) Create(ctx context.Context, createRequest *FloatingIPCreateRequest) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if err := s.client.Validate(); err != nil {
		return nil, nil, err
	}

	path := s.client.addServicePath(floatingipsBasePathV1)

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

// Delete the FloatingIP.
func (s *FloatingipsServiceOp) Delete(ctx context.Context, fipID string) (*TaskResponse, *Response, error) {
	if err := isValidUUID(fipID, "fipID"); err != nil {
		return nil, nil, err
	}

	if err := s.client.Validate(); err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addServicePath(floatingipsBasePathV1), fipID)

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
