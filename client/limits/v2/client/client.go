package client

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/common"

	"github.com/urfave/cli/v2"
)

func NewLimitClientV2(c *cli.Context) (*edgecloud.ServiceClient, error) {
	return common.BuildClient(c, "limits_request", "v2")
}
