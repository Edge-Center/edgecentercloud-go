package availablefloatingips

import (
	"github.com/urfave/cli/v2"

	"github.com/Edge-Center/edgecentercloud-go/client/floatingips/v1/client"
	"github.com/Edge-Center/edgecentercloud-go/client/utils"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/floatingip/v1/floatingips"
)

var availableFloatingIPListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "Available floating ips list",
	Category: "availablefloatingip",
	Action: func(c *cli.Context) error {
		client, err := client.NewAvailableFloatingIPClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := floatingips.ListAll(client, nil)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var AvailableFloatingIPCommands = cli.Command{
	Name:  "available",
	Usage: "EdgeCloud available floating ips API",
	Subcommands: []*cli.Command{
		&availableFloatingIPListSubCommand,
	},
}
