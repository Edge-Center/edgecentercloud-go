package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	quotasClientBasePathV2   = "/v2/quotas_client"
	quotasGlobalBasePathV2   = "/v2/quotas_global"
	quotasRegionalBasePathV2 = "/v2/quotas_regional"
)

const (
	quotasNotificationThreshold = "notification_threshold"
)

// QuotasService is an interface for creating and managing Quotas with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/quotas
type QuotasService interface {
	ListCombined(context.Context, *ListCombinedOptions) (*CombinedQuota, *Response, error)
	ListGlobal(context.Context, *ListGlobalOptions) (*Quota, *Response, error)
	ListRegional(context.Context, *ListRegionalOptions) (*Quota, *Response, error)
	DeleteNotificationThreshold(context.Context, int) (*Response, error)
	GetNotificationThreshold(context.Context, int) (*QuotaNotificationThreshold, *Response, error)
	UpdateNotificationThreshold(context.Context, int, *NotificationThresholdUpdateRequest) (*QuotaNotificationThreshold, *Response, error)
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

type QuotaNotificationThreshold struct {
	LastMessage CombinedQuota `json:"last_message"`
	LastSending string        `json:"last_sending"`
	Threshold   int           `json:"threshold"`
	ClientID    int           `json:"client_id"`
}

type NotificationThresholdUpdateRequest struct {
	LastMessage CombinedQuota `json:"last_message,omitempty" validate:"omitempty"`
	LastSending string        `json:"last_sending,omitempty" validate:"omitempty"`
	Threshold   int           `json:"threshold" validate:"required"`
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

// DeleteNotificationThreshold delete a client's quota notification threshold.
func (s *QuotasServiceOp) DeleteNotificationThreshold(ctx context.Context, clientID int) (*Response, error) {
	path := fmt.Sprintf("%s/%d/%s", quotasClientBasePathV2, clientID, quotasNotificationThreshold)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// GetNotificationThreshold get a client's quota notification threshold.
func (s *QuotasServiceOp) GetNotificationThreshold(ctx context.Context, clientID int) (*QuotaNotificationThreshold, *Response, error) {
	path := fmt.Sprintf("%s/%d/%s", quotasClientBasePathV2, clientID, quotasNotificationThreshold)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	quotaNotificationThreshold := new(QuotaNotificationThreshold)
	resp, err := s.client.Do(ctx, req, quotaNotificationThreshold)
	if err != nil {
		return nil, resp, err
	}

	return quotaNotificationThreshold, resp, err
}

// UpdateNotificationThreshold update or create a client's quota notification threshold.
func (s *QuotasServiceOp) UpdateNotificationThreshold(ctx context.Context, clientID int, reqBody *NotificationThresholdUpdateRequest) (*QuotaNotificationThreshold, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d/%s", quotasClientBasePathV2, clientID, quotasNotificationThreshold)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	quotaNotificationThreshold := new(QuotaNotificationThreshold)
	resp, err := s.client.Do(ctx, req, quotaNotificationThreshold)
	if err != nil {
		return nil, resp, err
	}

	return quotaNotificationThreshold, resp, err
}
