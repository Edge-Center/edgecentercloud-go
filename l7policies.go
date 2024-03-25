package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	l7policiesBasePathV1 = "/v1/l7policies"
)

// L7PoliciesService is an interface for creating and managing L7Policies with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/l7policies
type L7PoliciesService interface {
	List(context.Context) ([]L7Policy, *Response, error)
	Create(context.Context, *L7PolicyCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
	Get(context.Context, string) (*L7Policy, *Response, error)
	Update(context.Context, string, *L7PolicyUpdateRequest) (*TaskResponse, *Response, error)
}

// L7PoliciesServiceOp handles communication with L7Policies methods of the EdgecenterCloud API.
type L7PoliciesServiceOp struct {
	client *Client
}

var _ L7PoliciesService = &L7PoliciesServiceOp{}

// L7Policy represents an EdgecenterCloud L7Policy.
type L7Policy struct {
	RegionID           int            `json:"region_id"`
	ProjectID          int            `json:"project_id"`
	Name               string         `json:"name"`
	Region             string         `json:"region"`
	ID                 string         `json:"id"`
	TaskID             string         `json:"task_id"`
	RedirectHTTPCode   *int           `json:"redirect_http_code"`
	Tags               []string       `json:"tags"`
	ListenerID         string         `json:"listener_id"`
	RedirectPoolID     *string        `json:"redirect_pool_id"`
	OperatingStatus    string         `json:"operating_status"`
	ProvisioningStatus string         `json:"provisioning_status"`
	RedirectURL        *string        `json:"redirect_url"`
	Position           int            `json:"position"`
	RedirectPrefix     *string        `json:"redirect_prefix"`
	Action             L7PolicyAction `json:"action"`
	Rules              []L7Rule       `json:"rules"`
	CreatedAt          string         `json:"created_at"`
	UpdatedAt          string         `json:"updated_at,omitempty"`
}

type L7PolicyAction string

const (
	L7PolicyActionRedirectPrefix L7PolicyAction = "REDIRECT_PREFIX"
	L7PolicyActionRedirectToPool L7PolicyAction = "REDIRECT_TO_POOL"
	L7PolicyActionRedirectToURL  L7PolicyAction = "REDIRECT_TO_URL"
	L7PolicyActionReject         L7PolicyAction = "REJECT"
)

type L7PolicyCreateRequest struct {
	Tags             []string       `json:"tags,omitempty"`
	RedirectHTTPCode int            `json:"redirect_http_code,omitempty"`
	ListenerID       string         `json:"listener_id" required:"true" validate:"required"`
	Position         int            `json:"position,omitempty"`
	Name             string         `json:"name,omitempty"`
	Action           L7PolicyAction `json:"action" required:"true" validate:"required"`
	RedirectURL      string         `json:"redirect_url,omitempty"`
	RedirectPrefix   string         `json:"redirect_prefix,omitempty"`
	RedirectPoolID   string         `json:"redirect_pool_id,omitempty"`
}

type L7PolicyUpdateRequest struct {
	Tags             []string       `json:"tags,omitempty"`
	RedirectHTTPCode int            `json:"redirect_http_code,omitempty"`
	Position         int            `json:"position,omitempty"`
	Name             string         `json:"name,omitempty"`
	Action           L7PolicyAction `json:"action" required:"true" validate:"required"`
	RedirectURL      string         `json:"redirect_url,omitempty"`
	RedirectPrefix   string         `json:"redirect_prefix,omitempty"`
	RedirectPoolID   string         `json:"redirect_pool_id,omitempty"`
}

// l7policiesRoot represents a L7Policy root.
type l7policiesRoot struct {
	Count      int
	L7policies []L7Policy `json:"results"`
}

// List get L7policies.
func (s *L7PoliciesServiceOp) List(ctx context.Context) ([]L7Policy, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(l7policiesBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(l7policiesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.L7policies, resp, err
}

// Create a L7Policy.
func (s *L7PoliciesServiceOp) Create(ctx context.Context, reqBody *L7PolicyCreateRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(l7policiesBasePathV1)

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

// Delete a L7Policy.
func (s *L7PoliciesServiceOp) Delete(ctx context.Context, l7PolicyID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(l7PolicyID, "l7PolicyID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(l7policiesBasePathV1), l7PolicyID)

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

// Get a L7Policy.
func (s *L7PoliciesServiceOp) Get(ctx context.Context, l7PolicyID string) (*L7Policy, *Response, error) {
	if resp, err := isValidUUID(l7PolicyID, "l7PolicyID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(l7policiesBasePathV1), l7PolicyID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	l7Policy := new(L7Policy)
	resp, err := s.client.Do(ctx, req, l7Policy)
	if err != nil {
		return nil, resp, err
	}

	return l7Policy, resp, err
}

// Update replace L7Policy properties.
func (s *L7PoliciesServiceOp) Update(ctx context.Context, l7PolicyID string, reqBody *L7PolicyUpdateRequest) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(l7PolicyID, "l7PolicyID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(l7policiesBasePathV1), l7PolicyID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, reqBody)
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
