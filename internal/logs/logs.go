package logs

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/executor"
)

func Command(e executor.Executor) *cli.Command {
	return &cli.Command{
		Name:  "logs",
		Usage: "Show logs for a machine",
		Action: func(c *cli.Context) error {
			return cli.ShowSubcommandHelp(c)
		},
		Subcommands: []*cli.Command{
			{
				Name:        "nginx",
				Description: "Show logs from nginx",
				Action: func(c *cli.Context) error {
					return logsNginx(c, e)
				},
			},
		},
	}
}

func logsNginx(c *cli.Context, e executor.Executor) error {
	machine := c.String("machine")

	args := []string{"multipass", "exec", machine, "--", "sudo", "bash", "/opt/nitro/nginx/tail-logs.sh"}

	return e.Exec(e.Path(), args, os.Environ())
}
