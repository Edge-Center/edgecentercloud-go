package quotas

import (
	"strconv"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func getCombinedURL(c *edgecloud.ServiceClient) string {
	return c.BaseServiceURL("quotas_client")
}

func getGlobalURL(c *edgecloud.ServiceClient, clientID int) string {
	return c.BaseServiceURL("quotas_global", strconv.Itoa(clientID))
}

func getRegionURL(c *edgecloud.ServiceClient, clientID, regionID int) string {
	return c.BaseServiceURL("quotas_regional", strconv.Itoa(clientID), strconv.Itoa(regionID))
}
