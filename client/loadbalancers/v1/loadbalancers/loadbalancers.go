package loadbalancers

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/flags"
	"github.com/Edge-Center/edgecentercloud-go/client/loadbalancers/v1/client"
	"github.com/Edge-Center/edgecentercloud-go/client/loadbalancers/v1/lbpools"
	"github.com/Edge-Center/edgecentercloud-go/client/loadbalancers/v1/listeners"
	"github.com/Edge-Center/edgecentercloud-go/client/utils"
	cmeta "github.com/Edge-Center/edgecentercloud-go/client/utils/metadata"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/loadbalancer/v1/lbflavors"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/loadbalancer/v1/loadbalancers"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/task/v1/tasks"
)

var loadBalancerIDText = "loadbalancer_id is mandatory argument"

var loadBalancerListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "loadbalancers list",
	Category: "loadbalancer",
	Action: func(c *cli.Context) error {
		client, err := client.NewLoadbalancerClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		results, err := loadbalancers.ListAll(client, nil)
		if err != nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(results, c.String("format"))

		return nil
	},
}

var loadBalancerCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "create loadbalancer",
	Category: "loadbalancer",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Loadbalancer name",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "flavor",
			Aliases: []string{"fl"},
			Usage:   "Loadbalancer flavor",
		},
		&cli.StringFlag{
			Name:        "vip-port-id",
			Usage:       "Loadbalancer vip port ID",
			DefaultText: "<nil>",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "vip-network-id",
			Usage:       "Loadbalancer vip network ID",
			DefaultText: "<nil>",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "vip-subnet-id",
			Usage:       "Loadbalancer vip subnet ID",
			DefaultText: "<nil>",
			Required:    false,
		},
		&cli.StringSliceFlag{
			Name:    "tags",
			Aliases: []string{"t"},
			Usage:   "Loadbalancer tags",
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		client, err := client.NewLoadbalancerClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}

		opts := loadbalancers.CreateOpts{
			Name:         c.String("name"),
			Listeners:    []loadbalancers.CreateListenerOpts{},
			VipPortID:    c.String("vip-port-id"),
			VipNetworkID: c.String("vip-network-id"),
			VipSubnetID:  c.String("vip-subnet-id"),
			Tags:         c.StringSlice("tags"),
		}
		flavor := c.String("flavor")
		if flavor != "" {
			opts.Flavor = &flavor
		}

		results, err := loadbalancers.Create(client, opts).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			loadBalancerID, err := loadbalancers.ExtractLoadBalancerIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve loadbalancer ID from task info: %w", err)
			}
			loadBalancer, err := loadbalancers.Get(client, loadBalancerID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get loadbalancer with ID: %s. Error: %w", loadBalancerID, err)
			}
			utils.ShowResults(loadBalancer, c.String("format"))

			return nil, nil
		})
	},
}

var loadBalancerGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "show loadbalancer",
	ArgsUsage: "<loadbalancer_id>",
	Category:  "loadbalancer",
	Action: func(c *cli.Context) error {
		loadBalancerID, err := flags.GetFirstStringArg(c, loadBalancerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewLoadbalancerClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		result, err := loadbalancers.Get(client, loadBalancerID).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(result, c.String("format"))

		return nil
	},
}

var loadBalancerDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "delete loadbalancer",
	ArgsUsage: "<loadbalancer_id>",
	Category:  "loadbalancer",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		loadBalancerID, err := flags.GetFirstStringArg(c, loadBalancerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewLoadbalancerClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		results, err := loadbalancers.Delete(client, loadBalancerID).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			loadbalancer, err := loadbalancers.Get(client, loadBalancerID).Extract()
			if err == nil {
				if loadbalancer != nil && loadbalancer.IsDeleted() {
					return nil, nil
				}
				return nil, fmt.Errorf("cannot delete loadbalancer with ID: %s", loadBalancerID)
			}
			var e edgecloud.Default404Error
			if errors.As(err, &e) {
				return nil, nil
			}

			return nil, err
		})
	},
}

var loadBalancerUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "update loadbalancer",
	ArgsUsage: "<loadbalancer_id>",
	Category:  "loadbalancer",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Loadbalancer name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		loadBalancerID, err := flags.GetFirstStringArg(c, loadBalancerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}
		client, err := client.NewLoadbalancerClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}

		opts := loadbalancers.UpdateOpts{Name: c.String("name")}

		result, err := loadbalancers.Update(client, loadBalancerID, opts).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		if result == nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(result, c.String("format"))

		return nil
	},
}

var flavorListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "List loadbalancer flavor",
	Category: "loadbalancer flavor",
	Action: func(c *cli.Context) error {
		client, err := client.NewLBFlavorClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		results, err := lbflavors.ListAll(client)
		if err != nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(results, c.String("format"))

		return nil
	},
}

var flavorSubCommand = cli.Command{
	Name:  "flavor",
	Usage: "EdgeCloud loadbalancer flavor API",
	Subcommands: []*cli.Command{
		&flavorListSubCommand,
	},
}

var Commands = cli.Command{
	Name:  "loadbalancer",
	Usage: "EdgeCloud loadbalancers API",
	Subcommands: []*cli.Command{
		&loadBalancerListSubCommand,
		&loadBalancerGetSubCommand,
		&loadBalancerUpdateSubCommand,
		&loadBalancerDeleteSubCommand,
		&loadBalancerCreateSubCommand,
		&flavorSubCommand,
		&listeners.ListenerCommands,
		&lbpools.PoolCommands,
		{
			Name:  "metadata",
			Usage: "Loadbalancer metadata",
			Subcommands: []*cli.Command{
				cmeta.NewMetadataListCommand(
					client.NewLoadbalancerClientV1,
					"Get loadbalancer metadata",
					"<loadbalancer_id>",
					"loadbalancer_id is mandatory argument",
				),
				cmeta.NewMetadataGetCommand(
					client.NewLoadbalancerClientV1,
					"Show loadbalancer metadata by key",
					"<loadbalancer_id>",
					"loadbalancer_id is mandatory argument",
				),
				cmeta.NewMetadataDeleteCommand(
					client.NewLoadbalancerClientV1,
					"Delete loadbalancer metadata by key",
					"<loadbalancer_id>",
					"loadbalancer_id is mandatory argument",
				),
				cmeta.NewMetadataCreateCommand(
					client.NewLoadbalancerClientV1,
					"Create loadbalancer metadata. It would update existing keys",
					"<loadbalancer_id>",
					"loadbalancer_id is mandatory argument",
				),
				cmeta.NewMetadataUpdateCommand(
					client.NewLoadbalancerClientV1,
					"Update loadbalancer metadata. It overriding existing records",
					"<loadbalancer_id>",
					"loadbalancer_id is mandatory argument",
				),
				cmeta.NewMetadataReplaceCommand(
					client.NewLoadbalancerClientV1,
					"Replace loadbalancer metadata. It replace existing records",
					"<loadbalancer_id>",
					"loadbalancer_id is mandatory argument",
				),
			},
		},
	},
}
