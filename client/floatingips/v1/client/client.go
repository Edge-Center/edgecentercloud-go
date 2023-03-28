package client

import (
	"github.com/urfave/cli/v2"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/common"
)

func NewFloatingIPClientV1(c *cli.Context) (*edgecloud.ServiceClient, error) {
	return common.BuildClient(c, "floatingips", "v1")
}

func NewAvailableFloatingIPClientV1(c *cli.Context) (*edgecloud.ServiceClient, error) {
	return common.BuildClient(c, "availablefloatingips", "v1")
}
