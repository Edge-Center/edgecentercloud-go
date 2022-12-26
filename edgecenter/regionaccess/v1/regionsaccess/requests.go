package regionsaccess

import (
	"net/http"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToRegionAccessListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	ResellerID int `q:"reseller_id"`
	ClientID   int `q:"client_id"`
}

// ToRegionAccessListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToRegionAccessListQuery() (string, error) {
	if err := edgecloud.ValidateStruct(opts); err != nil {
		return "", err
	}
	q, err := edgecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create requets.
type CreateOptsBuilder interface {
	ToRegionAccessCreateMap() (map[string]interface{}, error)
}

// CreateOpts set parameters for Create operation
type CreateOpts struct {
	AccessAllEdgeRegions bool  `json:"access_all_edge_regions"`
	RegionIDs            []int `json:"region_ids"`
	ClientID             *int  `json:"client_id"`
	ResellerID           *int  `json:"reseller_id"`
}

// Validate
func (opts CreateOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

// ToRegionAccessCreateMap builds a request body form CreateOpts
func (opts CreateOpts) ToRegionAccessCreateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return edgecloud.BuildRequestBody(opts, "")
}

func List(c *edgecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToRegionAccessListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return RegionAccessPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func ListAll(c *edgecloud.ServiceClient, opts ListOptsBuilder) ([]RegionAccess, error) {
	pages, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractRegionsAccess(pages)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func Delete(c *edgecloud.ServiceClient, resellerID int) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, resellerID), &edgecloud.RequestOpts{OkCodes: []int{http.StatusNoContent}})
	return
}

func Create(c *edgecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRegionAccessCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(rootURL(c), b, &r.Body, &edgecloud.RequestOpts{OkCodes: []int{http.StatusOK}})
	return
}
