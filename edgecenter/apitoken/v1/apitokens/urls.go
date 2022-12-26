package apitokens

import (
	"strconv"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func resourceURL(c *edgecloud.ServiceClient, clientID, tokenID int) string {
	return c.ServiceURL("clients", strconv.Itoa(clientID), "tokens", strconv.Itoa(tokenID))
}

func rootURL(c *edgecloud.ServiceClient, clientID int) string {
	return c.ServiceURL("clients", strconv.Itoa(clientID), "tokens")
}

func getURL(c *edgecloud.ServiceClient, clientID, tokenID int) string {
	return resourceURL(c, clientID, tokenID)
}

func listURL(c *edgecloud.ServiceClient, clientID int) string {
	return rootURL(c, clientID)
}

func createURL(c *edgecloud.ServiceClient, clientID int) string {
	return rootURL(c, clientID)
}

func deleteURL(c *edgecloud.ServiceClient, clientID, tokenID int) string {
	return resourceURL(c, clientID, tokenID)
}
