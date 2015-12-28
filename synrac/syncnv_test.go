package qeconv

import (
	. "github.com/hiwane/qeconv/def"
	"strings"
	"testing"
	"text/scanner"
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

func TestToSyn(t *testing.T) {
	var data = []struct {
		input  string
		expect string
	}{
		{"true:", "true:"},
		{"false:", "false:"},
		{"x>0:", "0<x:"},
		{"x+y>0:", "0<x+y:"},
		{"x-2*y<>0:", "x-2*y <> 0:"},
		{"-x+y<=0:", "(-x)+y<=0:"},
		{"-x+(y*2)<=0:", "(-x)+y*2<=0:"},
		{"-(x+y/3)=0:", "(-(x+y/3)) = 0:"},
		{"-(x+y)*2<=0:", "(-(x+y))*2<=0:"},
		{"Not(y=0):", "Not(y  =  0):"},
		{"And(x<=0, y=0):", "And(x <= 0, y  =  0):"},
		{"And():", "true:"},
		{"Or():", "false:"},
		{"Or(x<=0, y<>0):", "Or(x<=0, y<>0):"},
		{"And(Or(x*y>0,z>0),Or(x*y<=-z,+z>=0)):", "And(Or(0<x*y, 0<z), Or(x*y<=(-z),  0<=(+z))):"},
		{"Or(And(x*y>0,z>0),And(x*y<=-z,+z>=0)):", "Or(And(0<x*y, 0<z), And(x*y<=(-z),  0<=(+z))):"},
		{"Ex(x, x^2=-1):", "Ex([x], x^2 = (-1)):"},
		{"Ex({x}, x^2=-1):", "Ex([x], x^2 = (-1)):"},
		{"Ex([x], x^2=-1):", "Ex([x], x^2 = (-1)):"},
		{"Ex(x, And(x^2=-1, x>0)):", "Ex([x], And(x^2 = (-1), 0<x)):"},
		{"All([x], a*x^2+b*x+c>0):", "All([x], 0<a*x^2+b*x+c):"},
		{"All([x], Ex([y], x+y+a=0)):", "All([x],Ex([y],x+y+a = 0)):"},
		{"Equiv(x<0, y=0):", "Equiv(x<0, y = 0):"},
		{"Impl(x<0, y=0):", "Impl(x<0, y = 0):"},
		{"Repl(y=0, x<0):", "Impl(x<0, y = 0):"},
		{"# comment line\n(1+a)*x+(3+b)*y=0:", "(1+a)*x+(3+b)*y  =  0:"},
		{"x+abs(y+z)=0:", "x+abs(y+z) = 0:"},
		{"x+(-y)=1:", "x+(-y) = 1:"},
		{"x*(-2)=1:", "x*(-2) = 1:"},
	}

	m := new(SynConv)
	parser := NewSynParse()

	for i, p := range data {
		fml, cmts, err := parser.Parse(p.input, false)
		if err != nil {
			t.Errorf("err invalid input=%s\n", p.input)
		}
		co := NewCnvOut(cmts)
		actual0, _ := m.Convert(fml, co)
		actual := removeLineComment(actual0 + m.Sep(), '#')
		if !cmpIgnoreSpace(actual, p.expect) {
			t.Errorf("err %d\nactual=%s\nexpect=%s\ninput=%s\n", i, actual0, p.expect, p.input)
		}
	}
}
