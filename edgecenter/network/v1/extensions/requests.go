package extensions

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

func List(c *edgecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ExtensionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific extension based on its alias.
func Get(c *edgecloud.ServiceClient, alias string) (r GetResult) {
	url := getURL(c, alias)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

func ListAll(c *edgecloud.ServiceClient) ([]Extension, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractExtensions(page)
}
