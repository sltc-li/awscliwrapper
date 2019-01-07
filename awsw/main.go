package main

import (
	"log"
	"os"

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
	}
	app.Commands = cmds.Commands()
	app.EnableBashCompletion = true

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
