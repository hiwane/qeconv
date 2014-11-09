package qeconv

import (
	"strings"
)

type MathConv struct {
}

func (m *MathConv) All(f Formula, co *cnv_out) {
	prefix(f, m, "ForAll[", "]", co)
}

func (m *MathConv) Ex(f Formula, co *cnv_out) {
	prefix(f, m, "Exists[", "]", co)
}

func (m *MathConv) And(f Formula, co *cnv_out) {
	infix(f, m, " && ", co)
}

func (m *MathConv) Or(f Formula, co *cnv_out) {
	infix(f, m, " || ", co)
}

func (m *MathConv) Not(f Formula, co *cnv_out) {
	prefix(f, m, "Not[", "]", co)
}

func (m *MathConv) Impl(f Formula, co *cnv_out) {
	prefix(f, m, "Implies[", "]", co)
}

func (m *MathConv) Equiv(f Formula, co *cnv_out) {
	prefix(f, m, "Equivalent[", "]", co)
}

func (m *MathConv) Leop(f Formula, co *cnv_out) {
	infix(f, m, " <= ", co)
}

func (m *MathConv) Ltop(f Formula, co *cnv_out) {
	infix(f, m, " < ", co)
}

func (m *MathConv) Eqop(f Formula, co *cnv_out) {
	infix(f, m, " == ", co)
}

func (m *MathConv) Neop(f Formula, co *cnv_out) {
	infix(f, m, " != ", co)
}

func (m *MathConv) List(f Formula, co *cnv_out) {
	prefix(f, m, "{", "}", co)
}

func (m *MathConv) Ftrue() string {
	return "True"
}

func (m *MathConv) Ffalse() string {
	return "False"
}
func (m *MathConv) Comment(str string) string {
	return "(*" + str + "*)"
}


func ToMath(str string) string {
	stack = new(Stack)
	l := new(SynLex)
	l.Init(strings.NewReader(str))
	yyParse(l)
	fml := tofml(stack)
	return conv(fml, new(MathConv), l.comment)
}
