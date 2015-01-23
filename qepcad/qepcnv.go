package qeconv

import (
	"errors"
	. "github.com/hiwane/qeconv/def"
)

type QepConv struct {
	err error
}

func (m *QepConv) quantifier(f Formula, co *CnvOut, qstr string) {
	q := f.Arg(0)
	for i := 0; i < len(q.Args()); i++ {
		co.Append("(" + qstr + " ")
		Conv2(q.Arg(i), m, co)
		co.Append(")")
	}

	q = f.Arg(1)
	if q.IsQuantifier() {
		co.Append("[")
		Conv2(q, m, co)
		co.Append("]")
	} else {
		Conv2(q, m, co)
	}
}

func (m *QepConv) All(f Formula, co *CnvOut) {
	m.quantifier(f, co, "A")
}

func (m *QepConv) Ex(f Formula, co *CnvOut) {
	m.quantifier(f, co, "E")
}

func (m *QepConv) And(f Formula, co *CnvOut) {
	Infixm(f, m, " /\\ ", co, "[", "]")
}

func (m *QepConv) Or(f Formula, co *CnvOut) {
	Infixm(f, m, " \\/ ", co, "[", "]")
}

func (m *QepConv) Not(f Formula, co *CnvOut) {
	Prefix(f, m, "~[", "]", co)
}

func (m *QepConv) Impl(f Formula, co *CnvOut) {
	Prefixm(f, m, "[[", "] ==> [", "]]", co)
}

func (m *QepConv) Equiv(f Formula, co *CnvOut) {
	Prefixm(f, m, "[", " <==> ", "]", co)
}

func (m *QepConv) Abs(f Formula, co *CnvOut) {
	// unsupported
	m.err = errors.New("unsupport the abs function")
}

func (m *QepConv) Leop(f Formula, co *CnvOut) {
	Infix(f, m, " <= ", co)
}

func (m *QepConv) Ltop(f Formula, co *CnvOut) {
	Infix(f, m, " < ", co)
}

func (m *QepConv) Eqop(f Formula, co *CnvOut) {
	Infix(f, m, " = ", co)
}

func (m *QepConv) Neop(f Formula, co *CnvOut) {
	Infix(f, m, " /= ", co)
}

func (m *QepConv) List(f Formula, co *CnvOut) {
	Prefix(f, m, "[", "]", co)
}

func (m *QepConv) Plus(fml Formula, co *CnvOut) {
	Mop(fml, m, "+", co)
}
func (m *QepConv) Minus(fml Formula, co *CnvOut) {
	Mop(fml, m, "-", co)
}
func (m *QepConv) Mult(fml Formula, co *CnvOut) {
	Mop(fml, m, " ", co)
}
func (m *QepConv) Div(fml Formula, co *CnvOut) {
	// unsupported
	m.err = errors.New("unsupport /")
}
func (m *QepConv) Pow(fml Formula, co *CnvOut) {
	Mop(fml, m, "^", co)
}

func (m *QepConv) Uniop(fml Formula, ope string, co *CnvOut) {
	Uniop(fml, m, ope, co)
}

func (m *QepConv) Ftrue() string {
	return "true"
}

func (m *QepConv) Ffalse() string {
	return "false"
}
func (m *QepConv) Comment(str string) string {
	return ""
}

func (m *QepConv) Convert(fml Formula, co *CnvOut) (string, error) {
	m.err = nil
	Conv2(fml, m, co)
	return co.String(), m.err
}

func (m *QepConv) Sep() string {
	return ".\n"
}
