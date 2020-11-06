package awscliwrapper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSM struct {
	svc *ssm.SSM
}

func NewSSM(sess *session.Session) *SSM {
	return &SSM{svc: ssm.New(sess)}
}

func (s *SSM) GetParameters(names []*string) ([]SSMParameter, error) {
	o, err := s.svc.GetParameters(&ssm.GetParametersInput{Names: names, WithDecryption: aws.Bool(true)})
	if err != nil {
		return nil, err
	}

	params := make([]SSMParameter, len(o.Parameters))
	for i, p := range o.Parameters {
		params[i] = SSMParameter{
			Name:  *p.Name,
			Type:  *p.Type,
			Value: *p.Value,
			ARN:   ARN(*p.ARN),
		}
	}

	return params, err
}

type SSMParameter struct {
	Name  string
	Type  string
	Value string
	ARN   ARN
}
