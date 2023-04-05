package keystones

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/flags"
	"github.com/Edge-Center/edgecentercloud-go/client/keystones/v1/client"
	"github.com/Edge-Center/edgecentercloud-go/client/utils"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/keystone/v1/keystones"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/keystone/v1/types"
)

var (
	keystoneIDText     = "keystone_id is mandatory argument"
	keystoneStatesList = types.KeystoneState("").StringList()
)

var keystoneListCommand = cli.Command{
	Name:     "list",
	Usage:    "List keystones",
	Category: "keystone",
	Action: func(c *cli.Context) error {
		client, err := client.NewKeystoneClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		results, err := keystones.ListAll(client)
		if err != nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(results, c.String("format"))

		return nil
	},
}

var keystoneGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get keystone",
	ArgsUsage: "<keystone_id>",
	Category:  "keystone",
	Action: func(c *cli.Context) error {
		keystoneID, err := flags.GetFirstIntArg(c, keystoneIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewKeystoneClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		task, err := keystones.Get(client, keystoneID).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(task, c.String("format"))

		return nil
	},
}

var keystoneUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update keystone",
	ArgsUsage: "<keystone_id>",
	Category:  "keystone",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "url",
			Usage:    "keystone API url",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "domain-id",
			Usage:    "keystone federated domain ID",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "state",
			Aliases: []string{"s"},
			Value: &utils.EnumValue{
				Enum: keystoneStatesList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(keystoneStatesList, ", ")),
			Required: false,
		},
		&cli.StringFlag{
			Name:     "password",
			Usage:    "keystone password",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		keystoneID, err := flags.GetFirstIntArg(c, keystoneIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}

		url, err := edgecloud.ParseURLNonMandatory(c.String("spice-url"))
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.Exit(err, 1)
		}

		opts := keystones.UpdateOpts{
			URL:                       url,
			State:                     types.KeystoneState(c.String("state")),
			KeystoneFederatedDomainID: c.String("domain-id"),
			AdminPassword:             c.String("password"),
		}

		err = edgecloud.TranslateValidationError(opts.Validate())
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.Exit(err, 1)
		}

		client, err := client.NewKeystoneClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}

		result, err := keystones.Update(client, keystoneID, opts).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(result, c.String("format"))

		return nil
	},
}

var keystoneCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create keystone",
	Category: "keystone",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "url",
			Usage:    "keystone API url",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "domain-id",
			Usage:    "keystone federated domain ID",
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "state",
			Aliases: []string{"s"},
			Value: &utils.EnumValue{
				Enum: keystoneStatesList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(keystoneStatesList, ", ")),
			Required: true,
		},
		&cli.StringFlag{
			Name:     "password",
			Usage:    "keystone password",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		url, err := edgecloud.ParseURL(c.String("spice-url"))
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.Exit(err, 1)
		}

		opts := keystones.CreateOpts{
			URL:                       *url,
			State:                     types.KeystoneState(c.String("state")),
			KeystoneFederatedDomainID: c.String("domain-id"),
			AdminPassword:             c.String("password"),
		}

		err = edgecloud.TranslateValidationError(opts.Validate())
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.Exit(err, 1)
		}

		client, err := client.NewKeystoneClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}

		result, err := keystones.Create(client, opts).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(result, c.String("format"))

		return nil
	},
}

var Commands = cli.Command{
	Name:  "keystone",
	Usage: "EdgeCloud keystones API",
	Subcommands: []*cli.Command{
		&keystoneListCommand,
		&keystoneGetCommand,
		&keystoneUpdateCommand,
		&keystoneCreateCommand,
	},
}
