package repository

import (
	"fmt"

	"github.com/davidjspooner/dsmatch/pkg/matcher"
)

var MaxGlobParts = 100

func NewGlob(path []byte, seperator byte) (matcher.Sequence, error) {
	matchers := make(matcher.Sequence, 0, 2)

	base := 0
	for i := 0; i < len(path); {
		switch path[i] {
		case '*':
			if i > base {
				matchers = append(matchers, &matcher.Text{Pattern: path[base:i], Seperator: seperator})
				base = i
			}
			if i+1 < len(path) && path[i+1] == '*' {
				if i+2 < len(path) && path[i+2] != seperator {
					return nil, fmt.Errorf("invalid glob %q", path)
				}
				matchers = append(matchers, &matcher.PathParts{Seperator: seperator, Min: 0, Max: MaxGlobParts})
				i++
				base = i + 1
			}
		case '?':
			if i > base {
				matchers = append(matchers, &matcher.Text{Pattern: path[base:i], Seperator: seperator})
				base = i
			}
		default:
		}
		i++
	}
	if base < len(path) {
		matchers = append(matchers, &matcher.Text{Pattern: path[base:], Seperator: seperator})
	}
	matchers = append(matchers, &matcher.EndOfText{})

	return matchers, nil
}
