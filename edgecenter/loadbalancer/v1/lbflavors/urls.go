package lbflavors

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func listURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL()
}
