package quotas

import (
	"strconv"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func getCombinedURL(c *edgecloud.ServiceClient) string {
	return c.BaseServiceURL("client_quotas")
}

func getGlobalURL(c *edgecloud.ServiceClient, clientID int) string {
	return c.BaseServiceURL("global_quotas", strconv.Itoa(clientID))
}

func getRegionURL(c *edgecloud.ServiceClient, clientID, regionID int) string {
	return c.BaseServiceURL("regional_quotas", strconv.Itoa(clientID), strconv.Itoa(regionID))
}
