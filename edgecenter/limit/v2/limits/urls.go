package limits

import (
	"strconv"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func resourceURL(c *edgecloud.ServiceClient, id int) string {
	return c.BaseServiceURL("limits_request", strconv.Itoa(id))
}

func rootURL(c *edgecloud.ServiceClient) string {
	return c.BaseServiceURL("limits_request")
}

func createURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func getURL(c *edgecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}

func listURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func deleteURL(c *edgecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}
