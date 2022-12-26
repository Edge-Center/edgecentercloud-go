package bmcapacity

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func rootURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL()
}

func getAvailableNodesURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}
