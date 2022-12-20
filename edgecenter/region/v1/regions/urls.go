package regions

import (
	"strconv"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func resourceURL(c *edgecloud.ServiceClient, id int) string {
	return c.BaseServiceURL("regions", strconv.Itoa(id))
}

func rootURL(c *edgecloud.ServiceClient) string {
	return c.BaseServiceURL("regions")
}

func getURL(c *edgecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}

func listURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func updateURL(c *edgecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}
