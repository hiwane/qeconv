package qeconv

import (
	"errors"
	. "github.com/hiwane/qeconv/def"
	"strconv"
)

type Smt2Parse struct {
}

func NewSmt2Parse() *Smt2Parse {
	p := new(Smt2Parse)
	return p
}

func (self *Smt2Parse) Next(str string) int {
	if str == "" {
		return -1
	}
	return len(str)
}

func (self *Smt2Parse) Parse(str string) (Formula, []Comment, error) {
	stack, c, e := parse(str)
	return ToFml(stack), c, e
}

type Smt2Conv struct {
	*CnvInfStrstruct
	err error
}

func (m *Smt2Conv) quantifier(f Formula, co *CnvOut, qstr string) {
	q := f.Args()[0]
	co.Append("(" + qstr + " (")
	for i := 0; i < len(q.Args()); i++ {
		co.Append(" (")
		Conv2(q.Arg(i), m, co)
		co.Append(" Real)")
	}
	co.Append(" ) ")
	Conv2(f.Args()[1], m, co)
	co.Append(" )")
}

func (m *Smt2Conv) All(f Formula, co *CnvOut) {
	m.quantifier(f, co, "forall")
}

func (m *Smt2Conv) Ex(f Formula, co *CnvOut) {
	m.quantifier(f, co, "exists")
}

func (m *Smt2Conv) And(f Formula, co *CnvOut) {
	Prefixm(f, m, "(and ", " ", ")", co)
}

func (m *Smt2Conv) Or(f Formula, co *CnvOut) {
	Prefixm(f, m, "(or ", " ", ")", co)
}

func (m *Smt2Conv) Not(f Formula, co *CnvOut) {
	Prefixm(f, m, "(not ", " ", ")", co)
}

func (m *Smt2Conv) Impl(f Formula, co *CnvOut) {
	Prefixm(f, m, "(=> ", " ", ")", co)
}

func (m *Smt2Conv) Equiv(f Formula, co *CnvOut) {
	co.Append("(and (=> ")
	Conv2(f.Args()[0], m, co)
	co.Append(" ")
	Conv2(f.Args()[1], m, co)
	co.Append(") (=> ")
	Conv2(f.Args()[1], m, co)
	co.Append(" ")
	Conv2(f.Args()[0], m, co)
	co.Append("))")
}

func (m *Smt2Conv) Abs(f Formula, co *CnvOut) {
	m.err = errors.New("unsupport the abs function")
}

func (m *Smt2Conv) Leop(f Formula, co *CnvOut) {
	Prefixm(f, m, "(<= ", " ", ")", co)
}

func (m *Smt2Conv) Ltop(f Formula, co *CnvOut) {
	Prefixm(f, m, "(< ", " ", ")", co)
}

func (m *Smt2Conv) Eqop(f Formula, co *CnvOut) {
	Prefixm(f, m, "(= ", " ", ")", co)
}

func (m *Smt2Conv) Neop(f Formula, co *CnvOut) {
	Prefixm(f, m, "(not (= ", " ", "))", co)
}

func (m *Smt2Conv) List(f Formula, co *CnvOut) {
	Prefixm(f, m, "(", " ", ")", co)
}

func (m *Smt2Conv) Plus(fml Formula, co *CnvOut) {
	Prefixm(fml, m, "(+ ", " ", ")", co)
}

func (m *Smt2Conv) Minus(fml Formula, co *CnvOut) {
	Prefixm(fml, m, "(- ", " ", ")", co)
}

func (m *Smt2Conv) Mult(fml Formula, co *CnvOut) {
	Prefixm(fml, m, "(* ", " ", ")", co)
}

func (m *Smt2Conv) Div(fml Formula, co *CnvOut) {
	Prefixm(fml, m, "(/ ", " ", ")", co)
}

func (m *Smt2Conv) Pow(fml Formula, co *CnvOut) {
	exp := fml.Arg(1)
	if exp.Cmd() != NUMBER {
		m.err = errors.New("unsupport rational exponential")
	}
	co.Append("(*")
	n, _ := strconv.Atoi(exp.String())
	for i := 0; i < n; i++ {
		co.Append(" ")
		Conv2(fml.Args()[0], m, co)
	}
	co.Append(")")
}

func (m *Smt2Conv) Uniop(fml Formula, ope string, co *CnvOut) {
	co.Append("(" + ope + " 0 ")
	Conv2(fml.Args()[0], m, co)
	co.Append(")")
}

func (m *Smt2Conv) Ftrue() string {
	return "(= 0 0)"
}

func (m *Smt2Conv) Ffalse() string {
	return "(= 0 1)"
}

func (m *Smt2Conv) Comment(str string) string {
	return ";" + str
}

func smt2header(fml Formula) string {
	var str string
	if fml.IsQff() {
		str = "(set-logic QF_NRA)\n"
	} else {
		str = "(set-logic NRA)\n"
	}

	vs := fml.FreeVars()
	for i := 0; i < vs.Len(); i++ {
		v := vs.Get(i)
		str += "(declare-fun " + v + " () Real)\n"
	}

	return str
}

func smt2footer(fml Formula) string {
	return "(check-sat)\n"
}

func (m *Smt2Conv) Convert(fml Formula, co *CnvOut) (string, error) {
	if fml.IsList() {
		return "", errors.New("unsupported input. use -n option")
	}
	qc := new(Smt2Conv)
	qc.err = nil

	header := smt2header(fml)
	if fml.Cmd() == AND {
		for _, v := range fml.Args() {
			co.Reset()
			Conv2(v, qc, co)
			header += "(assert " + co.String() + ")\n"
		}
	} else {
		Conv2(fml, qc, co)
		header += "(assert " + co.String() + ")\n"
	}
	header += smt2footer(fml)
	return header, qc.err

}

func (m *Smt2Conv) Sep() string {
	return ""
}

func (m *Smt2Conv) LineAlign() bool {
	return false
}

type smt2letdat struct {
	name []string
	data []*QeStack
}

func (lmap *smt2letdat) update_letmap(s *QeStack, pos int, sym smt2node) {
	nstack := s.Popn(s.Length() - pos)
	lmap.name = append(lmap.name, sym.str)
	lmap.data = append(lmap.data, nstack)
}

func (lmap *smt2letdat) get(str string) (*QeStack, bool) {
	for i := len(lmap.data) - 1; i >= 0; i-- {
		if lmap.name[i] == str {
			return lmap.data[i], true
		}
	}
	return nil, false
}

func (lmap *smt2letdat) reset() {
	lmap.name = make([]string, 0)
	lmap.data = make([]*QeStack, 0)
}

func (lmap *smt2letdat) popn(n int) {
	m := len(lmap.name) - n
	lmap.name = lmap.name[:m]
	lmap.data = lmap.data[:m]
}

