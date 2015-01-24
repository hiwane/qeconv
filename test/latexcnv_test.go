package qeconv

import (
	. "github.com/hiwane/qeconv"
	"testing"
)

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
		{"-x+y<=0:", "-x+y \\leq 0"},
		{"-x+(y*2)<=0:", "-x+y 2 \\leq 0"},
		{"-(x+y/3)=0:", "-(x+y/3)=0"},
		{"-(x+y)/3=0:", "-(x+y)/3=0"},
		{"-(x+y)*2<=0:", "-(x+y) 2 \\leq 0"},
		{"Not(y=0):", "\\neg (y = 0)"},
		{"And(x<=0, y=0):", "x \\leq 0 \\land y = 0"},
		{"Or(x<=0, y<>0):", "x \\leq 0 \\lor y \\neq 0"},
		{"Impl(x>0, y<=0):", "0 < x \\rightarrow y \\leq 0"},
		{"Repl(x<0, y>=0):", "0 \\leq y \\rightarrow x < 0"},
		{"Equiv(x>0, y<=0):", "0 < x \\leftrightarrow y \\leq 0"},
		{"Ex(x, x^2=-1):", "\\exists x (x^{2}=-1)"},
		{"All([x], a*x^2+b*x+c>0):", "\\forall x (0<a x^{2}+b x+c)"},
		{"All([x], Ex([y], x+y+a=0)):", "\\forall x \\exists y(x+y+a=0)"},
		{"abs(x)>0:", "0<|x|"},
	}

	for _, p := range data {
		t.Log("inp=%s\n", p.input)
		m, err := Str2cinf("tex")
		if err != nil {
			t.Errorf("err str2cnf: str=%s\n", p.input)
		}
		parser, err := Str2Parser("syn")
		if err != nil {
			t.Errorf("err str2parser: str=%s\n", p.input)
		}
		actual0, _ := Convert(parser, m, p.input, false, 0)
		t.Log("ac0=%s\n", actual0)
		actual := removeLineComment(actual0, '%')
		t.Log("rem=%s\n", actual)
		t.Log("exp=%s\n", p.expect)
		if !cmpIgnoreSpace(actual, p.expect+"\\\\") {
			t.Errorf("err actual=%s\nexpect=%s\ninput=%s\n", actual0, p.expect, p.input)
		}
	}
}
