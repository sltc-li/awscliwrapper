package cmds

import (
	"fmt"
	"os"

	"github.com/tcnksm/go-input"
	"github.com/urfave/cli"

	"github.com/sltc-li/awscliwrapper"
)

var (
	InputUI = &input.UI{Writer: os.Stdout, Reader: os.Stdin}
)

func Commands() cli.Commands {
	return cli.Commands{
		{
			Name:        "eb",
			Usage:       "EB commands",
			Subcommands: EBCommands(),
		},
		{
			Name:        "s3",
			Usage:       "S3 commands",
			Subcommands: S3Commands(),
		},
		{
			Name:        "ecs",
			Usage:       "ECS commands",
			Subcommands: ECSCommands(),
		},
		{
			Name:  "whoami",
			Usage: "Show information of the current AWS user",
			Action: ActionFunc(func(w *awscliwrapper.Wrapper) error {
				u, err := w.IAM.GetCurrentUser()
				if err != nil {
					return err
				}

				fmt.Printf("Current AWS user: %s (%s)\n", u.Name, u.ARN)
				return nil
			}),
		},
	}
}
