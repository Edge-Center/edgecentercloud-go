/*
Package edgecenter contains resources for the individual projects.
It also includes functions to authenticate to an
EdgeCloud cloud and for provisioning various service-level clients.

Example of Creating a EdgeCenter Magnum Cluster Client

	ao, err := edgecenter.AuthOptionsFromEnv()
	provider, err := edgecenter.AuthenticatedClient(ao)
	client, err := edgecenter.MagnumClusterV1(client, edgecloud.EndpointOpts{
	})
*/
package edgecenter
