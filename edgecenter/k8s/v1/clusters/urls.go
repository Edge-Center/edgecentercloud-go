package clusters

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func versionsURL(c *edgecloud.ServiceClient) string {
	return c.BaseServiceURL("k8s", "versions")
}

func resourceURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id)
}

func rootURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL("")
}

func configURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "config")
}

func resizeURL(c *edgecloud.ServiceClient, clusterID, poolID string) string {
	return c.ServiceURL(clusterID, "pools", poolID, "resize")
}

func upgradeURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "upgrade")
}

func instancesURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "instances")
}

func certificatesURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "certificates")
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
