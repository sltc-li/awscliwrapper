package cmds

import (
	"github.com/urfave/cli"
)

func IAMCommands() cli.Commands {
	return cli.Commands{
		{
			Name: "whoami",
		},
	}
}
