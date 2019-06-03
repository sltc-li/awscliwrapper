package cmds

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"time"

	eb "github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/fatih/color"
	"github.com/tcnksm/go-input"

	"github.com/li-go/awscliwrapper"
)

func describeEB(wrapper *awscliwrapper.EBWrapper, ec2Wrapper *awscliwrapper.EC2Wrapper) error {
	appName, _, err := selectEBApp(wrapper)
	if err != nil {
		return err
	}

	envName, _, err := selectEBEnv(wrapper, appName)
	if err != nil {
		return err
	}

	// print env vars
	vars, err := wrapper.GetEnvVars(appName, envName)
	if err != nil {
		return err
	}
	varStrings := make([]string, len(vars))
	for i, v := range vars {
		varStrings[i] = v.String()
	}
	sort.Strings(varStrings)
	fmt.Printf("Environment Variables:\n%s\n", strings.Join(varStrings, "\n"))

	// print resources
	resourceDesc, err := wrapper.GetEnvResource(envName)
	if err != nil {
		return err
	}
	lbNames := make([]string, len(resourceDesc.LoadBalancers))
	for i, lb := range resourceDesc.LoadBalancers {
		lbNames[i] = *lb.Name
	}
	instanceIDs := make([]string, len(resourceDesc.Instances))
	for i, instance := range resourceDesc.Instances {
		instanceIDs[i] = *instance.Id
	}
	fmt.Printf("\nResources:\nLoad Balancers: %s\nInstances:\n", strings.Join(lbNames, ","))
	instances, err := ec2Wrapper.DescribeInstances(instanceIDs...)
	if err != nil {
		return err
	}
	for _, instance := range instances {
		fmt.Printf("%s: %s / %s\n", *instance.InstanceId, *instance.KeyName, *instance.PublicDnsName)
	}
	return nil
}

func deployEB(wrapper *awscliwrapper.EBWrapper) error {
	appName, _, err := selectEBApp(wrapper)
	if err != nil {
		return err
	}

	envName, envs, err := selectEBEnv(wrapper, appName)
	if err != nil {
		return err
	}

	// select version
	vers, err := wrapper.GetVersions(appName)
	verLabels := make([]string, len(vers))
	for i, ver := range vers {
		var envNames []string
		for _, env := range envs {
			if *env.VersionLabel == *ver.VersionLabel {
				name := *env.EnvironmentName
				if name == envName {
					name = color.GreenString(name)
				}
				envNames = append(envNames, name)
			}
		}
		verLabels[i] = *ver.VersionLabel
		if len(envNames) > 0 {
			verLabels[i] = verLabels[i] + " / " + strings.Join(envNames, ",")
		}
	}
	verLabel, err := ui.Select("select a version?", verLabels, &input.Options{
		Required: true,
		Loop:     true,
	})
	if err != nil {
		return err
	}
	version := strings.Split(verLabel, " / ")[0]
	fmt.Printf("version: %s\n\n", color.GreenString(version))

	// deploy
	query := fmt.Sprintf("application: %s\ndeploy %s to %s?",
		color.GreenString(appName), color.GreenString(envName), color.GreenString(version))
	answer, err := ui.Select(query, []string{"yes", "no"}, &input.Options{
		Required: true,
		Loop:     true,
	})
	if err != nil {
		return err
	}
	if answer == "yes" {
		fmt.Printf("%s ...\n", color.GreenString("deploying"))
		if err := wrapper.DeployEB(envName, version); err != nil {
			return err
		}
		fmt.Printf("%s !!!\n\n", color.GreenString("deployed"))
	}

	// fetch events
	quitCh := make(chan os.Signal)
	signal.Notify(quitCh, os.Interrupt)
	errCh := make(chan error)
	var lastEvents []*eb.EventDescription
	filterEvents := func(evs []*eb.EventDescription) []*eb.EventDescription {
		var filtered []*eb.EventDescription
	loop:
		for _, ev := range evs {
			for _, lev := range lastEvents {
				if *ev.EventDate == *lev.EventDate && *ev.Message == *lev.Message {
					continue loop
				}
			}
			filtered = append(filtered, ev)
		}
		return filtered
	}
	go func() {
		for {
			start := time.Now()
			if len(lastEvents) > 0 {
				start = *lastEvents[len(lastEvents)-1].EventDate
			}
			events, err := wrapper.GetEvents(appName, envName, start, 10)
			if err != nil {
				errCh <- err
				break
			}
			sort.Slice(events, func(i, j int) bool {
				return events[i].EventDate.Sub(*events[j].EventDate) < 0
			})
			for _, ev := range filterEvents(events) {
				fmt.Printf("%s: %s / %s\n", ev.EventDate.Format("2006-01-02 15:04:05.000"), *ev.EnvironmentName, *ev.Message)
			}
			if len(events) > 0 {
				ev := events[len(events)-1]
				if *ev.EnvironmentName == envName && *ev.Message == "Environment update completed successfully." {
					errCh <- nil
					break
				}
			}
			lastEvents = events
			time.Sleep(time.Second)
		}
	}()
	select {
	case <-quitCh:
		return nil
	case err := <-errCh:
		return err
	}
}

func selectEBApp(wrapper *awscliwrapper.EBWrapper) (string, []*eb.ApplicationDescription, error) {
	apps, err := wrapper.GetApps()
	if err != nil {
		return "", nil, err
	}
	appNames := make([]string, len(apps))
	for i, app := range apps {
		appNames[i] = *app.ApplicationName
	}
	sort.Strings(appNames)
	appName, err := ui.Select("select an application?", appNames, &input.Options{
		Required: true,
		Loop:     true,
	})
	if err != nil {
		return "", nil, err
	}
	fmt.Printf("application: %s\n\n", color.GreenString(appName))
	return appName, apps, nil
}

func selectEBEnv(wrapper *awscliwrapper.EBWrapper, appName string) (string, []*eb.EnvironmentDescription, error) {
	envs, err := wrapper.GetEnvironments(appName)
	if err != nil {
		return "", nil, err
	}
	envNames := make([]string, len(envs))
	for i, env := range envs {
		status := *env.Status
		if status == "Ready" {
			status = color.GreenString(status)
		} else {
			status = color.RedString(status)
		}
		envNames[i] = *env.EnvironmentName + " / " + color.GreenString(*env.Status)
	}
	sort.Strings(envNames)
	envName, err := ui.Select("select an environment?", envNames, &input.Options{
		Required: true,
		Loop:     true,
	})
	if err != nil {
		return "", nil, err
	}
	envName = envName[0:strings.Index(envName, " / ")]
	fmt.Printf("environment: %s\n\n", color.GreenString(envName))

	return envName, envs, nil
}
