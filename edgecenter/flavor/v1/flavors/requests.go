package flavors

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToFlavorListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the cluster templates attributes you want to see returned. SortKey allows you to sort
// by a particular cluster templates attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	IncludePrices *bool `q:"include_prices"`
}

// ToFlavorListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFlavorListQuery() (string, error) {
	q, err := edgecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List retrieves list of flavors.
func List(c *edgecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToFlavorListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return FlavorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll retrieves list of flavors.
func ListAll(c *edgecloud.ServiceClient, opts ListOptsBuilder) ([]Flavor, error) {
	results, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractFlavors(results)
}

// IDFromName is a convenience function that returns a flavor ID, given its name.
func IDFromName(client *edgecloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""

	pages, err := List(client, nil).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ExtractFlavors(pages)
	if err != nil {
		return "", err
	}

	for _, s := range all {
		if s.FlavorName == name {
			count++
			id = s.FlavorID
		}
	}

	switch count {
	case 0:
		return "", edgecloud.ResourceNotFoundError{Name: name, ResourceType: "flavors"}
	case 1:
		return id, nil
	default:
		return "", edgecloud.MultipleResourcesFoundError{Name: name, Count: count, ResourceType: "flavors"}
	}
}
