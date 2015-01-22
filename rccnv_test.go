package qeconv

import (
	"testing"
)

func TestToRC(t *testing.T) {
	var data = []struct {
		input  string
		expect string
	}{
		{"true:", "true"},
		{"false:", "false"},
		{"x>0:", "0<x"},
		{"x+y>0:", "0<x+y"},
		{"x-2*y<>0:", "x-2*y <> 0"},
		{"-x+y<=0:", "-x+y<=0"},
		{"-x+(y*2)<=0:", "-x+y*2<=0"},
		{"-(x+y/3)=0:", "-(x+y/3) = 0"},
		{"-(x+y)*2<=0:", "-(x+y)*2<=0"},
		{"Not(y=0):", "`&not`(y  =  0)"},
		{"And(x<=0, y=0):", "`&and`(x <= 0, y  =  0)"},
		{"And():", "true"},
		{"Or():", "false"},
		{"Or(x<=0, y<>0):", "`&or`(x<=0, y<>0)"},
		{"And(Or(x*y>0,z>0),Or(x*y<=-z,+z>=0)):", "`&and`(`&or`(0<x*y, 0<z), `&or`(x*y<=-z,  0<=+z))"},
		{"Or(And(x*y>0,z>0),And(x*y<=-z,+z>=0)):", "`&or`(`&and`(0<x*y, 0<z), `&and`(x*y<=-z,  0<=+z))"},
		{"Ex([x], x^2=-1):", "`&E`([x]), x^2 = -1"},
		{"Ex(x, And(x^2=-1, x>0)):", "`&E`([x]), `&and`(x^2 = -1, 0<x)"},
		{"All([x], a*x^2+b*x+c>0):", "`&A`([x]), 0<a*x^2+b*x+c"},
		{"All([x], Ex([y], x+y+a=0)):", "`&A`([x]),`&E`([y]),x+y+a = 0"},
		{"Equiv(x<0, y=0):", "`&iff`(x<0, y = 0)"},
		{"Impl(x<0, y=0):", "`&implies`(x<0, y = 0)"},
		{"Repl(y=0, x<0):", "`&implies`(x<0, y = 0)"},
		{"# comment line\n(1+a)*x+(3+b)*y=0:", "(1+a)*x+(3+b)*y  =  0"},
	}

	for i, p := range data {
		actual0, _ := Convert(p.input, "rc", false, 0)
		actual := removeLineComment(actual0, '#')
		if !cmpIgnoreSpace(actual, p.expect+":") {
			t.Errorf("err %d\nactual=%s\nexpect=%s\ninput=%s\n", i, actual0, p.expect, p.input)
		}
	}
}
