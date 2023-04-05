package keypairs

import (
	"github.com/urfave/cli/v2"

	"github.com/Edge-Center/edgecentercloud-go/client/flags"
	"github.com/Edge-Center/edgecentercloud-go/client/keypairs/v1/client"
	"github.com/Edge-Center/edgecentercloud-go/client/utils"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/keypair/v1/keypairs"
)

var keyPairIDText = "keypair_id is mandatory argument"

var keypairListCommand = cli.Command{
	Name:     "list",
	Usage:    "List keypairs",
	Category: "keypair",
	Action: func(c *cli.Context) error {
		client, err := client.NewKeypairClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		pages, err := keypairs.List(client).AllPages()
		if err != nil {
			return cli.Exit(err, 1)
		}
		results, err := keypairs.ExtractKeyPairs(pages)
		if err != nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(results, c.String("format"))

		return nil
	},
}

var keypairGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get keypair information",
	ArgsUsage: "<keypair_id>",
	Category:  "keypair",
	Action: func(c *cli.Context) error {
		keypairID, err := flags.GetFirstStringArg(c, keyPairIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewKeypairClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		task, err := keypairs.Get(client, keypairID).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(task, c.String("format"))

		return nil
	},
}

var keypairDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete keypair by ID",
	ArgsUsage: "<keypair_id>",
	Category:  "keypair",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		keypairID, err := flags.GetFirstStringArg(c, keyPairIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewKeypairClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		err = keypairs.Delete(client, keypairID).ExtractErr()
		if err != nil {
			return cli.Exit(err, 1)
		}

		return nil
	},
}

var keypairCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create keypair",
	Category: "keypair",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Keypair name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "ssh-public-key",
			Usage:    "Keypair SSH public key file",
			Aliases:  []string{"k"},
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewKeypairClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		data, err := utils.ReadFile(c.String("ssh-public-key"))
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.Exit(err, 1)
		}
		opts := keypairs.CreateOpts{
			Name:      c.String("name"),
			PublicKey: string(data),
		}
		result, err := keypairs.Create(client, opts).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(result, c.String("format"))

		return nil
	},
}

var KeypairCommands = cli.Command{
	Name:  "keypair",
	Usage: "EdgeCloud keypairs API",
	Subcommands: []*cli.Command{
		&keypairListCommand,
		&keypairGetCommand,
		&keypairDeleteCommand,
		&keypairCreateCommand,
	},
}
