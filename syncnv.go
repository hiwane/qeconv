package qeconv

import (
	"strings"
)

type SynConv struct {
}

func (m *SynConv) All(f Formula, co *cnv_out) {
	prefix(f, m, "All(", ")", co)
}

func (m *SynConv) Ex(f Formula, co *cnv_out) {
	prefix(f, m, "Ex(", ")", co)
}

func (m *SynConv) And(f Formula, co *cnv_out) {
	prefix(f, m, "All(", ")", co)
}

func (m *SynConv) Or(f Formula, co *cnv_out) {
	prefix(f, m, "Or(", ")", co)
}

func (m *SynConv) Not(f Formula, co *cnv_out) {
	prefix(f, m, "Not(", ")", co)
}

func (m *SynConv) Impl(f Formula, co *cnv_out) {
	prefix(f, m, "Impl(", ")", co)
}

func (m *SynConv) Equiv(f Formula, co *cnv_out) {
	prefix(f, m, "Equiv(", ")", co)
}

func (m *SynConv) Abs(f Formula, co *cnv_out) {
	co.append("abs(")
	convm(f.args[0], m, co)
	co.append(")")
}

func (m *SynConv) Leop(f Formula, co *cnv_out) {
	infix(f, m, " <= ", co)
}

func (m *SynConv) Ltop(f Formula, co *cnv_out) {
	infix(f, m, " < ", co)
}

func (m *SynConv) Eqop(f Formula, co *cnv_out) {
	infix(f, m, " == ", co)
}

func (m *SynConv) Neop(f Formula, co *cnv_out) {
	infix(f, m, " <> ", co)
}

func (m *SynConv) List(f Formula, co *cnv_out) {
	prefix(f, m, "[", "]", co)
}

func (m *SynConv) Plus(fml Formula, co *cnv_out) {
	mop(fml, m, "+", co)
}
func (m *SynConv) Minus(fml Formula, co *cnv_out) {
	mop(fml, m, "-", co)
}
func (m *SynConv) Mult(fml Formula, co *cnv_out) {
	mop(fml, m, "*", co)
}
func (m *SynConv) Div(fml Formula, co *cnv_out) {
	mop(fml, m, "/", co)
}
func (m *SynConv) Pow(fml Formula, co *cnv_out) {
	mop(fml, m, "^", co)
}

func (m *SynConv) Ftrue() string {
	return "true"
}

func (m *SynConv) Ffalse() string {
	return "false"
}
func (m *SynConv) Comment(str string) string {
	return "#" + str
}

func ToSyn(str string) string {
	stack = new(Stack)
	l := new(SynLex)
	l.Init(strings.NewReader(str))
	yyParse(l)
	fml := tofml(stack)
	return conv(fml, new(SynConv), l.comment)
}
