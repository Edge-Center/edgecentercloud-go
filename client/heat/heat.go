package heat

import (
	"github.com/urfave/cli/v2"

	"github.com/Edge-Center/edgecentercloud-go/client/heat/v1/resources"
	"github.com/Edge-Center/edgecentercloud-go/client/heat/v1/stacks"
)

var Commands = cli.Command{
	Name:  "heat",
	Usage: "EdgeCloud Heat API",
	Subcommands: []*cli.Command{
		&resources.ResourceCommands,
		&stacks.StackCommands,
	},
}
