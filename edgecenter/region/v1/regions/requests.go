package regions

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/region/v1/types"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

func List(c *edgecloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, listURL(c), func(r pagination.PageResult) pagination.Page {
		return RegionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific region based on its unique ID.
func Get(c *edgecloud.ServiceClient, id int) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToRegionCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a region.
type CreateOpts struct {
	DisplayName       string             `json:"display_name" required:"true" validate:"required"`
	KeystoneName      string             `json:"keystone_name" required:"true" validate:"required"`
	State             types.RegionState  `json:"state" required:"true" validate:"required,enum"`
	EndpointType      types.EndpointType `json:"endpoint_type,omitempty" validate:"omitempty,enum"`
	ExternalNetworkID string             `json:"external_network_id" required:"true" validate:"required,uuid4"`
	SpiceProxyURL     *edgecloud.URL     `json:"spice_proxy_url,omitempty"`
	KeystoneID        int                `json:"keystone_id" required:"true" validate:"required"`
}

// ToRegionCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToRegionCreateMap() (map[string]interface{}, error) {
	return edgecloud.BuildRequestBody(opts, "")
}

// Validate CreateOpts.
func (opts CreateOpts) Validate() error {
	return edgecloud.Validate.Struct(opts)
}

// Create accepts a CreateOpts struct and creates a new region using the values provided.
func Create(c *edgecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRegionCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToRegionUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a region.
type UpdateOpts struct {
	DisplayName       string             `json:"display_name,omitempty" validate:"required_without_all=State EndpointType ExternalNetworkID SpiceProxyURL,omitempty"`
	State             types.RegionState  `json:"state,omitempty" validate:"required_without_all=DisplayName EndpointType ExternalNetworkID SpiceProxyURL,omitempty,enum"`
	EndpointType      types.EndpointType `json:"endpoint_type,omitempty" validate:"required_without_all=DisplayName State ExternalNetworkID SpiceProxyURL,omitempty,enum"`
	ExternalNetworkID string             `json:"external_network_id,omitempty" validate:"required_without_all=DisplayName State EndpointType SpiceProxyURL,omitempty,uuid4"`
	SpiceProxyURL     *edgecloud.URL     `json:"spice_proxy_url,omitempty" validate:"required_without_all=DisplayName State EndpointType ExternalNetworkID,omitempty"`
}

// ToRegionUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToRegionUpdateMap() (map[string]interface{}, error) {
	return edgecloud.BuildRequestBody(opts, "")
}

// Validate UpdateOpts.
func (opts UpdateOpts) Validate() error {
	return edgecloud.Validate.Struct(opts)
}

// Update accepts a UpdateOpts struct and updates an existing region using the values provided.
func Update(c *edgecloud.ServiceClient, id int, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToRegionUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, id), b, &r.Body, &edgecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// ListAll is a convenience function that returns all regions.
func ListAll(client *edgecloud.ServiceClient) ([]Region, error) {
	pages, err := List(client).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractRegions(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}
