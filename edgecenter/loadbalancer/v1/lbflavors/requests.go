package lbflavors

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

func List(c *edgecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return FlavorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll returns all LB flavors
func ListAll(c *edgecloud.ServiceClient) ([]Flavor, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractFlavors(page)
}
