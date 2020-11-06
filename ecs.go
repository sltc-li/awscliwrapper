package awscliwrapper

import (
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

func (e *ECS) GetServices(cluster ARN) ([]ARN, error) {
	o, err := e.svc.ListServices(&ecs.ListServicesInput{Cluster: cluster.AWSString()})
	if err != nil {
		return nil, err
	}

	arns := make([]ARN, len(o.ServiceArns))
	for i, arn := range o.ServiceArns {
		arns[i] = ARN(*arn)
	}

	return arns, err
}

func (e *ECS) GetTaskDefinition(cluster, service ARN) (ARN, error) {
	o, err := e.svc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  cluster.AWSString(),
		Services: []*string{service.AWSString()},
	})
	if err != nil {
		return "", nil
	}

	return ARN(*o.Services[0].TaskDefinition), nil
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

type ContainerDefinition struct {
	Name         string
	Image        string
	Environments []*ecs.KeyValuePair
	Secrets      []*ecs.Secret
}
