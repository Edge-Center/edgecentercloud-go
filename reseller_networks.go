package edgecloud

import (
	"context"
	"net/http"
)

const (
	resellerNetworksBasePathV1 = "/v1/reseller_networks"
)

// ResellerNetworksService
// Returns the list of networks with subnet details that are available to the reseller and its clients in all regions..
// See: https://apidocs.edgecenter.ru/cloud_resellers#tag/Networks/operation/ResellerNetworkViewSet.get
type ResellerNetworksService interface {
	List(context.Context, *ResellerNetworksListRequest) (*ResellerNetworks, *Response, error)
}

// ResellerNetworksServiceOp handles communication with Reseller Networks methods of the EdgecenterCloud API.
type ResellerNetworksServiceOp struct {
	client *Client
}

var _ ResellerNetworksService = &ResellerNetworksServiceOp{}

// ResellerNetwork represents an EdgecenterCloud reseller network.
type ResellerNetwork struct {
	CreatedAt      string             `json:"created_at"`
	Default        bool               `json:"default"`
	External       bool               `json:"external"`
	Shared         bool               `json:"shared"`
	ID             string             `json:"id"`
	MTU            int                `json:"mtu"`
	Name           string             `json:"name"`
	RegionID       int                `json:"region_id"`
	Region         string             `json:"region"`
	Type           string             `json:"type"`
	Subnets        []Subnetwork       `json:"subnets"`
	CreatorTaskID  string             `json:"creator_task_id"`
	TaskID         string             `json:"task_id"`
	SegmentationID int                `json:"segmentation_id"`
	UpdatedAt      string             `json:"updated_at"`
	Metadata       []MetadataDetailed `json:"metadata,omitempty"`
	ClientID       int                `json:"client_id"`
	ProjectID      int                `json:"project_id"`
}

type ResellerNetworks struct {
	Count   int               `json:"count"`
	Results []ResellerNetwork `json:"results"`
}

type ResellerNetworksListRequest struct {
	NetworkType string `url:"network_type,omitempty"`
	OrderBy     string `url:"order_by,omitempty"`
	Shared      bool   `url:"shared,omitempty"`
	MetadataKV  string `url:"metadata_kv,omitempty"`
	MetadataK   string `url:"metadata_k,omitempty"`
}

// List get reseller networks.
func (s *ResellerNetworksServiceOp) List(ctx context.Context, opts *ResellerNetworksListRequest) (*ResellerNetworks, *Response, error) {
	path, err := addOptions(resellerNetworksBasePathV1, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ResellerNetworks)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
