package client

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/common"
	"github.com/urfave/cli/v2"
)

func NewQuotaClientV2(c *cli.Context) (*edgecloud.ServiceClient, error) {
	return common.BuildClient(c, "quotas", "v2")
}
