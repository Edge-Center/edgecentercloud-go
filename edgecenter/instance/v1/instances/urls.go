package instances

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

func renameInstanceURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *edgecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func resourceActionURL(c *edgecloud.ServiceClient, id string, action string) string {
	return c.ServiceURL(id, action)
}

func resourceActionDetailsURL(c *edgecloud.ServiceClient, id string, action string, actionID string) string {
	return c.ServiceURL(id, action, actionID)
}

func interfacesListURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "interfaces")
}

func securityGroupsListURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "securitygroups")
}

func portsListURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "ports")
}

func addSecurityGroupsURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "addsecuritygroup")
}

func deleteSecurityGroupsURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "delsecuritygroup")
}

func attachInterfaceURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "attach_interface")
}

func detachInterfaceURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "detach_interface")
}

func startInstanceURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "start")
}

func stopInstanceURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "stop")
}

func powerCycleInstanceURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "powercycle")
}

func rebootInstanceURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "reboot")
}

func suspendInstanceURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "suspend")
}

func resumeInstanceURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "resume")
}

func changeFlavorInstanceURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "changeflavor")
}

func metadataURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "metadata")
}

func metadataDetailsURL(c *edgecloud.ServiceClient, id string, actionID string) string {
	return resourceActionDetailsURL(c, id, "metadata", actionID)
}

func listAvailableFlavorsURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "available_flavors")
}

func listInstanceMetricsURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "metrics")
}

func createURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func getSpiceConsoleURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "get_spice_console")
}

func getInstanceConsoleURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "get_console")
}

func listInstanceLocationURL(c *edgecloud.ServiceClient) string {
	return c.BaseServiceURL("instances", "search")
}
