package access

import (
	"encoding/json"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type GlobPattern struct {
	regexp *regexp.Regexp
}

func (p *GlobPattern) MatchString(s string) bool {
	return p.regexp.MatchString(s)
}

func (p *GlobPattern) String() string {
	return p.regexp.String()
}

func (p *GlobPattern) Set(s string) error {

	glob := strings.Builder{}
	for _, c := range s {
		switch c {
		case '*':
			glob.WriteString(".*")
		case '?', '.', '+', '(', ')', '|', '^', '$', '\\':
			glob.WriteRune('\\')
			glob.WriteRune(c)
		default:
			glob.WriteRune(c)
		}
	}
	regexp, err := regexp.Compile(glob.String())
	if err != nil {
		return err
	}
	p.regexp = regexp
	return nil
}

func (p *GlobPattern) UnmarshalText(text []byte) error {
	return p.Set(string(text))
}

type PatternList []GlobPattern

func (pl *PatternList) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.SequenceNode {
		return &yaml.TypeError{Errors: []string{"expected sequence of strings"}}
	}
	for _, item := range value.Content {
		var pattern GlobPattern
		if err := item.Decode(&pattern); err != nil {
			return err
		}
		*pl = append(*pl, pattern)
	}
	return nil
}

func (pl PatternList) UnmarshalJSON(encoded []byte) error {
	strings := []string{}
	if err := json.Unmarshal(encoded, &strings); err != nil {
		return err
	}
	for _, s := range strings {
		var pattern GlobPattern
		if err := pattern.UnmarshalText([]byte(s)); err != nil {
			return err
		}
		pl = append(pl, pattern)
	}
	return nil
}

func (pl PatternList) MatchString(s string) bool {
	for _, p := range pl {
		if p.MatchString(s) {
			return true
		}
	}
	return false
}
