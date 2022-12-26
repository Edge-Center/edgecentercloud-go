package pools

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func resourceURL(c *edgecloud.ServiceClient, clusterID, id string) string {
	return c.ServiceURL(clusterID, "pools", id)
}

func resourceActionURL(c *edgecloud.ServiceClient, clusterID, id, action string) string {
	return c.ServiceURL(clusterID, "pools", id, action)
}

func rootURL(c *edgecloud.ServiceClient, clusterID string) string {
	return c.ServiceURL(clusterID, "pools")
}

func getURL(c *edgecloud.ServiceClient, clusterID string, id string) string {
	return resourceURL(c, clusterID, id)
}

func listURL(c *edgecloud.ServiceClient, clusterID string) string {
	return rootURL(c, clusterID)
}

func createURL(c *edgecloud.ServiceClient, clusterID string) string {
	return rootURL(c, clusterID)
}

func updateURL(c *edgecloud.ServiceClient, clusterID string, id string) string {
	return resourceURL(c, clusterID, id)
}

func deleteURL(c *edgecloud.ServiceClient, clusterID string, id string) string {
	return resourceURL(c, clusterID, id)
}

func instancesURL(c *edgecloud.ServiceClient, clusterID string, id string) string {
	return resourceActionURL(c, clusterID, id, "instances")
}

func volumesURL(c *edgecloud.ServiceClient, clusterID string, id string) string {
	return resourceActionURL(c, clusterID, id, "volumes")
}
