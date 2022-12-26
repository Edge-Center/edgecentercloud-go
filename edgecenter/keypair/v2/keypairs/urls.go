package keypairs

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func resourceURL(c *edgecloud.ServiceClient, id string) string {
	return c.BaseServiceURL("keypairs", id)
}

func rootURL(c *edgecloud.ServiceClient) string {
	return c.BaseServiceURL("keypairs")
}

func getURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func deleteURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}
