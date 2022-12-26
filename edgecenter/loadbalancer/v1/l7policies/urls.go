package l7policies

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func resourceURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id)
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

func replaceURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func rulesRootURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "rules")
}

func rulesCreateURL(c *edgecloud.ServiceClient, id string) string {
	return rulesRootURL(c, id)
}

func rulesGetURL(c *edgecloud.ServiceClient, plid string, rlid string) string {
	return c.ServiceURL(plid, "rules", rlid)
}

func rulesListURL(c *edgecloud.ServiceClient, id string) string {
	return rulesRootURL(c, id)
}

func rulesDeleteURL(c *edgecloud.ServiceClient, plid string, rlid string) string {
	return c.ServiceURL(plid, "rules", rlid)
}

func rulesReplaceURL(c *edgecloud.ServiceClient, plid string, rlid string) string {
	return c.ServiceURL(plid, "rules", rlid)
}
