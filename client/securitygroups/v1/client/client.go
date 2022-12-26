package client

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/common"

	"github.com/urfave/cli/v2"
)

func NewSecurityGroupRuleClientV1(c *cli.Context) (*edgecloud.ServiceClient, error) {
	return common.BuildClient(c, "securitygrouprules", "v1")
}

func NewSecurityGroupClientV1(c *cli.Context) (*edgecloud.ServiceClient, error) {
	return common.BuildClient(c, "securitygroups", "v1")
}
