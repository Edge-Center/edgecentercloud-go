package lbpools

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func resourceURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id)
}

func resourceActionURL(c *edgecloud.ServiceClient, id string, action string) string {
	return c.ServiceURL(id, action)
}

func resourceActionDetailURL(c *edgecloud.ServiceClient, id string, action string, actorID string) string {
	return c.ServiceURL(id, action, actorID)
}

func rootURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL()
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

func updateURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func createMemberURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "member")
}

func deleteMemberURL(c *edgecloud.ServiceClient, id string, memberID string) string {
	return resourceActionDetailURL(c, id, "member", memberID)
}

func healthMonitorURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "healthmonitor")
}
