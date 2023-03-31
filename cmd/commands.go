package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/Edge-Center/edgecentercloud-go/client/apitokens/v1/apitokens"
	"github.com/Edge-Center/edgecentercloud-go/client/apptemplates/v1/apptemplates"
	"github.com/Edge-Center/edgecentercloud-go/client/flags"
	"github.com/Edge-Center/edgecentercloud-go/client/flavors/v1/flavors"
	"github.com/Edge-Center/edgecentercloud-go/client/floatingips/v1/floatingips"
	"github.com/Edge-Center/edgecentercloud-go/client/heat"
	"github.com/Edge-Center/edgecentercloud-go/client/images/v1/images"
	"github.com/Edge-Center/edgecentercloud-go/client/instances/v1/instances"
	"github.com/Edge-Center/edgecentercloud-go/client/k8s/v1/k8s"
	"github.com/Edge-Center/edgecentercloud-go/client/keypairs/v2/keypairs"
	"github.com/Edge-Center/edgecentercloud-go/client/keystones/v1/keystones"
	"github.com/Edge-Center/edgecentercloud-go/client/l7policies/v1/l7policies"
	"github.com/Edge-Center/edgecentercloud-go/client/lifecyclepolicy/v1/lifecyclepolicy"
	"github.com/Edge-Center/edgecentercloud-go/client/limits/v2/limits"
	"github.com/Edge-Center/edgecentercloud-go/client/loadbalancers/v1/loadbalancers"
	"github.com/Edge-Center/edgecentercloud-go/client/networks/v1/networks"
	"github.com/Edge-Center/edgecentercloud-go/client/ports/v1/ports"
	"github.com/Edge-Center/edgecentercloud-go/client/projects/v1/projects"
	"github.com/Edge-Center/edgecentercloud-go/client/quotas/v2/quotas"
	"github.com/Edge-Center/edgecentercloud-go/client/regions/v1/regions"
	"github.com/Edge-Center/edgecentercloud-go/client/regionsaccess/v1/regionsaccess"
	"github.com/Edge-Center/edgecentercloud-go/client/reservedfixedips/v1/reservedfixedips"
	"github.com/Edge-Center/edgecentercloud-go/client/routers/v1/routers"
	"github.com/Edge-Center/edgecentercloud-go/client/schedules/v1/schedules"
	"github.com/Edge-Center/edgecentercloud-go/client/secrets/v1/secrets"
	"github.com/Edge-Center/edgecentercloud-go/client/securitygroups/v1/securitygroups"
	"github.com/Edge-Center/edgecentercloud-go/client/servergroups/v1/servergroups"
	"github.com/Edge-Center/edgecentercloud-go/client/snapshots/v1/snapshots"
	"github.com/Edge-Center/edgecentercloud-go/client/subnets/v1/subnets"
	"github.com/Edge-Center/edgecentercloud-go/client/tasks/v1/tasks"
	"github.com/Edge-Center/edgecentercloud-go/client/volumes/v1/volumes"
)

var AppVersion = "v0.1.2"

var commands = []*cli.Command{
	&networks.Commands,
	&tasks.Commands,
	&keypairs.Commands,
	&volumes.Commands,
	&subnets.Commands,
	&flavors.Commands,
	&loadbalancers.Commands,
	&instances.Commands,
	&heat.Commands,
	&securitygroups.Commands,
	&floatingips.Commands,
	&schedules.Commands,
	&ports.Commands,
	&snapshots.Commands,
	&images.Commands,
	&regions.Commands,
	&projects.Commands,
	&keystones.Commands,
	&quotas.Commands,
	&limits.Commands,
	&k8s.Commands,
	&l7policies.Commands,
	&routers.Commands,
	&reservedfixedips.Commands,
	&servergroups.Commands,
	&secrets.Commands,
	&lifecyclepolicy.Commands,
	&regionsaccess.Commands,
	&apptemplates.Commands,
	&apitokens.Commands,
}

type clientCommands struct {
	commands []*cli.Command
	flags    []cli.Flag
	usage    string
}

func buildClientCommands(commands []*cli.Command) clientCommands {
	clientType := os.Getenv("EC_CLOUD_CLIENT_TYPE")
	tokenClientUsage := fmt.Sprintf("EdgeCloud API client\n%s", flags.TokenClientHelpText)
	platformClientUsage := fmt.Sprintf("EdgeCloud API client\n%s", flags.PlatformClientHelpText)
	apiTokenClientUsage := fmt.Sprintf("EdgeCloud API client\n%s", flags.APITokenClientHelpText)
	switch clientType {
	case flags.ClientTypeToken:
		flags.ClientType = clientType
		return clientCommands{
			commands: commands,
			flags:    flags.TokenClientFlags,
			usage:    tokenClientUsage,
		}
	case flags.ClientTypePlatform:
		flags.ClientType = clientType
		return clientCommands{
			commands: commands,
			flags:    flags.PlatformClientFlags,
			usage:    platformClientUsage,
		}
	case flags.ClientTypeAPIToken:
		flags.ClientType = clientType
		return clientCommands{
			commands: commands,
			flags:    flags.APITokenClientFlags,
			usage:    apiTokenClientUsage,
		}
	}
	mainClientUsage := fmt.Sprintf("EdgeCloud API client\n%s", flags.MainClientHelpText)

	return clientCommands{
		commands: []*cli.Command{
			{
				Name:        "token",
				Usage:       tokenClientUsage,
				Flags:       flags.TokenClientFlags,
				Subcommands: commands,
				Before: func(c *cli.Context) error {
					return c.Set("client-type", "token")
				},
			},
			{
				Name:        "platform",
				Usage:       platformClientUsage,
				Flags:       flags.PlatformClientFlags,
				Subcommands: commands,
				Before: func(c *cli.Context) error {
					return c.Set("client-type", "platform")
				},
			},
			{
				Name:        "api-token",
				Usage:       apiTokenClientUsage,
				Flags:       flags.APITokenClientFlags,
				Subcommands: commands,
				Before: func(c *cli.Context) error {
					return c.Set("client-type", "api-token")
				},
			},
		},
		flags: nil,
		usage: mainClientUsage,
	}
}

func NewApp(args []string) *cli.App {
	flags.AddOutputFlags(commands)
	clientCommands := buildClientCommands(commands)

	app := new(cli.App)
	app.Name = filepath.Base(args[0])
	app.HelpName = filepath.Base(args[0])
	app.Version = AppVersion
	app.EnableBashCompletion = true
	app.Commands = clientCommands.commands
	if clientCommands.flags != nil {
		app.Flags = clientCommands.flags
	}
	if len(clientCommands.usage) > 0 {
		app.Usage = clientCommands.usage
	}

	return app
}

func RunCommand(args []string) {
	app := NewApp(args)
	err := app.Run(args)
	if err != nil {
		logrus.Errorf("Cannot initialize application: %+v", err)
		os.Exit(1)
	}
}
