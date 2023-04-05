package securitygrouprules

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/Edge-Center/edgecentercloud-go/client/flags"
	"github.com/Edge-Center/edgecentercloud-go/client/securitygroups/v1/client"
	"github.com/Edge-Center/edgecentercloud-go/client/utils"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/securitygroup/v1/securitygrouprules"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/securitygroup/v1/securitygroups"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/securitygroup/v1/types"
)

var (
	securityGroupRuleIDText = "securitygrouprule_id is mandatory argument"
	securityGroupIDText     = "securitygroup_id is mandatory argument"
	protocolTypeList        = types.Protocol("").StringList()
	directionTypeList       = types.RuleDirection("").StringList()
	etherTypeTypeList       = types.EtherType("").StringList()
)

var securityGroupRuleDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete security group rule",
	ArgsUsage: "<securitygrouprule_id>",
	Category:  "securitygrouprule",
	Action: func(c *cli.Context) error {
		securityGroupRuleID, err := flags.GetFirstStringArg(c, securityGroupRuleIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewSecurityGroupRuleClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		err = securitygrouprules.Delete(client, securityGroupRuleID).ExtractErr()
		if err != nil {
			return cli.Exit(err, 1)
		}

		return nil
	},
}

var securityGroupRuleUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "Update security group rule",
	ArgsUsage: "<securitygrouprule_id>",
	Category:  "securitygrouprule",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "security-group-id",
			Usage:    "Security group ID",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "description",
			Usage:    "Security group rule description",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "remote-group-id",
			Usage:    "Security group rule remote group id",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "remote-ip-prefix",
			Usage:    "Security group rule remote ip prefix",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "port-range-max",
			Usage:    "Security group rule port max range",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "port-range-min",
			Usage:    "Security group rule port min range",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "protocol",
			Aliases: []string{"p"},
			Value: &utils.EnumValue{
				Enum: protocolTypeList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(protocolTypeList, ", ")),
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "ethertype",
			Aliases: []string{"e"},
			Value: &utils.EnumValue{
				Enum: etherTypeTypeList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(etherTypeTypeList, ", ")),
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "direction",
			Aliases: []string{"dr"},
			Value: &utils.EnumValue{
				Enum: directionTypeList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(directionTypeList, ", ")),
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		securityGroupRuleID, err := flags.GetFirstStringArg(c, securityGroupRuleIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "add-rule")
			return err
		}
		client, err := client.NewSecurityGroupRuleClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}

		opts := securitygroups.CreateSecurityGroupRuleOpts{
			SecurityGroupID: utils.StringToPointer(c.String("security-group-id")),
			Direction:       types.RuleDirection(c.String("direction")),
			RemoteGroupID:   utils.StringToPointer(c.String("remote-group-id")),
			EtherType:       types.EtherType(c.String("ethertype")),
			Protocol:        types.Protocol(c.String("protocol")),
			PortRangeMax:    utils.IntToPointer(c.Int("port-range-max")),
			PortRangeMin:    utils.IntToPointer(c.Int("port-range-min")),
			Description:     utils.StringToPointer(c.String("description")),
			RemoteIPPrefix:  utils.StringToPointer(c.String("remote-ip-prefix")),
		}

		results, err := securitygrouprules.Replace(client, securityGroupRuleID, opts).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(results, c.String("format"))

		return nil
	},
}

var securityGroupRuleAddSubCommand = cli.Command{
	Name:      "add",
	Usage:     "Add security group rule",
	ArgsUsage: "<securitygroup_id>",
	Category:  "securitygrouprule",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "description",
			Usage:    "Security group rule description",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "remote-group-id",
			Usage:    "Security group rule remote group id",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "remote-ip-prefix",
			Usage:    "Security group rule remote ip prefix",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "port-range-max",
			Usage:    "Security group rule port max range",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "port-range-min",
			Usage:    "Security group rule port min range",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "protocol",
			Aliases: []string{"p"},
			Value: &utils.EnumValue{
				Enum: protocolTypeList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(protocolTypeList, ", ")),
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "ethertype",
			Aliases: []string{"e"},
			Value: &utils.EnumValue{
				Enum: etherTypeTypeList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(etherTypeTypeList, ", ")),
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "direction",
			Aliases: []string{"dr"},
			Value: &utils.EnumValue{
				Enum: directionTypeList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(directionTypeList, ", ")),
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		securityGroupID, err := flags.GetFirstStringArg(c, securityGroupIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "add-rule")
			return err
		}
		client, err := client.NewSecurityGroupClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}

		opts := securitygroups.CreateSecurityGroupRuleOpts{
			Direction:      types.RuleDirection(c.String("direction")),
			RemoteGroupID:  utils.StringToPointer(c.String("remote-group-id")),
			EtherType:      types.EtherType(c.String("ethertype")),
			Protocol:       types.Protocol(c.String("protocol")),
			PortRangeMax:   utils.IntToPointer(c.Int("port-range-max")),
			PortRangeMin:   utils.IntToPointer(c.Int("port-range-min")),
			Description:    utils.StringToPointer(c.String("description")),
			RemoteIPPrefix: utils.StringToPointer(c.String("remote-ip-prefix")),
		}

		results, err := securitygroups.AddRule(client, securityGroupID, opts).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(results, c.String("format"))

		return nil
	},
}

var SecurityGroupRuleCommands = []*cli.Command{
	&securityGroupRuleUpdateSubCommand,
	&securityGroupRuleDeleteSubCommand,
	&securityGroupRuleAddSubCommand,
}
