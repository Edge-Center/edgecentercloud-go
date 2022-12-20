package client

import (
	edgecenter "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/common"

	"github.com/urfave/cli/v2"
)

func NewInstanceClientV2(c *cli.Context) (*edgecenter.ServiceClient, error) {
	return common.BuildClient(c, "instances", "v2")
}
