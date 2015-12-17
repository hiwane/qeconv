package qeconv

import (
	. "github.com/hiwane/qeconv/def"
)

type SynParse struct {
}

func NewSynParse() *SynParse {
	p := new(SynParse)
	return p
}


func (m *SynParse) charIndex(str string, sep uint8) int {
	comment := false
	for i := 0; i < len(str); i++ {
		if str[i] == '#' {
			comment = true
		} else if str[i] == '\n' {
			comment = false
		} else if str[i] == sep && !comment {
			return i
		}
	}

	return -1
}

func (self *SynParse) Next(str string) int {
	idx := self.charIndex(str, ':')
	if idx >= 0 {
		idx += 1
	}
	return idx
}

func (self *SynParse) Parse(str string) (Formula, []Comment, error) {
	stack, c, e := parse(str)
	return ToFml(stack), c, e
}


type SynConv struct {
	*CnvInfStrstruct
}

func (m *SynConv) All(f Formula, co *CnvOut) {
	f0 := f.Arg(0)
	if f0.IsList() {
		Prefix(f, m, "All(", ")", co)
	} else {
		Prefixm(f, m, "All([", "], ", ")", co)
	}
}

func (m *SynConv) Ex(f Formula, co *CnvOut) {
	f0 := f.Arg(0)
	if f0.IsList() {
		Prefix(f, m, "Ex(", ")", co)
	} else {
		Prefixm(f, m, "Ex([", "], ", ")", co)
	}
}

func (m *SynConv) And(f Formula, co *CnvOut) {
	Prefix(f, m, "And(", ")", co)
}

func (m *SynConv) Or(f Formula, co *CnvOut) {
	Prefix(f, m, "Or(", ")", co)
}

func (m *SynConv) Not(f Formula, co *CnvOut) {
	Prefix(f, m, "Not(", ")", co)
}

func (m *SynConv) Impl(f Formula, co *CnvOut) {
	Prefix(f, m, "Impl(", ")", co)
}

func (m *SynConv) Equiv(f Formula, co *CnvOut) {
	Prefix(f, m, "Equiv(", ")", co)
}

func (m *SynConv) Abs(f Formula, co *CnvOut) {
	co.Append("abs(")
	Convm(f.Arg(0), m, co)
	co.Append(")")
}

func (m *SynConv) Leop(f Formula, co *CnvOut) {
	Infix(f, m, " <= ", co)
}

func (m *SynConv) Ltop(f Formula, co *CnvOut) {
	Infix(f, m, " < ", co)
}

func (m *SynConv) Eqop(f Formula, co *CnvOut) {
	Infix(f, m, " = ", co)
}

func (m *SynConv) Neop(f Formula, co *CnvOut) {
	Infix(f, m, " <> ", co)
}

func (m *SynConv) List(f Formula, co *CnvOut) {
	Prefix(f, m, "[", "]", co)
}

func (m *SynConv) Plus(fml Formula, co *CnvOut) {
	Mop(fml, m, "+", co)
}
func (m *SynConv) Minus(fml Formula, co *CnvOut) {
	Mop(fml, m, "-", co)
}
func (m *SynConv) Mult(fml Formula, co *CnvOut) {
	Mop(fml, m, "*", co)
}
func (m *SynConv) Div(fml Formula, co *CnvOut) {
	Mop(fml, m, "/", co)
}
func (m *SynConv) Pow(fml Formula, co *CnvOut) {
	Mop(fml, m, "^", co)
}

func (m *SynConv) Uniop(fml Formula, ope string, co *CnvOut) {
	Uniop(fml, m, ope, co)
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

func (m *SynConv) Convert(fml Formula, co *CnvOut) (string, error) {
	Conv2(fml, m, co)
	return co.String(), nil
}

func (m *SynConv) Sep() string {
	return ":"
}
