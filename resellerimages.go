package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	resellerImageBasePathV1 = "/v1/reseller_image"
)

// ResellerImageService is an interface for managing Reseller images with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud_resellers#tag/Images
type ResellerImageService interface {
	List(context.Context, int) (*ResellerImageList, *Response, error)
	Delete(context.Context, int) (*Response, error)
	Update(context.Context, *ResellerImageUpdateRequest) (*ResellerImage, *Response, error)
}

// ResellerImageServiceOp handles communication with Reseller images methods of the EdgecenterCloud API.
type ResellerImageServiceOp struct {
	client *Client
}

var _ ResellerImageService = &ResellerImageServiceOp{}

type ImageIDs []string

// ResellerImageList represents an EdgecenterCloud reseller images list.
type ResellerImageList struct {
	Count   int             `json:"count"`
	Results []ResellerImage `json:"results"`
}

// ResellerImage represents an EdgecenterCloud reseller image.
type ResellerImage struct {
	ImageIDs  *ImageIDs `json:"image_ids"`
	RegionID  int       `json:"region_id"`
	CreatedAt string    `json:"created_at,omitempty"`
	UpdatedAt string    `json:"updated_at,omitempty"`
}

// ResellerImageUpdateRequest represents a request to update available image list for a reseller.
type ResellerImageUpdateRequest struct {
	ImageIDs   *ImageIDs `json:"image_ids"`
	RegionID   int       `json:"region_id"`
	ResellerID int       `json:"reseller_id"`
	CreatedAt  string    `json:"created_at,omitempty"`
	UpdatedAt  string    `json:"updated_at,omitempty"`
}

// List get available image IDs limits for a reseller.
func (s *ResellerImageServiceOp) List(ctx context.Context, resellerID int) (*ResellerImageList, *Response, error) {
	pathReq := fmt.Sprintf("%s/%d", resellerImageBasePathV1, resellerID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, pathReq, nil)
	if err != nil {
		return nil, nil, err
	}

	ril := new(ResellerImageList)
	resp, err := s.client.Do(ctx, req, ril)
	if err != nil {
		return nil, resp, err
	}

	return ril, resp, nil
}

// Update set or update available image list for a reseller.
func (s *ResellerImageServiceOp) Update(ctx context.Context, reqBody *ResellerImageUpdateRequest) (*ResellerImage, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, resellerImageBasePathV1, reqBody)
	if err != nil {
		return nil, nil, err
	}

	rImage := new(ResellerImage)
	resp, err := s.client.Do(ctx, req, rImage)
	if err != nil {
		return nil, resp, err
	}

	return rImage, resp, nil
}

// Delete image limits for reseller clients.
func (s *ResellerImageServiceOp) Delete(ctx context.Context, resellerID int) (*Response, error) {
	pathReq := fmt.Sprintf("%s/%d", resellerImageBasePathV1, resellerID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, pathReq, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
