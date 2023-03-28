package client

import (
	"github.com/urfave/cli/v2"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/common"
)

func NewRegionClientV1(c *cli.Context) (*edgecloud.ServiceClient, error) {
	return common.BuildClient(c, "reseller_region", "v1")
}
