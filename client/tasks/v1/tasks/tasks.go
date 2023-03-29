package tasks

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/Edge-Center/edgecentercloud-go/client/tasks/v1/client"
	"github.com/Edge-Center/edgecentercloud-go/client/utils"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/task/v1/tasks"
)

var taskIDText = "task_id is mandatory argument"

var taskListCommand = cli.Command{
	Name:     "list",
	Usage:    "List active tasks",
	Category: "task",
	Action: func(c *cli.Context) error {
		client, err := client.NewTaskClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		pages, err := tasks.List(client).AllPages()
		if err != nil {
			return cli.Exit(err, 1)
		}
		results, err := tasks.ExtractTasks(pages)
		if len(results) == 0 {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(results, c.String("format"))

		return nil
	},
}

var taskGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get task information",
	ArgsUsage: "<task_id>",
	Category:  "task",
	Action: func(c *cli.Context) error {
		taskID := c.Args().First()
		if taskID == "" {
			_ = cli.ShowCommandHelp(c, "show")
			return cli.Exit(fmt.Errorf(taskIDText), 1)
		}
		client, err := client.NewTaskClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		task, err := tasks.Get(client, taskID).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		if task == nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(task, c.String("format"))

		return nil
	},
}

var Commands = cli.Command{
	Name:  "task",
	Usage: "EdgeCloud tasks API",
	Subcommands: []*cli.Command{
		&taskListCommand,
		&taskGetCommand,
	},
}
