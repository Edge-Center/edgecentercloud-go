package availablefloatingips

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/floatingip/v1/floatingips"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

func List(c *edgecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return floatingips.FloatingIPPage{
			LinkedPageBase: pagination.LinkedPageBase{PageResult: r},
		}
	})
}

// ListAll returns all floating IPs
func ListAll(c *edgecloud.ServiceClient) ([]floatingips.FloatingIPDetail, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return floatingips.ExtractFloatingIPs(page)
}
