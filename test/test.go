package qeconv

import (
	"text/scanner"
	"strings"
)

func removeLineComment(s string, p rune) string {
	var ret []rune
	var l scanner.Scanner
	l.Init(strings.NewReader(s))

	for l.Peek() != scanner.EOF {
		s := l.Next()
		if s == p {
			s = l.Next()
			for s != '\n' && s != scanner.EOF {
				s = l.Next()
			}
			continue
		}
		ret = append(ret, s)
	}

	return string(ret)
}


func cmpIgnoreSpace(str1, str2 string) bool {
	var l1 scanner.Scanner
	var l2 scanner.Scanner
	l1.Init(strings.NewReader(str1))
	l2.Init(strings.NewReader(str2))

	for {
		for l1.Peek() == ' ' || l1.Peek() == '\t' ||
			l1.Peek() == '\r' || l1.Peek() == '\n' {
			l1.Next()
		}
		for l2.Peek() == ' ' || l2.Peek() == '\t' ||
			l2.Peek() == '\r' || l2.Peek() == '\n' {
			l2.Next()
		}
		if l1.Peek() == scanner.EOF || l2.Peek() == scanner.EOF {
			break
		}

		if l1.Next() != l2.Next() {
			return false
		}
	}

	return l1.Peek() == l2.Peek()
}


