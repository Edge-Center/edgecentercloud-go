package edgecloud

import (
	"context"
	"net/http"
)

const (
	flavorsBasePathV1 = "/v1/flavors"
)

// FlavorsService is an interface for creating and managing Flavors with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/flavors
type FlavorsService interface {
	List(context.Context, *FlavorListOptions) ([]Flavor, *Response, error)
}

// FlavorsServiceOp handles communication with Flavors methods of the EdgecenterCloud API.
type FlavorsServiceOp struct {
	client *Client
}

var _ FlavorsService = &FlavorsServiceOp{}

// Flavor represents an EdgecenterCloud Flavor.
type Flavor struct {
	FlavorID            string              `json:"flavor_id"`
	FlavorName          string              `json:"flavor_name"`
	VCPUS               int                 `json:"vcpus"`
	RAM                 int                 `json:"ram"`
	HardwareDescription HardwareDescription `json:"hardware_description,omitempty"`
	Disabled            bool                `json:"disabled"`
	ResourceClass       string              `json:"resource_class"`
}

type HardwareDescription struct {
	CPU         string `json:"cpu"`
	RAM         string `json:"ram"`
	Disk        string `json:"disk"`
	Network     string `json:"network"`
	GPU         string `json:"gpu"`
	IPU         string `json:"ipu,omitempty"`
	PoplarCount string `json:"poplar_count,omitempty"`
	SgxEPCSize  string `json:"sgx_epc_size"`
}

// FlavorListOptions specifies the optional query parameters to List method.
type FlavorListOptions struct {
	IncludePrices  bool `url:"include_prices,omitempty"  validate:"omitempty"`
	Disabled       bool `url:"disabled,omitempty"  validate:"omitempty"`
	ExcludeWindows bool `url:"exclude_windows,omitempty"  validate:"omitempty"`
}

// floatingipsRoot represents a FloatingIPs root.
type flavorsRoot struct {
	Count   int
	Flavors []Flavor `json:"results"`
}

// List get flavors.
func (s *FlavorsServiceOp) List(ctx context.Context, opts *FlavorListOptions) ([]Flavor, *Response, error) {
	if err := s.client.Validate(); err != nil {
		return nil, nil, err
	}

	path := s.client.addServicePath(flavorsBasePathV1)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(flavorsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Flavors, resp, err
}
