package awscliwrapper

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

func NewEC2(region, profile string) (*EC2Wrapper, error) {
	sess, err := newSession(region, profile)
	if err != nil {
		return nil, err
	}

	return &EC2Wrapper{svc: ec2.New(sess)}, nil
}

type EC2Wrapper struct {
	svc *ec2.EC2
}

func (w *EC2Wrapper) DescribeInstances(instanceIDs ...string) ([]*ec2.Instance, error) {
	if len(instanceIDs) == 0 {
		return nil, nil
	}

	i := ec2.DescribeInstancesInput{}
	var pInstanceIDs []*string
	for i := range instanceIDs {
		pInstanceIDs = append(pInstanceIDs, &instanceIDs[i])
	}
	i.SetInstanceIds(pInstanceIDs)
	o, err := w.svc.DescribeInstances(&i)
	if err != nil {
		return nil, err
	}
	var instances []*ec2.Instance
	for _, r := range o.Reservations {
		instances = append(instances, r.Instances...)
	}
	return instances, nil
}
