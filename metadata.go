package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	metadataPath     = "metadata"
	metadataItemPath = "metadata_item"
)

type MetadataDetailed struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	ReadOnly bool   `json:"read_only"`
}

type Metadata map[string]string

// MetadataCreateRequest represent a metadata create struct.
type MetadataCreateRequest struct {
	Metadata
}

type MetadataItemOptions struct {
	Key string `url:"key,omitempty" validate:"omitempty"`
}

// MetadataRoot represents a Metadata root.
type MetadataRoot struct {
	Count    int
	Metadata []MetadataDetailed `json:"results"`
}

// metadataList helper for same logic methods.
func metadataList(ctx context.Context, client *Client, id, resourcePath string) ([]MetadataDetailed, *Response, error) {
	path := fmt.Sprintf("%s/%s/%s", client.addProjectRegionPath(resourcePath), id, metadataPath)

	req, err := client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	metadata := new(MetadataRoot)
	resp, err := client.Do(ctx, req, metadata)
	if err != nil {
		return nil, resp, err
	}

	return metadata.Metadata, resp, err
}

// metadataCreate helper for same logic methods.
func metadataCreate(ctx context.Context, client *Client, id, resourcePath string, metadata *Metadata) (*Response, error) {
	path := fmt.Sprintf("%s/%s/%s", client.addProjectRegionPath(resourcePath), id, metadataPath)

	req, err := client.NewRequest(ctx, http.MethodPost, path, metadata)
	if err != nil {
		return nil, err
	}

	return client.Do(ctx, req, nil)
}

// metadataUpdate helper for same logic methods.
func metadataUpdate(ctx context.Context, client *Client, id, resourcePath string, metadata *Metadata) (*Response, error) {
	path := fmt.Sprintf("%s/%s/%s", client.addProjectRegionPath(resourcePath), id, metadataPath)

	req, err := client.NewRequest(ctx, http.MethodPut, path, metadata)
	if err != nil {
		return nil, err
	}

	return client.Do(ctx, req, nil)
}

// metadataDeleteItem helper for same logic methods.
func metadataDeleteItem(ctx context.Context, client *Client, id, resourcePath string, opts *MetadataItemOptions) (*Response, error) {
	path, err := addOptions(client.addProjectRegionPath(resourcePath), opts)
	if err != nil {
		return nil, err
	}

	path = fmt.Sprintf("%s/%s/%s", path, id, metadataItemPath)

	req, err := client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return client.Do(ctx, req, nil)
}

// metadataGetItem helper for same logic methods.
func metadataGetItem(ctx context.Context, client *Client, id, resourcePath string, opts *MetadataItemOptions) (*MetadataDetailed, *Response, error) {
	path, err := addOptions(client.addProjectRegionPath(resourcePath), opts)
	if err != nil {
		return nil, nil, err
	}

	path = fmt.Sprintf("%s/%s/%s", path, id, metadataItemPath)

	req, err := client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	metadata := new(MetadataDetailed)
	resp, err := client.Do(ctx, req, metadata)
	if err != nil {
		return nil, resp, err
	}

	return metadata, resp, err
}
