package awscliwrapper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func newSession(region string, profile string) (*session.Session, error) {
	var cfg = aws.NewConfig()
	cfg = cfg.WithRegion(region)
	cfg = cfg.WithCredentials(credentials.NewSharedCredentials("", profile))
	return session.NewSession(cfg)
}
