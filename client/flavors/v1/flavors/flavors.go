package flavors

import (
	"github.com/urfave/cli/v2"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/flavors/v1/client"
	"github.com/Edge-Center/edgecentercloud-go/client/utils"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/flavor/v1/flavors"
)

var flavorListCommand = cli.Command{
	Name:     "list",
	Usage:    "List flavors",
	Category: "flavor",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "include_prices",
			Aliases: []string{"p"},
			Usage:   "Include prices",
		},
		&cli.BoolFlag{
			Name:    "baremetal",
			Aliases: []string{"bm"},
			Usage:   "show only baremetal flavors",
		},
	},
	Action: func(c *cli.Context) error {
		var err error
		var cl *edgecloud.ServiceClient
		cl, err = client.NewFlavorClientV1(c)

		if c.Bool("baremetal") {
			cl, err = client.NewBmFlavorClientV1(c)
		}
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		prices := c.Bool("include_prices")
		opts := flavors.ListOpts{
			IncludePrices: &prices,
		}
		results, err := flavors.ListAll(cl, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))

		return nil
	},
}

var Commands = cli.Command{
	Name:  "flavor",
	Usage: "EdgeCloud flavors API",
	Subcommands: []*cli.Command{
		&flavorListCommand,
	},
}
