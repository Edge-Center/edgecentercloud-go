package apptemplates

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func rootURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL()
}

func resourceURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id)
}
