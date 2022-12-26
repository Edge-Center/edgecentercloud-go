package resources

import (
	"fmt"
	"strings"

	"github.com/Edge-Center/edgecentercloud-go/client/flags"
	"github.com/Edge-Center/edgecentercloud-go/client/heat/v1/client"
	"github.com/Edge-Center/edgecentercloud-go/client/utils"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/heat/v1/stack/resources"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/heat/v1/stack/resources/types"

	"github.com/urfave/cli/v2"
)

var (
	resourceNameText      = "resource_id is mandatory argument"
	stackIDText           = "stack_id is mandatory argument"
	stackResourceActions  = types.StackResourceAction("").StringList()
	stackResourceStatuses = types.StackResourceStatus("").StringList()
)

var resourceMetadataSubCommand = cli.Command{
	Name:      "metadata",
	Usage:     "Get stack resource metadata",
	ArgsUsage: "<resource_name>",
	Category:  "resource",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "stack",
			Aliases:  []string{"s"},
			Usage:    "Stack ID",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		resourceName, err := flags.GetFirstStringArg(c, resourceNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "metadata")
			return err
		}

		client, err := client.NewHeatClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		metadata, err := resources.Metadata(client, c.String("stack"), resourceName).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(metadata, c.String("format"))
		return nil
	},
}

var resourceSignalSubCommand = cli.Command{
	Name:      "signal",
	Usage:     "Send stack resource signal",
	ArgsUsage: "<resource_name>",
	Category:  "resource",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "stack",
			Aliases:  []string{"s"},
			Usage:    "Stack ID",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "signal",
			Usage:    "Signal data",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		resourceName, err := flags.GetFirstStringArg(c, resourceNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "metadata")
			return err
		}

		client, err := client.NewHeatClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		data := c.String("signal")
		err = resources.Signal(client, c.String("stack"), resourceName, []byte(data)).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var resourceGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Stack resource",
	ArgsUsage: "<resource_name>",
	Category:  "resource",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "stack",
			Aliases:  []string{"s"},
			Usage:    "Stack ID",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		resourceName, err := flags.GetFirstStringArg(c, resourceNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}

		client, err := client.NewHeatClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := resources.Get(client, c.String("stack"), resourceName).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var resourceMarkUnhealthySubCommand = cli.Command{
	Name:      "unhealthy",
	Usage:     "Stack resource mark unhealthy",
	ArgsUsage: "<resource_name>",
	Category:  "resource",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "stack",
			Aliases:  []string{"s"},
			Usage:    "Stack ID",
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "mark-unhealthy",
			Aliases:  []string{"m"},
			Usage:    "Either mark stack resource as unhealthy or not",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "resource-status-reason",
			Aliases:  []string{"r"},
			Usage:    "Stack resource change status reason",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		resourceName, err := flags.GetFirstStringArg(c, resourceNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "unhealthy")
			return err
		}

		client, err := client.NewHeatClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := resources.MarkUnhealthyOpts{
			MarkUnhealthy:        c.Bool("mark-unhealthy"),
			ResourceStatusReason: c.String("resource-status-reason"),
		}

		err = resources.MarkUnhealthy(client, c.String("stack"), resourceName, opts).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var resourceListSubCommand = cli.Command{
	Name:      "list",
	Usage:     "Stack resources",
	ArgsUsage: "<stack_id>",
	Category:  "resource",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "type",
			Aliases:  []string{"t"},
			Usage:    "Stack resource type",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Stack resource name",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "physical-resource-id",
			Usage:    "Stack physical resource id",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "logical-resource-id",
			Usage:    "Stack logical resource id",
			Required: false,
		},
		&cli.GenericFlag{
			Name: "status",
			Value: &utils.EnumValue{
				Enum: stackResourceStatuses,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(stackResourceStatuses, ", ")),
			Required: false,
		},
		&cli.GenericFlag{
			Name: "action",
			Value: &utils.EnumValue{
				Enum: stackResourceActions,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(stackResourceActions, ", ")),
			Required: false,
		},
		&cli.IntFlag{
			Name:     "nested-depth",
			Usage:    "includes resources from nested stacks up to the nested-depth",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "with-detail",
			Usage:    "enables detailed resource information",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		stackID, err := flags.GetFirstStringArg(c, stackIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "list")
			return err
		}

		client, err := client.NewHeatClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := resources.ListOpts{
			Type:               c.String("type"),
			Name:               c.String("name"),
			Status:             types.StackResourceStatus(c.String("status")),
			Action:             types.StackResourceAction(c.String("action")),
			LogicalResourceID:  c.String("logical-resource-id"),
			PhysicalResourceID: c.String("physical-resource-id"),
			NestedDepth:        c.Int("nested-depth"),
			WithDetail:         c.Bool("with-detail"),
		}

		result, err := resources.ListAll(client, stackID, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var ResourceCommands = cli.Command{
	Name:  "resource",
	Usage: "Heat stack resource commands",
	Subcommands: []*cli.Command{
		&resourceMetadataSubCommand,
		&resourceSignalSubCommand,
		&resourceGetSubCommand,
		&resourceListSubCommand,
		&resourceMarkUnhealthySubCommand,
	},
}
