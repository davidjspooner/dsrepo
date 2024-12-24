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
				tm := &matcher.Pattern{Text: path[base:i]}
				matchers = append(matchers, tm)
				base = i
			}
			if i+1 < len(path) && path[i+1] == '*' {
				if i+2 < len(path) && path[i+2] != seperator {
					return nil, fmt.Errorf("invalid glob %q", path)
				}
				tm := &matcher.PathParts{Seperator: seperator, Min: 0, Max: MaxGlobParts}
				matchers = append(matchers, tm)
				i++
				base = i + 1
			}
		case '?':
			if i > base {
				tm := &matcher.Pattern{Text: path[base:i]}
				matchers = append(matchers, tm)
				base = i
			}
		default:
		}
		i++
	}
	matchers = append(matchers, &matcher.Pattern{Text: path[base:], EndOfText: true})
	return matchers, nil
}
