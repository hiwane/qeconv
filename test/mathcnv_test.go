package qeconv

import (
	. "github.com/hiwane/qeconv"
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
		{"Or(And(x*y>0,z>0),And(x*y<=-z,+z>=0)):", "0<x*y && 0<z || x*y<=-z &&  0<=+z"},
		{"Ex(x, x^2=-1):", "Exists[{x}, x^2==-1]"},
		{"Ex(x, And(x^2=-1, x>0)):", "Exists[{x}, x^2==-1 && 0<x]"},
		{"All([x], a*x^2+b*x+c>0):", "ForAll[{x}, 0<a*x^2+b*x+c]"},
		{"All([x], Ex([y], x+y+a=0)):", "ForAll[{x},Exists[{y},x+y+a==0]]"},
		{"Equiv(x<0, y=0):", "Equivalent[x<0, y==0]"},
		{"Impl(x<0, y=0):", "Implies[x<0, y==0]"},
		{"Repl(y=0, x<0):", "Implies[x<0, y==0]"},
		{"# comment line\n(1+a)*x+(3+b)*y=0:", "(1+a)*x+(3+b)*y == 0"},
		{"x+abs(y+z)=0:", "x+Abs[y+z]==0"},
	}

	for i, p := range data {
		m, err := Str2cinf("math")
		if err != nil {
			t.Errorf("err str2cnf: %d, str=%s\n", i, p.input)
		}
		parser, err := Str2Parser("syn")
		if err != nil {
			t.Errorf("err str2parser: str=%s\n", p.input)
		}
		actual0, _ := Convert(parser, m, p.input, false, 0)
		actual := removeMathComment(actual0)
		if !cmpIgnoreSpace(actual, p.expect+";") {
			t.Errorf("err %d\nactual=%s\nexpect=%s\ninput=%s\n", i, actual0, p.expect, p.input)
		}
	}
}
