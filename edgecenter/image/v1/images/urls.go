package images

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func resourceURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id)
}

func rootURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL()
}

func getURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func updateURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func uploadURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}
