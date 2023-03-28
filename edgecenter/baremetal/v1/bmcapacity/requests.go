package bmcapacity

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

// GetAvailableNodes retrieves available baremetal nodes.
func GetAvailableNodes(c *edgecloud.ServiceClient) (r GetAvailableNodesResult) {
	url := getAvailableNodesURL(c)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}
