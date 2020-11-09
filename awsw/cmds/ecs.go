package cmds

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/fatih/color"
	"github.com/tcnksm/go-input"
	"github.com/urfave/cli"

	"github.com/li-go/awscliwrapper"
)

func ECSCommands() cli.Commands {
	return cli.Commands{
		{
			Name:  "walk",
			Usage: "walk ECS",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster",
					Usage: "cluster name",
				},
				cli.StringFlag{
					Name:  "service",
					Usage: "service name",
				},
				cli.StringFlag{
					Name:  "container",
					Usage: "container name",
				},
			},
			Action: ActionFuncWithContext(walkCluster),
		},
	}
}

func walkCluster(c *cli.Context, w *awscliwrapper.Wrapper) error {
	cluster := c.String("cluster")
	service := c.String("service")
	container := c.String("container")

	if cluster == "" {
		arn, err := getARN("select a cluster?", w.ECS.ListClusters)
		if err != nil {
			return err
		}
		cluster = arn.Name()
	}
	fmt.Printf("cluster: %s\n\n", color.GreenString(cluster))

	if service == "" {
		arn, err := getARN("select a service?", func() ([]awscliwrapper.ARN, error) {
			return w.ECS.GetServices(cluster)
		})
		if err != nil {
			return err
		}
		service = arn.Name()
	}
	fmt.Printf("service: %s\n\n", color.GreenString(service))

	taskDef, err := w.ECS.GetTaskDefinition(cluster, service)
	if err != nil {
		return err
	}
	fmt.Printf("task definition: %s\n\n", color.GreenString(taskDef.Name()))

	containerDefs, err := w.ECS.GetContainerDefinitions(taskDef)
	if err != nil {
		return err
	}

	for _, d := range containerDefs {
		if container != "" && d.Name != container {
			continue
		}

		fmt.Printf("container: %s\n", d.Name)
		for _, e := range d.Environments {
			fmt.Printf("\t%s = %s\n", *e.Name, *e.Value)
		}
		if len(d.Secrets) > 0 {
			paramNames := make([]*string, len(d.Secrets))
			paramNameToSecret := make(map[string]*ecs.Secret)
			for i, e := range d.Secrets {
				paramNames[i] = e.ValueFrom
				paramNameToSecret[*e.ValueFrom] = e
			}
			params, err := w.SSM.GetParameters(paramNames)
			if err != nil {
				return err
			}
			for _, p := range params {
				fmt.Printf("\t%s = %s\n", *paramNameToSecret[p.Name].Name, p.Value)
			}
		}
		fmt.Println()
	}

	return nil
}

func getARN(query string, getter func() ([]awscliwrapper.ARN, error)) (awscliwrapper.ARN, error) {
	arns, err := getter()
	if err != nil {
		return "", err
	}

	if len(arns) == 0 {
		return "", nil
	}

	arn := arns[0]
	if len(arns) > 1 {
		arn, err = selectOneARN(query, arns)
		if err != nil {
			return "", err
		}
	}

	return arn, nil
}

func selectOneARN(query string, arns []awscliwrapper.ARN) (awscliwrapper.ARN, error) {
	names := make([]string, len(arns))
	nameToARN := make(map[string]awscliwrapper.ARN)
	for i, arn := range arns {
		name := arn.Name()
		names[i] = name
		nameToARN[name] = arn
	}

	sort.Strings(names)
	name, err := InputUI.Select(query, names, &input.Options{
		Required: true,
		Loop:     true,
	})
	if err != nil {
		return "", err
	}

	return nameToARN[name], nil
}
