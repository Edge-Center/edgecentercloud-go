package laas

import (
	"strconv"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func statusURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL("status")
}

func usersURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL("users")
}

func topicsURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL("topics")
}

func deleteTopicURL(c *edgecloud.ServiceClient, name string) string {
	return c.ServiceURL("topics", name)
}

func kafkaURL(c *edgecloud.ServiceClient) string {
	return c.BaseServiceURL("laas", strconv.Itoa(c.RegionID), "kafka_hosts")
}

func openSearchURL(c *edgecloud.ServiceClient) string {
	return c.BaseServiceURL("laas", strconv.Itoa(c.RegionID), "opensearch_hosts")
}
