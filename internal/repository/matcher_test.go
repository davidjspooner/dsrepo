package repository

import (
	"testing"

	"github.com/davidjspooner/dsmatch/pkg/matcher"
)

func TestNewGlob(t *testing.T) {
	tests := []struct {
		path      []byte
		seperator byte
		expected  matcher.Sequence
		err       bool
	}{
		{
			path:      []byte("foo*bar"),
			seperator: '/',
			expected: matcher.Sequence{
				&matcher.Pattern{Text: []byte("foo")},
				&matcher.Pattern{Text: []byte("*bar"), EndOfText: true},
			},
			err: false,
		},
		{
			path:      []byte("foo**bar"),
			seperator: '/',
			expected:  nil,
			err:       true,
		},
		{
			path:      []byte("foo**/bar"),
			seperator: '/',
			expected: matcher.Sequence{
				&matcher.Pattern{Text: []byte("foo")},
				&matcher.PathParts{Seperator: '/', Min: 0, Max: MaxGlobParts},
				&matcher.Pattern{Text: []byte("/bar"), EndOfText: true},
			},
			err: false,
		},
		{
			path:      []byte("foo**"),
			seperator: '/',
			expected: matcher.Sequence{
				&matcher.Pattern{Text: []byte("foo")},
				&matcher.PathParts{Seperator: '/', Min: 0, Max: MaxGlobParts},
				&matcher.Pattern{EndOfText: true},
			},
			err: false,
		},
		{
			path:      []byte("foo?bar"),
			seperator: '/',
			expected: matcher.Sequence{
				&matcher.Pattern{Text: []byte("foo")},
				&matcher.Pattern{Text: []byte("?bar"), EndOfText: true},
			},
			err: false,
		},
	}

	for _, test := range tests {
		result, err := NewGlob(test.path, test.seperator)
		if (err != nil) != test.err {
			t.Errorf("NewGlob(%q, %q) error = %v, wantErr %v", test.path, test.seperator, err, test.err)
			continue
		}
		if !sequenceEqual(result, test.expected) {
			t.Errorf("NewGlob(%q, %q) = %v, want %v", test.path, test.seperator, result, test.expected)
		}
	}
}

func sequenceEqual(a, b matcher.Sequence) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		common, left, right := a[i].Split(b[i])
		if common == nil || left != nil || right != nil {
			return false
		}
	}
	return true
}
