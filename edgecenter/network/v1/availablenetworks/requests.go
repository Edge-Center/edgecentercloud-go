package availablenetworks

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToAvailableNetworkListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	NetworkID   string `q:"network_id"`
	NetworkType string `q:"network_type"`
}

// ToAvailableNetworkListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToAvailableNetworkListQuery() (string, error) {
	q, err := edgecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(c *edgecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToAvailableNetworkListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return NetworkPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll is a convenience function that returns all networks.
func ListAll(client *edgecloud.ServiceClient, opts ListOptsBuilder) ([]Network, error) {
	pages, err := List(client, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractNetworks(pages)
	if err != nil {
		return nil, err
	}

	return all, nil
}
