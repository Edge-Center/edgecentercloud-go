package client

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/common"

	"github.com/urfave/cli/v2"
)

func NewLoadbalancerClientV1(c *cli.Context) (*edgecloud.ServiceClient, error) {
	return common.BuildClient(c, "loadbalancers", "v1")
}

func NewLBListenerClientV1(c *cli.Context) (*edgecloud.ServiceClient, error) {
	return common.BuildClient(c, "lblisteners", "v1")
}

func NewLBPoolClientV1(c *cli.Context) (*edgecloud.ServiceClient, error) {
	return common.BuildClient(c, "lbpools", "v1")
}

func NewLBFlavorClientV1(c *cli.Context) (*edgecloud.ServiceClient, error) {
	return common.BuildClient(c, "lbflavors", "v1")
}
