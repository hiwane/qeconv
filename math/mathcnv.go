package qeconv

import (
	. "github.com/hiwane/qeconv/def"
)

type MathConv struct {
}

func (m *MathConv) All(f Formula, co *CnvOut) {
	Prefix(f, m, "ForAll[", "]", co)
}

func (m *MathConv) Ex(f Formula, co *CnvOut) {
	Prefix(f, m, "Exists[", "]", co)
}

func (m *MathConv) And(f Formula, co *CnvOut) {
	Infix(f, m, " && ", co)
}

func (m *MathConv) Or(f Formula, co *CnvOut) {
	Infix(f, m, " || ", co)
}

func (m *MathConv) Not(f Formula, co *CnvOut) {
	Prefix(f, m, "Not[", "]", co)
}

func (m *MathConv) Impl(f Formula, co *CnvOut) {
	Prefix(f, m, "Implies[", "]", co)
}

func (m *MathConv) Equiv(f Formula, co *CnvOut) {
	Prefix(f, m, "Equivalent[", "]", co)
}

func (m *MathConv) Abs(f Formula, co *CnvOut) {
	co.Append("Abs[")
	Convm(f.Arg(0), m, co)
	co.Append("]")
}

func (m *MathConv) Leop(f Formula, co *CnvOut) {
	Infix(f, m, " <= ", co)
}

func (m *MathConv) Ltop(f Formula, co *CnvOut) {
	Infix(f, m, " < ", co)
}

func (m *MathConv) Eqop(f Formula, co *CnvOut) {
	Infix(f, m, " == ", co)
}

func (m *MathConv) Neop(f Formula, co *CnvOut) {
	Infix(f, m, " != ", co)
}

func (m *MathConv) List(f Formula, co *CnvOut) {
	Prefix(f, m, "{", "}", co)
}

func (m *MathConv) Plus(fml Formula, co *CnvOut) {
	Mop(fml, m, "+", co)
}
func (m *MathConv) Minus(fml Formula, co *CnvOut) {
	Mop(fml, m, "-", co)
}
func (m *MathConv) Mult(fml Formula, co *CnvOut) {
	Mop(fml, m, "*", co)
}
func (m *MathConv) Div(fml Formula, co *CnvOut) {
	Mop(fml, m, "/", co)
}
func (m *MathConv) Pow(fml Formula, co *CnvOut) {
	Mop(fml, m, "^", co)
}

func (m *MathConv) Uniop(fml Formula, ope string, co *CnvOut) {
	Uniop(fml, m, ope, co)
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

func (m *MathConv) Convert(fml Formula, co *CnvOut) (string, error) {
	Conv2(fml, m, co)
	return co.String(), nil
}

func (m *MathConv) Sep() string {
	return ";"
}
