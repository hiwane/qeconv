package qeconv

import (
	. "github.com/hiwane/qeconv/def"
)

type LatexConv struct {
	*CnvInfStrstruct
}

func (m *LatexConv) q(f Formula, qop string, co *CnvOut) {
	vars := f.Arg(0)
	for i := 0; i < len(vars.Args()); i++ {
		co.Append(" " + qop + " ")
		Conv2(vars.Arg(i), m, co)
	}

	q := f.Arg(1)
	if q.IsQuantifier() {
		Conv2(q, m, co)
	} else {
		co.Append("(")
		Conv2(q, m, co)
		co.Append(")")
	}
}

func (m *LatexConv) All(f Formula, co *CnvOut) {
	m.q(f, "\\forall", co)
}

func (m *LatexConv) Ex(f Formula, co *CnvOut) {
	m.q(f, "\\exists", co)
}

func (m *LatexConv) And(f Formula, co *CnvOut) {
	Infix(f, m, " \\land ", co)
}

func (m *LatexConv) Or(f Formula, co *CnvOut) {
	Infix(f, m, " \\lor ", co)
}

func (m *LatexConv) Not(f Formula, co *CnvOut) {
	Prefix(f, m, "\\neg (", ")", co)
}

func (m *LatexConv) Impl(f Formula, co *CnvOut) {
	Infix(f, m, " \\rightarrow ", co)
}

func (m *LatexConv) Equiv(f Formula, co *CnvOut) {
	Infix(f, m, " \\leftrightarrow ", co)
}

func (m *LatexConv) Abs(f Formula, co *CnvOut) {
	Prefix(f, m, "|", "|", co)
}

func (m *LatexConv) Leop(f Formula, co *CnvOut) {
	Infix(f, m, " \\leq ", co)
}

func (m *LatexConv) Ltop(f Formula, co *CnvOut) {
	Infix(f, m, " < ", co)
}

func (m *LatexConv) Eqop(f Formula, co *CnvOut) {
	Infix(f, m, " = ", co)
}

func (m *LatexConv) Neop(f Formula, co *CnvOut) {
	Infix(f, m, " \\neq ", co)
}

func (m *LatexConv) List(f Formula, co *CnvOut) {
	Prefix(f, m, "[", "]", co)
}

func (m *LatexConv) Ftrue() string {
	return "\\top"
}

func (m *LatexConv) Ffalse() string {
	return "\\bot"
}
func (m *LatexConv) Comment(str string) string {
	return "% " + str
}

func (m *LatexConv) Plus(fml Formula, co *CnvOut) {
	Mop(fml, m, "+", co)
}
func (m *LatexConv) Minus(fml Formula, co *CnvOut) {
	Mop(fml, m, "-", co)
}
func (m *LatexConv) Mult(fml Formula, co *CnvOut) {
	Mop(fml, m, " ", co)
}
func (m *LatexConv) Div(fml Formula, co *CnvOut) {
	Mop(fml, m, "/", co)
}
func (m *LatexConv) Pow(fml Formula, co *CnvOut) {
	for i := 0; i < len(fml.Args()); i++ {
		if i != 0 {
			co.Append("^{")
			Conv2(fml.Arg(i), m, co)
			co.Append("}")
		} else {
			Conv2(fml.Arg(i), m, co)
		}
	}
}

func (m *LatexConv) Uniop(fml Formula, ope string, co *CnvOut) {
	Uniop(fml, m, ope, co)
}

func (m *LatexConv) Convert(fml Formula, co *CnvOut) (string, error) {
	Conv2(fml, m, co)
	return co.String(), nil
}

func (m *LatexConv) Sep() string {
	return "\\\\"
}
