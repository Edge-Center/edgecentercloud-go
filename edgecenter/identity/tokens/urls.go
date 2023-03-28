package tokens

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

func tokenURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL("auth", "jwt", "login")
}

func refreshURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL("auth", "jwt", "refresh")
}

func refreshECCloudURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL("v1", "token", "refresh")
}

func selectAccountURL(c *edgecloud.ServiceClient, clientID string) string {
	return c.ServiceURL("auth", "jwt", "clients", clientID, "login")
}
