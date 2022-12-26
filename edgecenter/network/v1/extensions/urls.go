package extensions

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func resourceURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL("extensions", id)
}

func rootURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL("extensions")
}

func getURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}
