package availablenetworks

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func rootURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL()
}

func listURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}
