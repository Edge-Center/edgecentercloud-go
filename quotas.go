package edgecloud

import (
	"context"
	"net/http"
)

const (
	quotasClientBasePathV2   = "/v2/quotas_client"
	quotasGlobalBasePathV2   = "/v2/quotas_global"
	quotasRegionalBasePathV2 = "/v2/quotas_regional"
)

// QuotasService is an interface for creating and managing Quotas with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/quotas
type QuotasService interface {
	ListCombined(context.Context, *ListCombinedOptions) (*CombinedQuota, *Response, error)
	ListGlobal(context.Context, *ListGlobalOptions) (*Quota, *Response, error)
	ListRegional(context.Context, *ListRegionalOptions) (*Quota, *Response, error)
}

// QuotasServiceOp handles communication with Quotas methods of the EdgecenterCloud API.
type QuotasServiceOp struct {
	client *Client
}

var _ QuotasService = &QuotasServiceOp{}

type Quota map[string]int

type CombinedQuota struct {
	GlobalQuotas   Quota   `json:"global_quotas"`
	RegionalQuotas []Quota `json:"regional_quotas"`
}

// ListCombinedOptions specifies the query parameters to ListCombined method.
type ListCombinedOptions struct {
	ClientID int `url:"client_id,omitempty" validate:"omitempty"`
}

// ListGlobalOptions specifies the query parameters to ListGlobal method.
type ListGlobalOptions struct {
	ClientID int `url:"client_id"  required:"true" validate:"required"`
}

// ListRegionalOptions specifies the query parameters to ListRegional method.
type ListRegionalOptions struct {
	ClientID int `url:"client_id"  required:"true" validate:"required"`
	RegionID int `url:"region_id"  required:"true" validate:"required"`
}

// ListCombined get combined client quotas, regional and global.
func (s *QuotasServiceOp) ListCombined(ctx context.Context, opts *ListCombinedOptions) (*CombinedQuota, *Response, error) {
	path, err := addOptions(quotasClientBasePathV2, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	quota := new(CombinedQuota)
	resp, err := s.client.Do(ctx, req, quota)
	if err != nil {
		return nil, resp, err
	}

	return quota, resp, err
}

// ListGlobal get a global quota.
func (s *QuotasServiceOp) ListGlobal(ctx context.Context, opts *ListGlobalOptions) (*Quota, *Response, error) {
	path, err := addOptions(quotasGlobalBasePathV2, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	quota := new(Quota)
	resp, err := s.client.Do(ctx, req, quota)
	if err != nil {
		return nil, resp, err
	}

	return quota, resp, err
}

// ListRegional get a quota by region.
func (s *QuotasServiceOp) ListRegional(ctx context.Context, opts *ListRegionalOptions) (*Quota, *Response, error) {
	path, err := addOptions(quotasRegionalBasePathV2, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	quota := new(Quota)
	resp, err := s.client.Do(ctx, req, quota)
	if err != nil {
		return nil, resp, err
	}

	return quota, resp, err
}
