package access

import "strings"

type ErrorList []error

func (el ErrorList) Error() string {
	if el == nil || len(el) == 0 {
		return ""
	}
	sb := strings.Builder{}
	for n, err := range el {
		if n > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}
