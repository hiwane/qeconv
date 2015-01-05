package qeconv

import (
	"errors"
	"strings"
)

type regchainConv struct {
	err error
}

func (m *regchainConv) quantifier(f Formula, co *cnv_out, qstr string) {
	co.append("`&" + qstr + "`(")
	conv2(f.args[0], m, co)
	co.append("), ")
	conv2(f.args[1], m, co)
}

func (m *regchainConv) All(f Formula, co *cnv_out) {
	m.quantifier(f, co, "A")
}

func (m *regchainConv) Ex(f Formula, co *cnv_out) {
	m.quantifier(f, co, "E")
}

func (m *regchainConv) And(f Formula, co *cnv_out) {
	prefix(f, m, "`&and`(", ")", co)
}

func (m *regchainConv) Or(f Formula, co *cnv_out) {
	prefix(f, m, "`&or`(", ")", co)
}

func (m *regchainConv) Not(f Formula, co *cnv_out) {
	prefix(f, m, "`&not`(", ")", co)
}

func (m *regchainConv) Impl(f Formula, co *cnv_out) {
	prefix(f, m, "`&implies`(", ")", co)
}

func (m *regchainConv) Equiv(f Formula, co *cnv_out) {
	prefix(f, m, "`&iff`(", ")", co)
}

func (m *regchainConv) Abs(f Formula, co *cnv_out) {
	// unsupported
	m.err = errors.New("unsupport the abs function")
}

func (m *regchainConv) Leop(f Formula, co *cnv_out) {
	infix(f, m, " <= ", co)
}

func (m *regchainConv) Ltop(f Formula, co *cnv_out) {
	infix(f, m, " < ", co)
}

func (m *regchainConv) Eqop(f Formula, co *cnv_out) {
	infix(f, m, " = ", co)
}

func (m *regchainConv) Neop(f Formula, co *cnv_out) {
	infix(f, m, " <> ", co)
}

func (m *regchainConv) List(f Formula, co *cnv_out) {
	prefix(f, m, "[", "]", co)
}

func (m *regchainConv) Plus(fml Formula, co *cnv_out) {
	mop(fml, m, "+", co)
}
func (m *regchainConv) Minus(fml Formula, co *cnv_out) {
	mop(fml, m, "-", co)
}
func (m *regchainConv) Mult(fml Formula, co *cnv_out) {
	mop(fml, m, "*", co)
}
func (m *regchainConv) Div(fml Formula, co *cnv_out) {
	mop(fml, m, "/", co)
}
func (m *regchainConv) Pow(fml Formula, co *cnv_out) {
	mop(fml, m, "^", co)
}

func (m *regchainConv) Ftrue() string {
	return "true"
}

func (m *regchainConv) Ffalse() string {
	return "false"
}
func (m *regchainConv) Comment(str string) string {
	return "#" + str
}

func ToRegularChains(str string) (string, error) {
	stack = new(Stack)
	l := new(SynLex)
	l.Init(strings.NewReader(str))
	yyParse(l)
	fml := tofml(stack)
	qc := new(regchainConv)
	qc.err = nil
	qstr := conv(fml, qc, l.comment)
	return qstr, qc.err
}
