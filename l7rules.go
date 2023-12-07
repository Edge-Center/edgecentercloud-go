package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	l7rulesPath = "rules"
)

// L7RulesService is an interface for creating and managing L7Rules with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/l7rules
type L7RulesService interface {
	List(context.Context, string) ([]L7Rule, *Response, error)
	Create(context.Context, string, *L7RuleCreateRequest) (*TaskResponse, *Response, error)

	Delete(context.Context, string, string) (*TaskResponse, *Response, error)
	Get(context.Context, string, string) (*L7Rule, *Response, error)
	Update(context.Context, string, string, *L7RuleUpdateRequest) (*TaskResponse, *Response, error)
}

// L7RulesServiceOp handles communication with L7Rules methods of the EdgecenterCloud API.
type L7RulesServiceOp struct {
	client *Client
}

var _ L7RulesService = &L7RulesServiceOp{}

type L7Rule struct {
	RegionID           int               `json:"region_id"`
	ProjectID          int               `json:"project_id"`
	TaskID             string            `json:"task_id"`
	ID                 string            `json:"id"`
	Region             string            `json:"region"`
	ProvisioningStatus string            `json:"provisioning_status"`
	OperatingStatus    string            `json:"operating_status"`
	Tags               []string          `json:"tags"`
	Value              string            `json:"value"`
	Key                string            `json:"key"`
	Invert             bool              `json:"invert"`
	Type               L7RuleType        `json:"type"`
	CompareType        L7RuleCompareType `json:"compare_type"`
}

type L7RuleType string

const (
	L7RuleTypeCookie          L7RuleType = "COOKIE"
	L7RuleTypeFyleType        L7RuleType = "FILE_TYPE"
	L7RuleTypeHeader          L7RuleType = "HEADER"
	L7RuleTypeHostName        L7RuleType = "HOST_NAME"
	L7RuleTypePath            L7RuleType = "PATH"
	L7RuleTypeSSLConnHasCert  L7RuleType = "SSL_CONN_HAS_CERT"
	L7RuleTypeSSLVerifyResult L7RuleType = "SSL_VERIFY_RESULT"
	L7RuleTypeSSLDNField      L7RuleType = "SSL_DN_FIELD"
)

type L7RuleCompareType string

const (
	L7RuleCompareTypeContains   L7RuleCompareType = "CONTAINS"
	L7RuleCompareTypeEndsWith   L7RuleCompareType = "ENDS_WITH"
	L7RuleCompareTypeEqualTo    L7RuleCompareType = "EQUAL_TO"
	L7RuleCompareTypeRegex      L7RuleCompareType = "REGEX"
	L7RuleCompareTypeStartsWith L7RuleCompareType = "STARTS_WITH"
)

type L7RuleUpdateRequest struct {
	Tags        []string          `json:"tags,omitempty"`
	CompareType L7RuleCompareType `json:"compare_type,omitempty"`
	Value       string            `json:"value,omitempty"`
	Key         string            `json:"key,omitempty"`
	Type        L7RuleType        `json:"type,omitempty"`
	Invert      bool              `json:"invert,omitempty"`
}

type L7RuleCreateRequest struct {
	Tags        []string          `json:"tags,omitempty"`
	CompareType L7RuleCompareType `json:"compare_type" required:"true" validate:"required"`
	Value       string            `json:"value" required:"true" validate:"required"`
	Key         string            `json:"key,omitempty"`
	Type        L7RuleType        `json:"type" required:"true" validate:"required"`
	Invert      bool              `json:"invert,omitempty"`
}

// l7rulesRoot represents a L7Rule root.
type l7rulesRoot struct {
	Count   int
	L7Rules []L7Rule `json:"results"`
}

// List get L7Rules.
func (s *L7RulesServiceOp) List(ctx context.Context, l7PolicyID string) ([]L7Rule, *Response, error) {
	if resp, err := isValidUUID(l7PolicyID, "l7PolicyID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(l7policiesBasePathV1), l7PolicyID, l7rulesPath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(l7rulesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.L7Rules, resp, err
}

// Create a L7Rule.
func (s *L7RulesServiceOp) Create(ctx context.Context, l7PolicyID string, reqBody *L7RuleCreateRequest) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(l7PolicyID, "l7PolicyID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(l7policiesBasePathV1), l7PolicyID, l7rulesPath)

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

// Delete a L7Rule.
func (s *L7RulesServiceOp) Delete(ctx context.Context, l7PolicyID string, l7RuleID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(l7PolicyID, "l7PolicyID"); err != nil {
		return nil, resp, err
	}

	if resp, err := isValidUUID(l7RuleID, "l7RuleID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s/%s", s.client.addProjectRegionPath(l7policiesBasePathV1), l7PolicyID, l7rulesPath, l7RuleID)

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

// Get a L7Rule.
func (s *L7RulesServiceOp) Get(ctx context.Context, l7PolicyID string, l7RuleID string) (*L7Rule, *Response, error) {
	if resp, err := isValidUUID(l7PolicyID, "l7PolicyID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s/%s", s.client.addProjectRegionPath(l7policiesBasePathV1), l7PolicyID, l7rulesPath, l7RuleID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	l7Rule := new(L7Rule)
	resp, err := s.client.Do(ctx, req, l7Rule)
	if err != nil {
		return nil, resp, err
	}

	return l7Rule, resp, err
}

// Update replace L7Rule properties.
func (s *L7RulesServiceOp) Update(ctx context.Context, l7PolicyID string, l7RuleID string, reqBody *L7RuleUpdateRequest) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(l7PolicyID, "l7PolicyID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s/%s", s.client.addProjectRegionPath(l7policiesBasePathV1), l7PolicyID, l7rulesPath, l7RuleID)

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
