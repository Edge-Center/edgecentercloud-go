package resources

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func resourceActionURL(c *edgecloud.ServiceClient, stackID, resourceName, action string) string {
	return c.ServiceURL("stacks", stackID, "resources", resourceName, action)
}

func resourceURL(c *edgecloud.ServiceClient, stackID, resourceName string) string {
	return c.ServiceURL("stacks", stackID, "resources", resourceName)
}

func rootURL(c *edgecloud.ServiceClient, stackID string) string {
	return c.ServiceURL("stacks", stackID, "resources")
}

func MetadataURL(c *edgecloud.ServiceClient, stackID, resourceName string) string {
	return resourceActionURL(c, stackID, resourceName, "metadata")
}

func SignalURL(c *edgecloud.ServiceClient, stackID, resourceName string) string {
	return resourceActionURL(c, stackID, resourceName, "signal")
}

func listURL(c *edgecloud.ServiceClient, stackID string) string {
	return rootURL(c, stackID)
}

func getURL(c *edgecloud.ServiceClient, stackID, resourceName string) string {
	return resourceURL(c, stackID, resourceName)
}

func markUnhealthyURL(c *edgecloud.ServiceClient, stackID, resourceName string) string {
	return resourceURL(c, stackID, resourceName)
}
