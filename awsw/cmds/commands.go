package cmds

import (
	"os"

	"github.com/tcnksm/go-input"
	"github.com/urfave/cli"

	"github.com/li-go/awscliwrapper"
)

var (
	ui = &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}
)

func Commands() cli.Commands {
	return cli.Commands{
		{
			Name:  "eb-desc",
			Usage: "describe elasticbeanstalk",
			Action: func(c *cli.Context) error {
				wrapper, err := newEBWrapper(c)
				if err != nil {
					return err
				}
				ec2wrapper, err := newEC2Wrapper(c)
				if err != nil {
					return err
				}
				return describeEB(wrapper, ec2wrapper)
			},
		},
		{
			Name:  "eb-deploy",
			Usage: "deploy elasticbeanstalk",
			Action: func(c *cli.Context) error {
				wrapper, err := newEBWrapper(c)
				if err != nil {
					return err
				}
				return deployEB(wrapper)
			},
		},
		{
			Name:  "s3-ls",
			Usage: "",
			Action: func(c *cli.Context) error {
				wrapper, err := newS3Wrapper(c)
				if err != nil {
					return err
				}
				return listS3(wrapper)
			},
		},
	}
}

func newEBWrapper(c *cli.Context) (*awscliwrapper.EBWrapper, error) {
	region, profile := c.GlobalString("region"), c.GlobalString("profile")
	return awscliwrapper.NewEB(region, profile)
}

func newEC2Wrapper(c *cli.Context) (*awscliwrapper.EC2Wrapper, error) {
	region, profile := c.GlobalString("region"), c.GlobalString("profile")
	return awscliwrapper.NewEC2(region, profile)
}

func newS3Wrapper(c *cli.Context) (*awscliwrapper.S3Wrapper, error) {
	region, profile := c.GlobalString("region"), c.GlobalString("profile")
	return awscliwrapper.NewS3(region, profile)
}
