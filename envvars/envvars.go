package envvars

import (
	"regexp"
	"strings"
)

var (
	keyNameRegex = regexp.MustCompile(`\w+=\w*`)
)

type Var struct {
	Name  string
	Value string
}

func (env Var) String() string {
	return env.Name + "=" + env.Value
}

func Split(s string) []Var {
	var vars []Var
	for _, c := range strings.Split(s, ",") {
		if keyNameRegex.MatchString(c) {
			kv := strings.Split(c, "=")
			vars = append(vars, Var{Name: kv[0], Value: kv[1]})
		} else if len(vars) > 0 {
			vars[len(vars)-1].Value += "," + c
		}
	}
	return vars
}
