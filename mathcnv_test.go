package qeconv

import (
	"strings"
	"testing"
	"text/scanner"
)

func removeMathComment(s string) string {
	var ret []rune
	var l scanner.Scanner
	l.Init(strings.NewReader(s))

	for l.Peek() != scanner.EOF {
		s := l.Peek()
		if s == '(' {
			l.Next()
			if l.Peek() == '*' {
				// commment

				l.Next()
				for {
					if l.Peek() == '*' {
						l.Next()
						if l.Peek() == ')' {
							l.Next()
							break
						}
					}
					l.Next()
				}
				continue
			}
		} else {
			l.Next()
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

func TestToMath(t *testing.T) {
	var data = []struct {
		input  string
		expect string
	}{
		{"true:", "True"},
		{"false:", "False"},
		{"x>0:", "0<x"},
		{"x+y>0:", "0<x+y"},
		{"x-2*y<>0:", "x-2*y != 0"},
		{"-x+y<=0:", "-x+y<=0"},
		{"-x+(y*2)<=0:", "-x+y*2<=0"},
		{"-(x+y/3)=0:", "-(x+y/3)==0"},
		{"-(x+y)*2<=0:", "-(x+y)*2<=0"},
		{"Not(y=0):", "Not[y == 0]"},
		{"And(x<=0, y=0):", "x <= 0 && y == 0"},
		{"Or(x<=0, y<>0):", "x<=0 || y!=0"},
		{"And(Or(x*y>0,z>0),Or(x*y<=-z,+z>=0)):", "(0<x*y || 0<z) && (x*y<=-z ||  0<=+z)"},
		{"Ex(x, x^2=-1):", "Exists[{x}, x^2==-1]"},
		{"All([x], a*x^2+b*x+c>0):", "ForAll[{x}, 0<a*x^2+b*x+c]"},
		{"All([x], Ex([y], x+y+a=0)):", "ForAll[{x},Exists[{y},x+y+a==0]]"},
		{"# comment line\n(1+a)*x+(3+b)*y=0:", "(1+a)*x+(3+b)*y == 0"},
		{"x+abs(y+z)=0:", "x+Abs[y+z]==0"},
	}

	for i, p := range data {
		actual0 := ToMath(p.input)
		actual := removeMathComment(actual0)
		if !cmpIgnoreSpace(actual, p.expect) {
			t.Errorf("err %d\nactual=%s\nexpect=%s\ninput=%s\n", i, actual0, p.expect, p.input)
		}
	}
}
