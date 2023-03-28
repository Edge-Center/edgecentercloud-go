package keystones

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/keystone/v1/types"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

func List(c *edgecloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, listURL(c), func(r pagination.PageResult) pagination.Page {
		return KeystonePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific keystone based on its unique ID.
func Get(c *edgecloud.ServiceClient, id int) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToKeystoneCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a keystone.
type CreateOpts struct {
	URL                       edgecloud.URL       `json:"url" required:"true" validate:"required"`
	State                     types.KeystoneState `json:"state" required:"true" validate:"required,enum"`
	KeystoneFederatedDomainID string              `json:"keystone_federated_domain_id" required:"true" validate:"required"`
	AdminPassword             string              `json:"admin_password,omitempty"`
}

// ToKeystoneCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToKeystoneCreateMap() (map[string]interface{}, error) {
	return edgecloud.BuildRequestBody(opts, "")
}

// Validate CreateOpts.
func (opts CreateOpts) Validate() error {
	return edgecloud.Validate.Struct(opts)
}

// Create accepts a CreateOpts struct and creates a new keystone using the values provided.
func Create(c *edgecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToKeystoneCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToKeystoneUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a keystone.
type UpdateOpts struct {
	URL                       *edgecloud.URL      `json:"url,omitempty" validate:"required_without_all=State KeystoneFederatedDomainID AdminPassword"`
	State                     types.KeystoneState `json:"state,omitempty" validate:"required_without_all=URL KeystoneFederatedDomainID AdminPassword"`
	KeystoneFederatedDomainID string              `json:"keystone_federated_domain_id,omitempty" validate:"required_without_all=URL State AdminPassword"`
	AdminPassword             string              `json:"admin_password,omitempty" validate:"required_without_all=State KeystoneFederatedDomainID AdminPassword"`
}

// ToKeystoneUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToKeystoneUpdateMap() (map[string]interface{}, error) {
	return edgecloud.BuildRequestBody(opts, "")
}

// Validate UpdateOpts.
func (opts UpdateOpts) Validate() error {
	return edgecloud.Validate.Struct(opts)
}

// Update accepts a UpdateOpts struct and updates an existing keystone using the values provided.
func Update(c *edgecloud.ServiceClient, id int, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToKeystoneUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, id), b, &r.Body, &edgecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// ListAll is a convenience function that returns all keystones.
func ListAll(client *edgecloud.ServiceClient) ([]Keystone, error) {
	pages, err := List(client).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractKeystones(pages)
	if err != nil {
		return nil, err
	}

	return all, nil
}
