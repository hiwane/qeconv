package qeconv

import (
	"strings"
)

type LaTeXConv struct {
}

func (m *LaTeXConv) q(f Formula, qop string, co *cnv_out) {
	vars := f.args[0]
	for i := 0; i < len(vars.args); i++ {
		co.append(" " + qop + " ")
		conv2(vars.args[i], m, co)
	}

	if f.args[1].cmd == ALL || f.args[1].cmd == EX {
		conv2(f.args[1], m, co)
	} else {
		co.append("(")
		conv2(f.args[1], m, co)
		co.append(")")
	}
}

func (m *LaTeXConv) All(f Formula, co *cnv_out) {
	m.q(f, "\\forall", co)
}

func (m *LaTeXConv) Ex(f Formula, co *cnv_out) {
	m.q(f, "\\exists", co)
}

func (m *LaTeXConv) And(f Formula, co *cnv_out) {
	infix(f, m, " \\land ", co)
}

func (m *LaTeXConv) Or(f Formula, co *cnv_out) {
	infix(f, m, " \\lor ", co)
}

func (m *LaTeXConv) Not(f Formula, co *cnv_out) {
	prefix(f, m, "\\neg (", ")", co)
}

func (m *LaTeXConv) Impl(f Formula, co *cnv_out) {
	infix(f, m, " \\rightarrow ", co)
}

func (m *LaTeXConv) Equiv(f Formula, co *cnv_out) {
	infix(f, m, " \\leftrightarrow ", co)
}

func (m *LaTeXConv) Leop(f Formula, co *cnv_out) {
	infix(f, m, " \\leq ", co)
}

func (m *LaTeXConv) Ltop(f Formula, co *cnv_out) {
	infix(f, m, " < ", co)
}

func (m *LaTeXConv) Eqop(f Formula, co *cnv_out) {
	infix(f, m, " = ", co)
}

func (m *LaTeXConv) Neop(f Formula, co *cnv_out) {
	infix(f, m, " \\neq ", co)
}

func (m *LaTeXConv) List(f Formula, co *cnv_out) {
	prefix(f, m, "[", "]", co)
}

func (m *LaTeXConv) Ftrue() string {
	return "\\top"
}

func (m *LaTeXConv) Ffalse() string {
	return "\\bot"
}
func (m *LaTeXConv) Comment(str string) string {
	return "% " + str
}

func (m *LaTeXConv) Plus(fml Formula, co *cnv_out) {
	mop(fml, m, "+", co)
}
func (m *LaTeXConv) Minus(fml Formula, co *cnv_out) {
	mop(fml, m, "-", co)
}
func (m *LaTeXConv) Mult(fml Formula, co *cnv_out) {
	mop(fml, m, " ", co)
}
func (m *LaTeXConv) Div(fml Formula, co *cnv_out) {
	mop(fml, m, "/", co)
}
func (m *LaTeXConv) Pow(fml Formula, co *cnv_out) {
	for i := 0; i < len(fml.args); i++ {
		if i != 0 {
			co.append("^{")
			conv2(fml.args[i], m, co)
			co.append("}")
		} else {
			conv2(fml.args[i], m, co)
		}
	}
}

func ToLaTeX(str string) string {
	stack = new(Stack)
	l := new(SynLex)
	l.Init(strings.NewReader(str))
	yyParse(l)
	fml := tofml(stack)
	return conv(fml, new(LaTeXConv), l.comment)
}
