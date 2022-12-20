package volumes

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

func updateURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func resourceActionURL(c *edgecloud.ServiceClient, id, action string) string {
	return c.ServiceURL(id, action)
}

func attachURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "attach")
}

func detachURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "detach")
}

func retypeURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "retype")
}

func extendURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "extend")
}

func revertURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "revert")
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
