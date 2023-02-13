package awscliwrapper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type ECS struct {
	svc *ecs.ECS
}

func NewECS(sess *session.Session) *ECS {
	return &ECS{svc: ecs.New(sess)}
}

func (e *ECS) ListClusters() ([]ARN, error) {
	o, err := e.svc.ListClusters(&ecs.ListClustersInput{})
	if err != nil {
		return nil, err
	}

	arns := make([]ARN, len(o.ClusterArns))
	for i, arn := range o.ClusterArns {
		arns[i] = ARN(*arn)
	}

	return arns, err
}

func (e *ECS) GetServices(clusterName string) ([]ARN, error) {
	arns, nextToken, err := e._getSerivces(clusterName, nil)
	if err != nil {
		return nil, err
	}

	for nextToken != nil {
		var nextARNs []ARN
		nextARNs, nextToken, err = e._getSerivces(clusterName, nextToken)
		if err != nil {
			return nil, err
		}
		arns = append(arns, nextARNs...)
	}

	return arns, err
}

func (e *ECS) _getSerivces(clusterName string, nextToken *string) ([]ARN, *string /* nextToken */, error) {
	input := &ecs.ListServicesInput{Cluster: aws.String(clusterName)}
	if nextToken != nil {
		input.NextToken = nextToken
	}
	o, err := e.svc.ListServices(input)
	if err != nil {
		return nil, nil, err
	}

	arns := make([]ARN, len(o.ServiceArns))
	for i, arn := range o.ServiceArns {
		arns[i] = ARN(*arn)
	}

	return arns, o.NextToken, err
}

func (e *ECS) GetTaskDefinition(clusterName, serviceName string) (ARN, error) {
	o, err := e.svc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(clusterName),
		Services: []*string{aws.String(serviceName)},
	})
	if err != nil {
		return "", err
	}

	return ARN(*o.Services[0].TaskDefinition), nil
}

func (e *ECS) GetTaskDefinitions(family string) ([]ARN, error) {
	o, err := e.svc.ListTaskDefinitions(&ecs.ListTaskDefinitionsInput{
		FamilyPrefix: aws.String(family),
		Sort:         aws.String("DESC"),
	})
	if err != nil {
		return nil, err
	}

	arns := make([]ARN, len(o.TaskDefinitionArns))
	for i, a := range o.TaskDefinitionArns {
		arns[i] = ARN(*a)
	}
	return arns, nil
}

func (e *ECS) GetTasks(clusterName, serviceName string) ([]ARN, error) {
	o, err := e.svc.ListTasks(&ecs.ListTasksInput{
		Cluster:     aws.String(clusterName),
		ServiceName: aws.String(serviceName),
	})

	if err != nil {
		return nil, err
	}

	arns := make([]ARN, len(o.TaskArns))
	for i, arn := range o.TaskArns {
		arns[i] = ARN(*arn)
	}
	return arns, nil
}

func (e *ECS) GetContainerDefinitions(taskDef ARN) ([]ContainerDefinition, error) {
	o, err := e.svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{TaskDefinition: taskDef.AWSString()})
	if err != nil {
		return nil, err
	}

	defs := make([]ContainerDefinition, len(o.TaskDefinition.ContainerDefinitions))
	for i, d := range o.TaskDefinition.ContainerDefinitions {
		defs[i] = ContainerDefinition{
			Name:         *d.Name,
			Image:        *d.Image,
			Environments: d.Environment,
			Secrets:      d.Secrets,
		}
	}

	return defs, nil
}

func (e *ECS) Deploy(cluster, service, taskDefinition string) error {
	_, err := e.svc.UpdateService(&ecs.UpdateServiceInput{
		Cluster:            aws.String(cluster),
		Service:            aws.String(service),
		TaskDefinition:     aws.String(taskDefinition),
		ForceNewDeployment: aws.Bool(true),
	})
	if err != nil {
		return err
	}

	return e.svc.WaitUntilServicesStable(&ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: []*string{aws.String(service)},
	})
}

type ContainerDefinition struct {
	Name         string
	Image        string
	Environments []*ecs.KeyValuePair
	Secrets      []*ecs.Secret
}
