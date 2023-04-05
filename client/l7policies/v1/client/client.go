package client

import (
	"github.com/urfave/cli/v2"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/common"
)

func NewL7PoliciesClientV1(c *cli.Context) (*edgecloud.ServiceClient, error) {
	return common.BuildClient(c, "l7policies", "v1")
}

func NewL7RulesClientV1(c *cli.Context) (*edgecloud.ServiceClient, error) {
	return NewL7PoliciesClientV1(c)
}
