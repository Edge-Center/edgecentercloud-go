package edgecloud

import (
	"context"
	"net/http"
)

const (
	flavorsBasePathV1   = "/v1/flavors"
	bmflavorsBasePathV1 = "/v1/bmflavors"
)

// FlavorsService is an interface for creating and managing Flavors with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/flavors
type FlavorsService interface {
	List(context.Context, *FlavorListOptions) ([]Flavor, *Response, error)
	ListBaremetal(context.Context, *FlavorListOptions) ([]Flavor, *Response, error)
	ListBaremetalForClient(context.Context, *FlavorListOptions) ([]Flavor, *Response, error)
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

// FlavorsOptions specifies the optional query parameters to Get loadbalancer or instance flavor method.
type FlavorsOptions struct {
	IncludePrices bool `url:"include_prices,omitempty"  validate:"omitempty"`
}

// flavorsRoot represents a Flavors root.
type flavorsRoot struct {
	Count   int
	Flavors []Flavor `json:"results"`
}

// List get flavors.
func (s *FlavorsServiceOp) List(ctx context.Context, opts *FlavorListOptions) ([]Flavor, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(flavorsBasePathV1)
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

// ListBaremetal get baremetal flavors.
func (s *FlavorsServiceOp) ListBaremetal(ctx context.Context, opts *FlavorListOptions) ([]Flavor, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(bmflavorsBasePathV1)
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

// ListBaremetalForClient get baremetal flavors from default project for current client.
func (s *FlavorsServiceOp) ListBaremetalForClient(ctx context.Context, opts *FlavorListOptions) ([]Flavor, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addRegionPath(bmflavorsBasePathV1)
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
