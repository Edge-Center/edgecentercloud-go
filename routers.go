package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	routersBasePathV1 = "/v1/routers"
)

const (
	routersAttach = "attach"
	routersDetach = "detach"
)

// RoutersService is an interface for creating and managing Routers with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/routers
type RoutersService interface {
	List(context.Context) ([]Router, *Response, error)
	Create(context.Context, *RouterCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
	Get(context.Context, string) (*Router, *Response, error)
	Update(context.Context, string, *RouterUpdateRequest) (*Router, *Response, error)
	Attach(context.Context, string, *RouterAttachRequest) (*Router, *Response, error)
	Detach(context.Context, string, *RouterDetachRequest) (*Router, *Response, error)
}

// RoutersServiceOp handles communication with Routers methods of the EdgecenterCloud API.
type RoutersServiceOp struct {
	client *Client
}

var _ RoutersService = &RoutersServiceOp{}

// Router represents an EdgecenterCloud Router.
type Router struct {
	Region              string              `json:"region"`
	UpdatedAt           string              `json:"updated_at"`
	CreatedAt           string              `json:"created_at"`
	Name                string              `json:"name"`
	ID                  string              `json:"id"`
	RegionID            int                 `json:"region_id"`
	ProjectID           int                 `json:"project_id"`
	TaskID              string              `json:"task_id"`
	Status              string              `json:"status"`
	CreatorTaskID       string              `json:"creator_task_id"`
	ExternalGatewayInfo ExternalGatewayInfo `json:"external_gateway_info"`
	Interfaces          []RouterInterface   `json:"interfaces"`
	Routes              []HostRoute         `json:"routes"`
}

// RouterInterface represents a router instance interface.
type RouterInterface struct {
	PortID        string   `json:"port_id"`
	IPAssignments []PortIP `json:"ip_assignments"`
	MacAddress    string   `json:"mac_address"`
	NetworkID     string   `json:"network_id"`
}

type ExternalGatewayInfo struct {
	ExternalFixedIPs []ExternalFixedIP `json:"external_fixed_ips"`
	NetworkID        string            `json:"network_id"`
	EnableSnat       bool              `json:"enable_snat"`
}

type ExternalFixedIP struct {
	IPAddress string `json:"ip_address"`
	SubnetID  string `json:"subnet_id"`
}

type RouterCreateRequest struct {
	Interfaces          []RouterInterfaceCreate   `json:"interfaces"`
	ExternalGatewayInfo ExternalGatewayInfoCreate `json:"external_gateway_info"`
	Name                string                    `json:"name" required:"true" validate:"required"`
	Routes              []HostRoute               `json:"routes"`
}

type RouterUpdateRequest struct {
	ExternalGatewayInfo ExternalGatewayInfoCreate `json:"external_gateway_info"`
	Name                string                    `json:"name" required:"true" validate:"required"`
	Routes              []HostRoute               `json:"routes"`
}

type RouterInterfaceCreate struct {
	SubnetID string `json:"subnet_id" required:"true" validate:"required"`
	Type     string `json:"type"`
}

type ExternalGatewayInfoCreate struct {
	EnableSnat bool   `json:"enable_snat"`
	Type       string `json:"type"`
	NetworkID  string `json:"network_id"`
}

type RouterAttachRequest struct {
	SubnetID string `json:"subnet_id" required:"true" validate:"required"`
}

type RouterDetachRequest struct {
	SubnetID string `json:"subnet_id" required:"true" validate:"required"`
}

// routersRoot represents a Routers root.
type routersRoot struct {
	Count   int
	Routers []Router `json:"results"`
}

// List get routers.
func (s *RoutersServiceOp) List(ctx context.Context) ([]Router, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(routersBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(routersRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Routers, resp, err
}

// Create a Router.
func (s *RoutersServiceOp) Create(ctx context.Context, reqBody *RouterCreateRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(routersBasePathV1)

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

// Delete a Router.
func (s *RoutersServiceOp) Delete(ctx context.Context, routerID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(routerID, "routerID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(routersBasePathV1), routerID)

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

// Get a Router.
func (s *RoutersServiceOp) Get(ctx context.Context, routerID string) (*Router, *Response, error) {
	if resp, err := isValidUUID(routerID, "routerID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(routersBasePathV1), routerID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	router := new(Router)
	resp, err := s.client.Do(ctx, req, router)
	if err != nil {
		return nil, resp, err
	}

	return router, resp, err
}

// Update a Router.
func (s *RoutersServiceOp) Update(ctx context.Context, routerID string, reqBody *RouterUpdateRequest) (*Router, *Response, error) {
	if resp, err := isValidUUID(routerID, "routerID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(routersBasePathV1), routerID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	router := new(Router)
	resp, err := s.client.Do(ctx, req, router)
	if err != nil {
		return nil, resp, err
	}

	return router, resp, err
}

// Attach a subnet to a Router.
func (s *RoutersServiceOp) Attach(ctx context.Context, routerID string, reqBody *RouterAttachRequest) (*Router, *Response, error) {
	if resp, err := isValidUUID(routerID, "routerID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(routersBasePathV1), routerID, routersAttach)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	router := new(Router)
	resp, err := s.client.Do(ctx, req, router)
	if err != nil {
		return nil, resp, err
	}

	return router, resp, err
}

// Detach a subnet from a Router.
func (s *RoutersServiceOp) Detach(ctx context.Context, routerID string, reqBody *RouterDetachRequest) (*Router, *Response, error) {
	if resp, err := isValidUUID(routerID, "routerID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(routersBasePathV1), routerID, routersDetach)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	router := new(Router)
	resp, err := s.client.Do(ctx, req, router)
	if err != nil {
		return nil, resp, err
	}

	return router, resp, err
}
