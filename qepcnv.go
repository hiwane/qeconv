package qeconv

import (
	"errors"
	"strings"
)

type qepConv struct {
	err error
}

func (m *qepConv) quantifier(f Formula, co *cnv_out, qstr string) {
	q := f.args[0]
	for i := 0; i < len(q.args); i++ {
		co.append("(" + qstr + " ")
		conv2(q.args[i], m, co)
		co.append(")")
	}

	q = f.args[1]
	if q.cmd != ALL && q.cmd != EX {
		co.append("[")
		conv2(q, m, co)
		co.append("]")
	} else {
		conv2(q, m, co)
	}
}

func (m *qepConv) All(f Formula, co *cnv_out) {
	m.quantifier(f, co, "A")
}

func (m *qepConv) Ex(f Formula, co *cnv_out) {
	m.quantifier(f, co, "E")
}

func (m *qepConv) And(f Formula, co *cnv_out) {
	infixm(f, m, " /\\ ", co, "[", "]")
}

func (m *qepConv) Or(f Formula, co *cnv_out) {
	infixm(f, m, " \\/ ", co, "[", "]")
}

func (m *qepConv) Not(f Formula, co *cnv_out) {
	prefix(f, m, "~[", "]", co)
}

func (m *qepConv) Impl(f Formula, co *cnv_out) {
	prefixm(f, m, "[[", "] ==> [", "]]", co)
}

func (m *qepConv) Equiv(f Formula, co *cnv_out) {
	prefixm(f, m, "[", " <==> ", "]", co)
}

func (m *qepConv) Abs(f Formula, co *cnv_out) {
	// unsupported
	m.err = errors.New("unsupport the abs function")
}

func (m *qepConv) Leop(f Formula, co *cnv_out) {
	infix(f, m, " <= ", co)
}

func (m *qepConv) Ltop(f Formula, co *cnv_out) {
	infix(f, m, " < ", co)
}

func (m *qepConv) Eqop(f Formula, co *cnv_out) {
	infix(f, m, " = ", co)
}

func (m *qepConv) Neop(f Formula, co *cnv_out) {
	infix(f, m, " /= ", co)
}

func (m *qepConv) List(f Formula, co *cnv_out) {
	prefix(f, m, "[", "]", co)
}

func (m *qepConv) Plus(fml Formula, co *cnv_out) {
	mop(fml, m, "+", co)
}
func (m *qepConv) Minus(fml Formula, co *cnv_out) {
	mop(fml, m, "-", co)
}
func (m *qepConv) Mult(fml Formula, co *cnv_out) {
	mop(fml, m, " ", co)
}
func (m *qepConv) Div(fml Formula, co *cnv_out) {
	// unsupported
	m.err = errors.New("unsupport /")
}
func (m *qepConv) Pow(fml Formula, co *cnv_out) {
	mop(fml, m, "^", co)
}

func (m *qepConv) uniop(fml Formula, ope string, co *cnv_out) {
	uniop(fml, m, ope, co)
}


func (m *qepConv) Ftrue() string {
	return "true"
}

func (m *qepConv) Ffalse() string {
	return "false"
}
func (m *qepConv) Comment(str string) string {
	return ""
}

func ToQepcad(str string) (string, error) {
	stack = new(Stack)
	l := new(SynLex)
	l.Init(strings.NewReader(str))
	yyParse(l)
	fml := tofml(stack)
	qc := new(qepConv)
	qc.err = nil
	qstr := conv(fml, qc, l.comment)
	return qstr, qc.err
}
