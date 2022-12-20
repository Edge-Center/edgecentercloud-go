package networks

import (
	"fmt"
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func resourceURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id)
}

func resourceActionURL(c *edgecloud.ServiceClient, id string, action string) string {
	return c.ServiceURL(id, action)
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

func listInstancePortsURL(c *edgecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "ports")
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

func metadataURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "metadata")
}
func metadataItemURL(c *edgecloud.ServiceClient, id string, key string) string {
	return resourceActionURL(c, id, fmt.Sprintf("metadata_item?key=%s", key))
}
