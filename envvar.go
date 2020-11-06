package awscliwrapper

import (
	"regexp"
	"strings"
)

var (
	keyNameRegex = regexp.MustCompile(`\w+=\w*`)
)

type EnvVar struct {
	Name  string
	Value string
}

func (ev EnvVar) String() string {
	return ev.Name + "=" + ev.Value
}

func SplitIntoEnvVars(s string) []EnvVar {
	var vars []EnvVar
	for _, c := range strings.Split(s, ",") {
		if keyNameRegex.MatchString(c) {
			kv := strings.Split(c, "=")
			vars = append(vars, EnvVar{Name: kv[0], Value: kv[1]})
		} else if len(vars) > 0 {
			vars[len(vars)-1].Value += "," + c
		}
	}
	return vars
}
