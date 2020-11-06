package awscliwrapper

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	eb "github.com/aws/aws-sdk-go/service/elasticbeanstalk"
)

type EB struct {
	svc *eb.ElasticBeanstalk
}

func NewEB(sess *session.Session) *EB {
	return &EB{svc: eb.New(sess)}
}

func (e *EB) GetApps() ([]*eb.ApplicationDescription, error) {
	i := eb.DescribeApplicationsInput{}
	o, err := e.svc.DescribeApplications(&i)
	if err != nil {
		return nil, err
	}
	return o.Applications, nil
}

func (e *EB) GetVersions(appName string) ([]*eb.ApplicationVersionDescription, error) {
	i := eb.DescribeApplicationVersionsInput{}
	i.SetApplicationName(appName)
	o, err := e.svc.DescribeApplicationVersions(&i)
	if err != nil {
		return nil, err
	}
	return o.ApplicationVersions, nil
}

func (e *EB) GetEnvironments(appName string) ([]*eb.EnvironmentDescription, error) {
	i := eb.DescribeEnvironmentsInput{}
	i.SetApplicationName(appName)
	o, err := e.svc.DescribeEnvironments(&i)
	if err != nil {
		return nil, err
	}
	return o.Environments, nil
}

func (e *EB) GetEnvVars(appName, envName string) ([]EnvVar, error) {
	i := eb.DescribeConfigurationSettingsInput{}
	i.SetApplicationName(appName)
	i.SetEnvironmentName(envName)
	o, err := e.svc.DescribeConfigurationSettings(&i)
	if err != nil {
		return nil, err
	}
	var vars []EnvVar
	for _, settings := range o.ConfigurationSettings {
		for _, opt := range settings.OptionSettings {
			if *opt.OptionName == "EnvironmentVariables" {
				vars = append(vars, SplitIntoEnvVars(*opt.Value)...)
			}
		}
	}
	return vars, nil
}

func (e *EB) GetEnvResource(envName string) (*eb.EnvironmentResourceDescription, error) {
	i := eb.DescribeEnvironmentResourcesInput{}
	i.SetEnvironmentName(envName)
	o, err := e.svc.DescribeEnvironmentResources(&i)
	if err != nil {
		return nil, err
	}
	return o.EnvironmentResources, nil
}

func (e *EB) DeployEB(env, ver string) error {
	i := eb.UpdateEnvironmentInput{}
	i.SetEnvironmentName(env)
	i.SetVersionLabel(ver)
	_, err := e.svc.UpdateEnvironment(&i)
	return err
}

func (e *EB) GetEvents(app, env string, start time.Time, limit int) ([]*eb.EventDescription, error) {
	i := eb.DescribeEventsInput{}
	if len(app) > 0 {
		i.SetApplicationName(app)
	}
	if len(env) > 0 {
		i.SetEnvironmentName(env)
	}
	i.SetStartTime(start)
	i.SetMaxRecords(int64(limit))
	o, err := e.svc.DescribeEvents(&i)
	if err != nil {
		return nil, err
	}
	return o.Events, nil
}
