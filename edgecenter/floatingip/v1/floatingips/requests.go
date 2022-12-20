package floatingips

import (
	"net"

	"github.com/Edge-Center/edgecentercloud-go/edgecenter/task/v1/tasks"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

func List(c *edgecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return FloatingIPPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific floating ip based on its unique ID.
func Get(c *edgecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder adds additional parameters to the request.
type CreateOptsBuilder interface {
	ToFloatingIPCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a floating ip.
type CreateOpts struct {
	PortID         string `json:"port_id,omitempty"`
	FixedIPAddress net.IP `json:"fixed_ip_address,omitempty"`
}

// ToFloatingIPCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToFloatingIPCreateMap() (map[string]interface{}, error) {
	return edgecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new floating ip using the values provided.
func Create(c *edgecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToFloatingIPCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToFloatingIPUpdateMap() (map[string]interface{}, error)
}

// Delete accepts a unique ID and deletes the floating ip associated with it.
func Delete(c *edgecloud.ServiceClient, floatingID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, floatingID), &r.Body, nil)
	return
}

// ListAll returns all floating IPs
func ListAll(c *edgecloud.ServiceClient) ([]FloatingIPDetail, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractFloatingIPs(page)
}

// Assign accepts a CreateOpts struct and assign floating IP.
func Assign(c *edgecloud.ServiceClient, floatingIPID string, opts CreateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToFloatingIPCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(assignURL(c, floatingIPID), b, &r.Body, nil)
	return
}

// UnAssign
func UnAssign(c *edgecloud.ServiceClient, floatingIPID string) (r UpdateResult) {
	_, r.Err = c.Post(unAssignURL(c, floatingIPID), nil, &r.Body, nil)
	return
}
