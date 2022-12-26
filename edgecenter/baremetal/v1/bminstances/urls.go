package bminstances

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func rootURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL()
}

func listURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func rebuildURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "rebuild")
}
