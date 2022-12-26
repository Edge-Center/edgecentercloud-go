package reservedfixedips

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

func createURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func listURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func deleteURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func switchVIPURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func connectedDeviceListURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "connected_devices")
}

func availableDeviceListURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "available_devices")
}

func portsToShareVIPURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "connected_devices")
}
