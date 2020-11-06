package cmds

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fatih/color"
	"github.com/tcnksm/go-input"
	"github.com/urfave/cli"

	"github.com/li-go/awscliwrapper"
)

func S3Commands() cli.Commands {
	return cli.Commands{
		{
			Name:   "walk",
			Usage:  "walk S3",
			Action: ActionFunc(walkS3),
		},
	}
}

func walkS3(w *awscliwrapper.Wrapper) error {
	// select bucket
	buckets, err := w.S3.GetBuckets()
	if err != nil {
		return err
	}
	var bucketNames []string
	for _, b := range buckets {
		if skipBucket(b) {
			continue
		}
		bucketNames = append(bucketNames, *b.Name)
	}
	sort.Strings(bucketNames)
	bucketName, err := InputUI.Select("select a bucket?", bucketNames, &input.Options{
		Required: true,
		Loop:     true,
	})
	if err != nil {
		return err
	}
	fmt.Printf("bucket: %s\n\n", color.GreenString(bucketName))

	// print object
	var prefix string
	for {
		prefixes, objects, err := w.S3.GetObjects(bucketName, prefix)
		if err != nil {
			return err
		}
		if len(prefixes) == 0 && len(objects) == 0 {
			return nil
		}
		for _, obj := range objects {
			if *obj.Key == prefix {
				continue
			}
			fmt.Println(s3URL(bucketName, *obj.Key))
		}

		if len(prefixes) == 0 {
			return nil
		}
		sort.Strings(prefixes)
		selectedPrefix, err := InputUI.Select("\nselect a prefix?", prefixes, &input.Options{
			Required: true,
			Loop:     true,
		})
		if err != nil {
			return err
		}
		fmt.Printf("prefix: %s\n\n", color.GreenString(selectedPrefix))
		prefix = selectedPrefix
	}
}

func s3URL(bucket, objectKey string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, objectKey)
}

func getFiles(dir string) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var allFiles []string
	for _, info := range fileInfos {
		fullName := path.Join(dir, info.Name())
		if info.IsDir() {
			files, err := getFiles(fullName)
			if err != nil {
				return nil, err
			}
			allFiles = append(allFiles, files...)
			continue
		}
		allFiles = append(allFiles, fullName)
	}
	return allFiles, nil
}

var (
	skipBucketPrefixes = []string{"cloudtrail-", "elasticbeanstalk-", "aws-logs-", "cm-members-", "cf-templates-"}
	skipBucketSuffix   = []string{"-log", "-logs", "-accesslog"}
)

func skipBucket(bucket *s3.Bucket) bool {
	for _, p := range skipBucketPrefixes {
		if strings.HasPrefix(*bucket.Name, p) {
			return true
		}
	}
	for _, s := range skipBucketSuffix {
		if strings.HasSuffix(*bucket.Name, s) {
			return true
		}
	}
	return false
}
