package awscliwrapper

import (
	"strings"
)

type ARN string

func (a ARN) Name() string {
	return string(a)[strings.Index(string(a), "/")+1:]
}

func (a ARN) AWSString() *string {
	s := string(a)
	return &s
}
