package awscliwrapper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	svc *s3.S3
}

func NewS3(sess *session.Session) *S3 {
	return &S3{svc: s3.New(sess)}
}

func (s *S3) GetBuckets() ([]*s3.Bucket, error) {
	i := s3.ListBucketsInput{}
	o, err := s.svc.ListBuckets(&i)
	if err != nil {
		return nil, err
	}
	return o.Buckets, nil
}

func (s *S3) GetObjects(bucket string, prefix string) ([]string, []*s3.Object, error) {
	var prefixes []string
	var objects []*s3.Object
	var cToken *string
	for {
		i := s3.ListObjectsV2Input{}
		i.SetBucket(bucket)
		i.SetPrefix(prefix)
		i.SetDelimiter("/")
		i.SetMaxKeys(50)
		if cToken != nil {
			i.SetContinuationToken(*cToken)
		}
		o, err := s.svc.ListObjectsV2(&i)
		if err != nil {
			return nil, nil, err
		}
		if prefixes == nil {
			prefixes = make([]string, len(o.CommonPrefixes))
			for i, p := range o.CommonPrefixes {
				prefixes[i] = *p.Prefix
			}
		}
		objects = append(objects, o.Contents...)
		if o.NextContinuationToken == nil {
			break
		}
		cToken = o.NextContinuationToken
	}
	return prefixes, objects, nil
}

// TODO: upload s3 object
func (s *S3) getObject(bucket, key string) error {
	i := s3.GetObjectInput{}
	i.SetBucket(bucket)
	i.SetKey(key)
	o, err := s.svc.GetObject(&i)
	if err != nil {
		return err
	}
	fmt.Println(o)
	return nil
}

func detectContentType(file *os.File) (string, error) {
	buff := make([]byte, 512)
	n, err := file.Read(buff)
	if err != nil && err != io.EOF {
		return "", err
	}
	file.Seek(0, 0)
	return http.DetectContentType(buff[:n]), nil
}

func (s *S3) UploadFile(from, to string) error {
	file, err := os.Open(from)
	if err != nil {
		return err
	}
	contentType, err := detectContentType(file)
	contentType = "text/css"
	if err != nil {
		return err
	}
	ss := strings.Split(to, "/")
	bucket, keyPrefix := ss[0], strings.Join(ss[1:], "/")
	key := path.Join(keyPrefix, path.Base(from))
	i := s3.PutObjectInput{}
	i.SetBucket(bucket)
	i.SetKey(key)
	i.SetBody(file)
	i.SetACL("public-read")
	i.SetContentType(contentType)
	_, err = s.svc.PutObject(&i)
	return err
}
