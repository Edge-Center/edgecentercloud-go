package ports

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func resourceActionURL(c *edgecloud.ServiceClient, id string, action string) string {
	return c.ServiceURL(id, action)
}

func enablePortSecurityURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "enable_port_security")
}

func disablePortSecurityURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "disable_port_security")
}

func assignAllowedAddressPairsURL(c *edgecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "allow_address_pairs")
}
