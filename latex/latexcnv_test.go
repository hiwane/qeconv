package qeconv

import (
	. "github.com/hiwane/qeconv/def"
	syn "github.com/hiwane/qeconv/synrac"
	"testing"
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



func TestToLaTeX(t *testing.T) {
	var data = []struct {
		input  string
		expect string
	}{
		{"true:", "\\top"},
		{"false:", "\\bot"},
		{"x>0:", "0<x"},
		{"x+y>0:", "0<x+y"},
		{"x-2*y<>0:", "x-2 y \\neq 0"},
		{"-x+y<=0:", "(-x)+y \\leq 0"},
		{"-x+(y*2)<=0:", "(-x)+y 2 \\leq 0"},
		{"-(x+y/3)=0:", "(-(x+y/3))=0"},
		{"-(x+y)/3=0:", "(-(x+y))/3=0"},
		{"-(x+y)*2<=0:", "(-(x+y)) 2 \\leq 0"},
		{"Not(y=0):", "\\neg (y = 0)"},
		{"And(x<=0, y=0):", "x \\leq 0 \\land y = 0"},
		{"Or(x<=0, y<>0):", "x \\leq 0 \\lor y \\neq 0"},
		{"Impl(x>0, y<=0):", "0 < x \\rightarrow y \\leq 0"},
		{"Repl(x<0, y>=0):", "0 \\leq y \\rightarrow x < 0"},
		{"Equiv(x>0, y<=0):", "0 < x \\leftrightarrow y \\leq 0"},
		{"((x-1)^3+y)^(1/2) < 0:", "((x-1)^{3}+y)^{1/2} < 0"},
		{"Ex(x, x^2=-1):", "\\exists x (x^{2}=(-1))"},
		{"All([x], a*x^2+b*x+c>0):", "\\forall x (0<a x^{2}+b x+c)"},
		{"All([x], Ex([y], x+y+a=0)):", "\\forall x \\exists y(x+y+a=0)"},
		{"abs(x)>0:", "0<|x|"},
	}

	m := new(LatexConv)
	parser := syn.NewSynParse()
	for _, p := range data {
		t.Log("inp=%s\n", p.input)
		fml, cmts, err := parser.Parse(p.input, false)
		if err != nil {
			t.Errorf("err invalid input=%s\n", p.input)
		}

		co := NewCnvOut(cmts)
		actual0, _ := m.Convert(fml, co)
		t.Log("ac0=%s\n", actual0)
		actual := removeLineComment(actual0 + m.Sep(), '%')
		t.Log("rem=%s\n", actual)
		t.Log("exp=%s\n", p.expect)
		if !cmpIgnoreSpace(actual, p.expect+"\\\\") {
			t.Errorf("err actual=%s\nexpect=%s\ninput=%s\n", actual0, p.expect, p.input)
		}
	}
}
