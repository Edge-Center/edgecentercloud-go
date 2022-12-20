package stacks

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func resourceURL(c *edgecloud.ServiceClient, stackID string) string {
	return c.ServiceURL("stacks", stackID)
}

func rootURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL("stacks")
}

func getURL(c *edgecloud.ServiceClient, stackID string) string {
	return resourceURL(c, stackID)
}

func updateURL(c *edgecloud.ServiceClient, stackID string) string {
	return resourceURL(c, stackID)
}

func deleteURL(c *edgecloud.ServiceClient, stackID string) string {
	return resourceURL(c, stackID)
}

func listURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}
