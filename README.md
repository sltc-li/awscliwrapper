awscli-wrapper
==============

### Install with `go get`
```
$ go get github.com/li-go/awscliwrapper/awsw
```

### How to use
```
$ awsw
NAME:
   awsw - a simple wrapper command for awscli

USAGE:
   awsw [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     eb-desc    describe elasticbeanstalk
     eb-deploy  deploy elasticbeanstalk
     s3-ls
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --region value   aws region (default: "ap-northeast-1")
   --profile value  aws profile (default: "default")
   --fish           generate fish completion
   --help, -h       show help
   --version, -v    print the version
```
