package edgecloud

import (
	"context"
)

const (
	imagesBasePathV1 = "/v1/images"
)

// ImagesService is an interface for creating and managing Images with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/images
type ImagesService interface {
	ImagesMetadata
}

type ImagesMetadata interface {
	MetadataList(context.Context, string) ([]MetadataDetailed, *Response, error)
	MetadataCreate(context.Context, string, *MetadataCreateRequest) (*Response, error)
	MetadataUpdate(context.Context, string, *MetadataCreateRequest) (*Response, error)
	MetadataDeleteItem(context.Context, string, *MetadataItemOptions) (*Response, error)
	MetadataGetItem(context.Context, string, *MetadataItemOptions) (*MetadataDetailed, *Response, error)
}

// ImagesServiceOp handles communication with Images methods of the EdgecenterCloud API.
type ImagesServiceOp struct {
	client *Client
}

var _ ImagesService = &ImagesServiceOp{}

// MetadataList security group detailed metadata items.
func (s *ImagesServiceOp) MetadataList(ctx context.Context, imageID string) ([]MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(imageID, "imageID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataList(ctx, s.client, imageID, imagesBasePathV1)
}

// MetadataCreate or update security group metadata.
func (s *ImagesServiceOp) MetadataCreate(ctx context.Context, imageID string, reqBody *MetadataCreateRequest) (*Response, error) {
	if resp, err := isValidUUID(imageID, "imageID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataCreate(ctx, s.client, imageID, imagesBasePathV1, reqBody)
}

// MetadataUpdate security group metadata.
func (s *ImagesServiceOp) MetadataUpdate(ctx context.Context, imageID string, reqBody *MetadataCreateRequest) (*Response, error) {
	if resp, err := isValidUUID(imageID, "imageID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataUpdate(ctx, s.client, imageID, imagesBasePathV1, reqBody)
}

// MetadataDeleteItem a security group metadata item by key.
func (s *ImagesServiceOp) MetadataDeleteItem(ctx context.Context, imageID string, opts *MetadataItemOptions) (*Response, error) {
	if resp, err := isValidUUID(imageID, "imageID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataDeleteItem(ctx, s.client, imageID, imagesBasePathV1, opts)
}

// MetadataGetItem security group detailed metadata.
func (s *ImagesServiceOp) MetadataGetItem(ctx context.Context, imageID string, opts *MetadataItemOptions) (*MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(imageID, "imageID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataGetItem(ctx, s.client, imageID, imagesBasePathV1, opts)
}
