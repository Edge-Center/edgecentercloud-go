package lifecyclepolicy

import (
	"strconv"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func getURL(c *edgecloud.ServiceClient, id int) string {
	return c.ServiceURL(strconv.Itoa(id))
}

func listURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL()
}

func deleteURL(c *edgecloud.ServiceClient, id int) string {
	return c.ServiceURL(strconv.Itoa(id))
}

func createURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL()
}

func updateURL(c *edgecloud.ServiceClient, id int) string {
	return c.ServiceURL(strconv.Itoa(id))
}

func addVolumesURL(c *edgecloud.ServiceClient, id int) string {
	return c.ServiceURL(strconv.Itoa(id), "add_volumes_to_policy")
}

func removeVolumesURL(c *edgecloud.ServiceClient, id int) string {
	return c.ServiceURL(strconv.Itoa(id), "remove_volumes_from_policy")
}

func addSchedulesURL(c *edgecloud.ServiceClient, id int) string {
	return c.ServiceURL(strconv.Itoa(id), "add_schedules")
}

func removeSchedulesURL(c *edgecloud.ServiceClient, id int) string {
	return c.ServiceURL(strconv.Itoa(id), "remove_schedules")
}

func estimateURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL("estimate_max_policy_usage")
}
