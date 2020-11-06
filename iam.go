package awscliwrapper

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAM struct {
	svc *iam.IAM
}

func NewIAM(sess *session.Session) *IAM {
	return &IAM{svc: iam.New(sess)}
}

func (i *IAM) GetCurrentUser() (*IAMUser, error) {
	o, err := i.svc.GetUser(&iam.GetUserInput{})
	if err != nil {
		return nil, err
	}

	return &IAMUser{
		ID:   *o.User.UserId,
		Name: *o.User.UserName,
		ARN:  ARN(*o.User.Arn),
	}, nil
}

type IAMUser struct {
	ID   string
	Name string
	ARN  ARN
}
