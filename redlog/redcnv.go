package qeconv

import (
	. "github.com/hiwane/qeconv/def"
)

type RedConv struct {
}

func (m *RedConv) All(f Formula, co *CnvOut) {
	Prefix(f, m, "all(", ")", co)
}

func (m *RedConv) Ex(f Formula, co *CnvOut) {
	Prefix(f, m, "ex(", ")", co)
}

func (m *RedConv) And(f Formula, co *CnvOut) {
	Infix(f, m, " and ", co)
}

func (m *RedConv) Or(f Formula, co *CnvOut) {
	Infix(f, m, " or ", co)
}

func (m *RedConv) Not(f Formula, co *CnvOut) {
	Prefix(f, m, "not (", ")", co)
}

func (m *RedConv) Impl(f Formula, co *CnvOut) {
	Infix(f, m, " implies ", co)
}

func (m *RedConv) Equiv(f Formula, co *CnvOut) {
	Infix(f, m, " equiv ", co)
}

func (m *RedConv) Abs(f Formula, co *CnvOut) {
	co.Append("abs(")
	Convm(f.Arg(0), m, co)
	co.Append(")")
}

func (m *RedConv) Leop(f Formula, co *CnvOut) {
	Infix(f, m, " <= ", co)
}

func (m *RedConv) Ltop(f Formula, co *CnvOut) {
	Infix(f, m, " < ", co)
}

func (m *RedConv) Eqop(f Formula, co *CnvOut) {
	Infix(f, m, " = ", co)
}

func (m *RedConv) Neop(f Formula, co *CnvOut) {
	Infix(f, m, " <> ", co)
}

func (m *RedConv) List(f Formula, co *CnvOut) {
	Prefix(f, m, "[", "]", co)
}

func (m *RedConv) Plus(fml Formula, co *CnvOut) {
	Mop(fml, m, "+", co)
}
func (m *RedConv) Minus(fml Formula, co *CnvOut) {
	Mop(fml, m, "-", co)
}
func (m *RedConv) Mult(fml Formula, co *CnvOut) {
	Mop(fml, m, "*", co)
}
func (m *RedConv) Div(fml Formula, co *CnvOut) {
	Mop(fml, m, "/", co)
}
func (m *RedConv) Pow(fml Formula, co *CnvOut) {
	Mop(fml, m, "^", co)
}

func (m *RedConv) Uniop(fml Formula, ope string, co *CnvOut) {
	Uniop(fml, m, ope, co)
}

func (m *RedConv) Ftrue() string {
	return "true"
}

func (m *RedConv) Ffalse() string {
	return "false"
}
func (m *RedConv) Comment(str string) string {
	return "%" + str
}

func (m *RedConv) Convert(fml Formula, co *CnvOut) (string, error) {
	Conv2(fml, m, co)
	return co.String(), nil
}

func (m *RedConv) Sep() string {
	return ";"
}
