package cmds

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/fatih/color"
	"github.com/tcnksm/go-input"
	"github.com/urfave/cli"

	"github.com/li-go/awscliwrapper"
)

func ECSCommands() cli.Commands {
	return cli.Commands{
		{
			Name:   "walk",
			Usage:  "walk ECS",
			Action: ActionFunc(walkCluster),
		},
	}
}

func walkCluster(w *awscliwrapper.Wrapper) error {
	clusters, err := w.ECS.ListClusters()
	if err != nil {
		return err
	}

	if len(clusters) == 0 {
		return nil
	}

	cluster := clusters[0]
	if len(clusters) > 1 {
		cluster, err = selectOneARN("select a cluster?", clusters)
		if err != nil {
			return err
		}
	}
	fmt.Printf("cluster: %s\n\n", color.GreenString(cluster.Name()))

	services, err := w.ECS.GetServices(cluster)
	if err != nil {
		return err
	}

	if len(services) == 0 {
		return nil
	}

	service := services[0]
	if len(services) > 1 {
		service, err = selectOneARN("select a service?", services)
		if err != nil {
			return err
		}
	}
	fmt.Printf("service: %s\n\n", color.GreenString(service.Name()))

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

func selectOneARN(query string, arns []awscliwrapper.ARN) (awscliwrapper.ARN, error) {
	names := make([]string, len(arns))
	nameToARN := make(map[string]awscliwrapper.ARN)
	for i, arn := range arns {
		name := arn.Name()
		names[i] = name
		nameToARN[name] = arn
	}

	name, err := InputUI.Select(query, names, &input.Options{
		Required: true,
		Loop:     true,
	})
	if err != nil {
		return "", err
	}

	return nameToARN[name], nil
}
