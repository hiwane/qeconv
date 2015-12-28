package qeconv

import (
	syn "github.com/hiwane/qeconv/synrac"
	"testing"
)

func TestFromSmt2(t *testing.T) {
	var data = []struct {
		smt2 string
		syn  string
	}{
		{"(assert (= x y))", "x=y"},
		{"(assert (> x y))", "x>y"},
		{"(assert (< x y))", "x<y"},
		{"(assert (>= x y))", "x>=y"},
		{"(assert (<= x y))", "x<=y"},
		{"(assert (<= (+ x 3) (* y 4)))", "x+3<=y*4"},
		{"(assert (not (<= x y)))", "Not(x<=y)"},
		{"(assert (and (<= x y) (> z 1)))", "And(x<=y,z>1)"},
		{"(assert (or (<= x y) (> z 1)))", "Or(x<=y,z>1)"},
		{"(assert (=> (<= x y) (> z 1)))", "Impl(x<=y,z>1)"},
		{"(assert (implies (<= x y) (> z 1)))", "Impl(x<=y,z>1)"},
		{"(assert (exists ((x Real)) (= x 0)))", "Ex([x], x=0)"},
		{"(assert (forall ((x Real)) (< x 0)))", "All([x], x<0)"},
		// uni-minus
		{"(assert (= x (- 2)))", "x = -2"},
		{"(assert (= x (* u (- 2))))", "x = u*(-2)"},
		{"(assert (= y (+ 3 (- z))))", "y = 3 + (- z)"},
		// let
		{"(assert (let ((x 3)) (= y x)))", "y = 3"},
		{"(assert (let ((x (> y 3))) (and (= y 0) x)))", "And(y=0,y>3)"},
	}

	synp := syn.NewSynParse()
	smtp := NewSmt2Parse()

	for i, p := range data {
		fmlt, _, err := smtp.Parse(p.smt2, false)
		if err != nil {
			t.Errorf("err %v\n", err)
			t.Errorf("[%d] err invalid input1=%s\n", i, p.smt2)
			continue
		}

		fmly, _, err := synp.Parse(p.syn + ":", false)
		if err != nil {
			t.Errorf("err %v\n", err)
			t.Errorf("[%d] err invalid input2=%s\n", i, p.syn)
			continue
		}

		// fmlt と fmly は同じもの
		if !fmlt.Equals(fmly) {
			t.Errorf("[%d] err invalid input3=%s\n", i, p.smt2)
		}
	}
}

