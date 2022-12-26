package heat

import (
	"github.com/Edge-Center/edgecentercloud-go/client/heat/v1/resources"
	"github.com/Edge-Center/edgecentercloud-go/client/heat/v1/stacks"
	"github.com/urfave/cli/v2"
)

var Commands = cli.Command{
	Name:  "heat",
	Usage: "EdgeCloud Heat API",
	Subcommands: []*cli.Command{
		&resources.ResourceCommands,
		&stacks.StackCommands,
	},
}
