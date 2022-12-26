package regionsaccess

import (
	"strconv"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func rootURL(c *edgecloud.ServiceClient) string {
	return c.BaseServiceURL("reseller_region")
}

func resourceURL(c *edgecloud.ServiceClient, id int) string {
	return c.BaseServiceURL("reseller_region", strconv.Itoa(id))
}
