package metadata

import (
	"fmt"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func ResourceActionURL(c *edgecloud.ServiceClient, id string, action string) string {
	return c.ServiceURL(id, action)
}

func ResourceMetadataURL(c *edgecloud.ServiceClient, id string) string {
	return ResourceActionURL(c, id, "metadata")
}

func ResourceMetadataItemURL(c *edgecloud.ServiceClient, id string, key string) string {
	return ResourceActionURL(c, id, fmt.Sprintf("metadata_item?key=%s", key))
}
