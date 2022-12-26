package apptemplates

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

// List retrieves list of app templates
func List(c *edgecloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, rootURL(c), func(r pagination.PageResult) pagination.Page {
		return AppTemplatePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll retrieves list of app templates
func ListAll(c *edgecloud.ServiceClient) ([]AppTemplate, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractAppTemplates(page)
}

// Get retrieves a specific app template based on its unique ID.
func Get(c *edgecloud.ServiceClient, id string) (r GetResult) {
	url := resourceURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}
