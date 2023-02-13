package cmds

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

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
		{
			Name:  "exec",
			Usage: "execute-command in container",
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
				cli.StringFlag{
					Name:  "command",
					Usage: "command to execute in container",
					Value: "/bin/bash",
				},
			},
			Action: ActionFuncWithContext(execContainer),
		},
		{
			Name:  "deploy",
			Usage: "start new deployment of service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster",
					Usage: "cluster name",
				},
				cli.StringFlag{
					Name:  "service",
					Usage: "service name",
				},
			},
			Action: ActionFuncWithContext(deployService),
		},
	}
}

func walkCluster(c *cli.Context, w *awscliwrapper.Wrapper) error {
	containerDefs, err := listContainers(c, w)
	if err != nil {
		return err
	}

	for _, d := range containerDefs {
		fmt.Printf("container: %s\n", d.Name)
		for _, e := range d.Environments {
			fmt.Printf("\t%s = %s\n", *e.Name, *e.Value)
		}
		if len(d.Secrets) > 0 {
			names := make([]*string, len(d.Secrets))
			valueFroms := make([]*string, len(d.Secrets))
			for i, e := range d.Secrets {
				names[i] = e.Name
				valueFroms[i] = e.ValueFrom
			}
			params, err := w.SSM.GetParameters(valueFroms)
			if err != nil {
				return err
			}
			paramNameToSecretValue := make(map[string]string)
			for _, p := range params {
				paramNameToSecretValue[p.Name] = p.Value
			}
			for _, name := range names {
				fmt.Printf("\t%s = %s\n", *name, paramNameToSecretValue[*name])
			}
		}
		fmt.Println()
	}

	return nil
}

func execContainer(c *cli.Context, w *awscliwrapper.Wrapper) error {
	containerDefs, err := listContainers(c, w)
	if err != nil {
		return err
	}

	container := containerDefs[0].Name
	if len(containerDefs) > 1 {
		var names []string
		for _, c := range containerDefs {
			names = append(names, c.Name)
		}

		sort.Strings(names)
		name, err := InputUI.Select("select a container?", names, &input.Options{
			Required: true,
			Loop:     true,
		})
		if err != nil {
			return err
		}
		container = name
	}
	fmt.Printf("container: %s\n\n", color.GreenString(container))

	tasks, err := w.ECS.GetTasks(c.String("cluster"), c.String("service"))
	if err != nil {
		return nil
	}
	taskID := strings.Split(string(tasks[0]), "/")[2]

	exeCommand := fmt.Sprintf(
		`aws ecs execute-command --cluster %s --task %s --container %s --command "%s" --interactive`,
		c.String("cluster"),
		taskID,
		container,
		c.String("command"),
	)

	cmd := exec.Command("sh", "-c", exeCommand)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func deployService(c *cli.Context, w *awscliwrapper.Wrapper) error {
	taskDef, err := getTaskDefinition(c, w)
	if err != nil {
		return err
	}

	taskDef = taskDef[strings.Index(string(taskDef), "/")+1:]
	family := taskDef[:strings.LastIndex(string(taskDef), ":")]
	taskDefs, err := w.ECS.GetTaskDefinitions(string(family))
	if err != nil {
		return err
	}
	latestTaskDef := taskDefs[0][strings.Index(string(taskDefs[0]), "/")+1:]
	fmt.Printf("latest task definition: %s\n\n", color.GreenString(string(latestTaskDef)))

	if latestTaskDef != taskDef {
		fmt.Println("deploying to latest task definition ...")
		return w.ECS.Deploy(c.String("cluster"), c.String("service"), string(latestTaskDef))
	}
	return nil
}

func listContainers(c *cli.Context, w *awscliwrapper.Wrapper) ([]awscliwrapper.ContainerDefinition, error) {
	container := c.String("container")

	taskDef, err := getTaskDefinition(c, w)
	if err != nil {
		return nil, err
	}

	defs, err := w.ECS.GetContainerDefinitions(taskDef)
	if err != nil {
		return nil, err
	}

	var filtered []awscliwrapper.ContainerDefinition
	for _, def := range defs {
		if container != "" && def.Name != container {
			continue
		}
		filtered = append(filtered, def)
	}
	return filtered, nil
}

func getTaskDefinition(c *cli.Context, w *awscliwrapper.Wrapper) (awscliwrapper.ARN, error) {
	cluster := c.String("cluster")
	service := c.String("service")

	if cluster == "" {
		arn, err := getARN("select a cluster?", w.ECS.ListClusters)
		if err != nil {
			return "", err
		}
		cluster = arn.Name()
		c.Set("cluster", cluster)
	}
	fmt.Printf("cluster: %s\n\n", color.GreenString(cluster))

	if service == "" {
		arn, err := getARN("select a service?", func() ([]awscliwrapper.ARN, error) {
			return w.ECS.GetServices(cluster)
		})
		if err != nil {
			return "", err
		}
		service = arn.Name()
		c.Set("service", service)
	}
	fmt.Printf("service: %s\n\n", color.GreenString(service))

	taskDef, err := w.ECS.GetTaskDefinition(cluster, service)
	if err != nil {
		return "", err
	}
	fmt.Printf("task definition: %s\n\n", color.GreenString(taskDef.Name()))

	return taskDef, nil
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
