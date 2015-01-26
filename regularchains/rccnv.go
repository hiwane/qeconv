package qeconv

import (
	. "github.com/hiwane/qeconv/def"
)

import (
	"errors"
)

type RegChainConv struct {
	*CnvInfStrstruct
	err error
}

func (m *RegChainConv) quantifier(f Formula, co *CnvOut, qstr string) {
	co.Append("`&" + qstr + "`(")
	Conv2(f.Args()[0], m, co)
	co.Append("), ")
	Conv2(f.Args()[1], m, co)
}

func (m *RegChainConv) All(f Formula, co *CnvOut) {
	m.quantifier(f, co, "A")
}

func (m *RegChainConv) Ex(f Formula, co *CnvOut) {
	m.quantifier(f, co, "E")
}

func (m *RegChainConv) And(f Formula, co *CnvOut) {
	Prefix(f, m, "`&and`(", ")", co)
}

func (m *RegChainConv) Or(f Formula, co *CnvOut) {
	Prefix(f, m, "`&or`(", ")", co)
}

func (m *RegChainConv) Not(f Formula, co *CnvOut) {
	Prefix(f, m, "`&not`(", ")", co)
}

func (m *RegChainConv) Impl(f Formula, co *CnvOut) {
	Prefix(f, m, "`&implies`(", ")", co)
}

func (m *RegChainConv) Equiv(f Formula, co *CnvOut) {
	Prefix(f, m, "`&iff`(", ")", co)
}

func (m *RegChainConv) Abs(f Formula, co *CnvOut) {
	// unsupported
	m.err = errors.New("unsupport the abs function")
}

func (m *RegChainConv) Leop(f Formula, co *CnvOut) {
	Infix(f, m, " <= ", co)
}

func (m *RegChainConv) Ltop(f Formula, co *CnvOut) {
	Infix(f, m, " < ", co)
}

func (m *RegChainConv) Eqop(f Formula, co *CnvOut) {
	Infix(f, m, " = ", co)
}

func (m *RegChainConv) Neop(f Formula, co *CnvOut) {
	Infix(f, m, " <> ", co)
}

func (m *RegChainConv) List(f Formula, co *CnvOut) {
	Prefix(f, m, "[", "]", co)
}

func (m *RegChainConv) Plus(fml Formula, co *CnvOut) {
	Mop(fml, m, "+", co)
}
func (m *RegChainConv) Minus(fml Formula, co *CnvOut) {
	Mop(fml, m, "-", co)
}
func (m *RegChainConv) Mult(fml Formula, co *CnvOut) {
	Mop(fml, m, "*", co)
}
func (m *RegChainConv) Div(fml Formula, co *CnvOut) {
	Mop(fml, m, "/", co)
}
func (m *RegChainConv) Pow(fml Formula, co *CnvOut) {
	Mop(fml, m, "^", co)
}
func (m *RegChainConv) Uniop(fml Formula, ope string, co *CnvOut) {
	Uniop(fml, m, ope, co)
}

func (m *RegChainConv) Ftrue() string {
	return "true"
}

func (m *RegChainConv) Ffalse() string {
	return "false"
}
func (m *RegChainConv) Comment(str string) string {
	return "#" + str
}

func (m *RegChainConv) Convert(fml Formula, co *CnvOut) (string, error) {
	m.err = nil
	Conv2(fml, m, co)
	return co.String(), m.err
}

func (m *RegChainConv) Sep() string {
	return ":"
}
