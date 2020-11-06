package awscliwrapper

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

type Wrapper struct {
	EB  *EB
	EC2 *EC2
	ECS *ECS
	S3  *S3
	SSM *SSM
	IAM *IAM
}

func New() (*Wrapper, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	return &Wrapper{
		EB:  NewEB(sess),
		EC2: NewEC2(sess),
		ECS: NewECS(sess),
		S3:  NewS3(sess),
		SSM: NewSSM(sess),
		IAM: NewIAM(sess),
	}, nil
}
