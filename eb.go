package awscliwrapper

import (
	"time"

	eb "github.com/aws/aws-sdk-go/service/elasticbeanstalk"

	"github.com/li-go/awscliwrapper/envvars"
)

func NewEB(region, profile string) (*EBWrapper, error) {
	sess, err := newSession(region, profile)
	if err != nil {
		return nil, err
	}

	return &EBWrapper{svc: eb.New(sess)}, nil
}

type EBWrapper struct {
	svc *eb.ElasticBeanstalk
}

func (w *EBWrapper) GetApps() ([]*eb.ApplicationDescription, error) {
	i := eb.DescribeApplicationsInput{}
	o, err := w.svc.DescribeApplications(&i)
	if err != nil {
		return nil, err
	}
	return o.Applications, nil
}

func (w *EBWrapper) GetVersions(appName string) ([]*eb.ApplicationVersionDescription, error) {
	i := eb.DescribeApplicationVersionsInput{}
	i.SetApplicationName(appName)
	o, err := w.svc.DescribeApplicationVersions(&i)
	if err != nil {
		return nil, err
	}
	return o.ApplicationVersions, nil
}

func (w *EBWrapper) GetEnvironments(appName string) ([]*eb.EnvironmentDescription, error) {
	i := eb.DescribeEnvironmentsInput{}
	i.SetApplicationName(appName)
	o, err := w.svc.DescribeEnvironments(&i)
	if err != nil {
		return nil, err
	}
	return o.Environments, nil
}

func (w *EBWrapper) GetEnvVars(appName, envName string) ([]envvars.Var, error) {
	i := eb.DescribeConfigurationSettingsInput{}
	i.SetApplicationName(appName)
	i.SetEnvironmentName(envName)
	o, err := w.svc.DescribeConfigurationSettings(&i)
	if err != nil {
		return nil, err
	}
	var vars []envvars.Var
	for _, settings := range o.ConfigurationSettings {
		for _, opt := range settings.OptionSettings {
			if *opt.OptionName == "EnvironmentVariables" {
				vars = append(vars, envvars.Split(*opt.Value)...)
			}
		}
	}
	return vars, nil
}

func (w *EBWrapper) GetEnvResource(envName string) (*eb.EnvironmentResourceDescription, error) {
	i := eb.DescribeEnvironmentResourcesInput{}
	i.SetEnvironmentName(envName)
	o, err := w.svc.DescribeEnvironmentResources(&i)
	if err != nil {
		return nil, err
	}
	return o.EnvironmentResources, nil
}

func (w *EBWrapper) DeployEB(env, ver string) error {
	i := eb.UpdateEnvironmentInput{}
	i.SetEnvironmentName(env)
	i.SetVersionLabel(ver)
	_, err := w.svc.UpdateEnvironment(&i)
	return err
}

func (w *EBWrapper) GetEvents(app, env string, start time.Time, limit int) ([]*eb.EventDescription, error) {
	i := eb.DescribeEventsInput{}
	if len(app) > 0 {
		i.SetApplicationName(app)
	}
	if len(env) > 0 {
		i.SetEnvironmentName(env)
	}
	i.SetStartTime(start)
	i.SetMaxRecords(int64(limit))
	o, err := w.svc.DescribeEvents(&i)
	if err != nil {
		return nil, err
	}
	return o.Events, nil
}
