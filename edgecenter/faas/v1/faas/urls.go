package faas

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func rootURL(c *edgecloud.ServiceClient) string {
	return c.ServiceURL("")
}

func namespaceListURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func namespaceCreateURL(c *edgecloud.ServiceClient) string {
	return rootURL(c)
}

func namespaceURL(c *edgecloud.ServiceClient, namespaceName string) string {
	return c.ServiceURL(namespaceName)
}

func functionListURL(c *edgecloud.ServiceClient, namespaceName string) string {
	return c.ServiceURL(namespaceName, "functions")
}

func functionCreateURL(c *edgecloud.ServiceClient, namespaceName string) string {
	return c.ServiceURL(namespaceName, "functions")
}

func functionURL(c *edgecloud.ServiceClient, namespaceName, functionName string) string {
	return c.ServiceURL(namespaceName, "functions", functionName)
}
