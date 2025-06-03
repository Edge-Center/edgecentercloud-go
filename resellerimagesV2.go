package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	resellerImageBasePathV2 = "/v2/reseller_image"
)

type (
	EntityID   = int
	EntityType = string
)

const (
	ResellerType EntityType = "reseller"
	ClientType   EntityType = "client"
	ProjectType  EntityType = "project"
)

// ResellerImageV2Service is an interface for managing Reseller images with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud_resellers#tag/Images
type ResellerImageV2Service interface {
	List(context.Context, EntityType, EntityID) (*ResellerImageV2List, *Response, error)
	ListByRole(context.Context) (*ResellerImageV2List, *Response, error)
	Delete(context.Context, EntityType, EntityID, *ResellerImageV2DeleteOptions) (*Response, error)
	Update(context.Context, *ResellerImageV2UpdateRequest) (*ResellerImageV2, *Response, error)
}

// ResellerImageV2ServiceOp handles communication with Reseller images methods of the EdgecenterCloud API.
type ResellerImageV2ServiceOp struct {
	client *Client
}

var _ ResellerImageV2Service = &ResellerImageV2ServiceOp{}

// ResellerImageV2List represents an EdgecenterCloud reseller images list.
type ResellerImageV2List struct {
	Count   int               `json:"count"`
	Results []ResellerImageV2 `json:"results"`
}

type ImageIDs []string

// ResellerImageV2 represents an EdgecenterCloud reseller images.
type ResellerImageV2 struct {
	ImageIDs   *ImageIDs  `json:"image_ids"`
	RegionID   int        `json:"region_id"`
	EntityID   EntityID   `json:"entity_id"`
	EntityType EntityType `json:"entity_type"`
	CreatedAt  string     `json:"created_at,omitempty"`
	UpdatedAt  string     `json:"updated_at,omitempty"`
}

// ResellerImageV2UpdateRequest represents a request to update available image list for a reseller, client, project.
type ResellerImageV2UpdateRequest ResellerImageV2

type ResellerImageV2DeleteOptions struct {
	RegionID int `url:"region_id"`
}

// ListByRole get available shared image IDs by role in APIKey.
func (s *ResellerImageV2ServiceOp) ListByRole(ctx context.Context) (*ResellerImageV2List, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, resellerImageBasePathV2, nil)
	if err != nil {
		return nil, nil, err
	}

	ril := new(ResellerImageV2List)
	resp, err := s.client.Do(ctx, req, ril)
	if err != nil {
		return nil, resp, err
	}

	return ril, resp, nil
}

// List get available image IDs limits for a reseller, client, project.
func (s *ResellerImageV2ServiceOp) List(ctx context.Context, entityType EntityType, entityID EntityID) (*ResellerImageV2List, *Response, error) {
	pathReq := fmt.Sprintf("%s/%s/%d", resellerImageBasePathV2, entityType, entityID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, pathReq, nil)
	if err != nil {
		return nil, nil, err
	}

	ril := new(ResellerImageV2List)
	resp, err := s.client.Do(ctx, req, ril)
	if err != nil {
		return nil, resp, err
	}

	return ril, resp, nil
}

// Update set or update available image list for a reseller, client, project.
func (s *ResellerImageV2ServiceOp) Update(ctx context.Context, reqBody *ResellerImageV2UpdateRequest) (*ResellerImageV2, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, resellerImageBasePathV2, reqBody)
	if err != nil {
		return nil, nil, err
	}

	rImage := new(ResellerImageV2)
	resp, err := s.client.Do(ctx, req, rImage)
	if err != nil {
		return nil, resp, err
	}

	return rImage, resp, nil
}

// Delete image limits for reseller, client, project.
func (s *ResellerImageV2ServiceOp) Delete(ctx context.Context, entityType EntityType, entityID EntityID, opts *ResellerImageV2DeleteOptions) (*Response, error) {
	pathReq := fmt.Sprintf("%s/%s/%d", resellerImageBasePathV2, entityType, entityID)

	pathReq, err := addOptions(pathReq, opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, pathReq, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
