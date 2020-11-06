package awscliwrapper

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2 struct {
	svc *ec2.EC2
}

func NewEC2(sess *session.Session) *EC2 {
	return &EC2{svc: ec2.New(sess)}
}

func (e *EC2) DescribeInstances(instanceIDs ...string) ([]*ec2.Instance, error) {
	if len(instanceIDs) == 0 {
		return nil, nil
	}

	i := ec2.DescribeInstancesInput{}
	var pInstanceIDs []*string
	for i := range instanceIDs {
		pInstanceIDs = append(pInstanceIDs, &instanceIDs[i])
	}
	i.SetInstanceIds(pInstanceIDs)
	o, err := e.svc.DescribeInstances(&i)
	if err != nil {
		return nil, err
	}
	var instances []*ec2.Instance
	for _, r := range o.Reservations {
		instances = append(instances, r.Instances...)
	}
	return instances, nil
}
