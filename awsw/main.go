package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/urfave/cli"

	"github.com/li-go/awscliwrapper/awsw/cmds"
)

func main() {
	app := cli.NewApp()
	app.Name = "awsw"
	app.Usage = "a simple wrapper command for awscli"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "region",
			Usage: "aws region",
			Value: "ap-northeast-1",
		},
		cli.StringFlag{
			Name:  "profile",
			Usage: "aws profile",
			Value: "default",
		},
		cli.BoolFlag{
			Name:  "fish",
			Usage: "generate fish completion",
		},
	}
	app.Commands = cmds.Commands()
	app.EnableBashCompletion = true
	app.Action = func(c *cli.Context) error {
		if c.Bool("fish") {
			completion, err := c.App.ToFishCompletion()
			if err != nil {
				return err
			}
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			completionFile := path.Join(home, ".config", "fish", "completions", "awsw.fish")
			fmt.Printf("Installing to %s\n", completionFile)
			if err := ioutil.WriteFile(completionFile, []byte(completion), 0644); err != nil {
				return err
			}
			fmt.Println("Done!")
			return nil
		}
		return cli.ShowAppHelp(c)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
