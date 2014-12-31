package qeconv

import (
	"strings"
)

type redConv struct {
	err error
}

func (m *redConv) All(f Formula, co *cnv_out) {
	prefix(f, m, "all(", ")", co)
}

func (m *redConv) Ex(f Formula, co *cnv_out) {
	prefix(f, m, "ex(", ")", co)
}

func (m *redConv) And(f Formula, co *cnv_out) {
	infix(f, m, " and ", co)
}

func (m *redConv) Or(f Formula, co *cnv_out) {
	infix(f, m, " or ", co)
}

func (m *redConv) Not(f Formula, co *cnv_out) {
	prefix(f, m, "not (", ")", co)
}

func (m *redConv) Impl(f Formula, co *cnv_out) {
	infix(f, m, " implies ", co)
}

func (m *redConv) Equiv(f Formula, co *cnv_out) {
	infix(f, m, " equiv ", co)
}

func (m *redConv) Abs(f Formula, co *cnv_out) {
	co.append("abs(")
	convm(f.args[0], m, co)
	co.append(")")
}

func (m *redConv) Leop(f Formula, co *cnv_out) {
	infix(f, m, " <= ", co)
}

func (m *redConv) Ltop(f Formula, co *cnv_out) {
	infix(f, m, " < ", co)
}

func (m *redConv) Eqop(f Formula, co *cnv_out) {
	infix(f, m, " = ", co)
}

func (m *redConv) Neop(f Formula, co *cnv_out) {
	infix(f, m, " <> ", co)
}

func (m *redConv) List(f Formula, co *cnv_out) {
	prefix(f, m, "[", "]", co)
}

func (m *redConv) Plus(fml Formula, co *cnv_out) {
	mop(fml, m, "+", co)
}
func (m *redConv) Minus(fml Formula, co *cnv_out) {
	mop(fml, m, "-", co)
}
func (m *redConv) Mult(fml Formula, co *cnv_out) {
	mop(fml, m, "*", co)
}
func (m *redConv) Div(fml Formula, co *cnv_out) {
	mop(fml, m, "/", co)
}
func (m *redConv) Pow(fml Formula, co *cnv_out) {
	mop(fml, m, "^", co)
}

func (m *redConv) Ftrue() string {
	return "true"
}

func (m *redConv) Ffalse() string {
	return "false"
}
func (m *redConv) Comment(str string) string {
	return "%" + str
}

func ToRedlog(str string) string {
	stack = new(Stack)
	l := new(SynLex)
	l.Init(strings.NewReader(str))
	yyParse(l)
	fml := tofml(stack)
	qc := new(redConv)
	qstr := conv(fml, qc, l.comment)
	return qstr
}
