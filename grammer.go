package sdptransform

import (
	"regexp"
)

type Rule struct {
	Name       string
	Push       string
	Reg        *regexp.Regexp
	Names      []string
	Types      []rune
	Format     string
	FormatFunc interface{}
}

var RulesMap map[byte][]*Rule = map[byte][]*Rule{
	'v': []*Rule{
		&Rule{
			Name:   "version",
			Push:   "",
			Reg:    regexp.MustCompile("^(\\d*)$"),
			Names:  []string{},
			Types:  []rune{'d'},
			Format: "%d",
		},
	},
}
