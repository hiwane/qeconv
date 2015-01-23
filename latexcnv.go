package qeconv

type latexConv struct {
}

func (m *latexConv) q(f Formula, qop string, co *cnv_out) {
	vars := f.args[0]
	for i := 0; i < len(vars.args); i++ {
		co.append(" " + qop + " ")
		conv2(vars.args[i], m, co)
	}

	if f.args[1].cmd == ALL || f.args[1].cmd == EX {
		conv2(f.args[1], m, co)
	} else {
		co.append("(")
		conv2(f.args[1], m, co)
		co.append(")")
	}
}

func (m *latexConv) All(f Formula, co *cnv_out) {
	m.q(f, "\\forall", co)
}

func (m *latexConv) Ex(f Formula, co *cnv_out) {
	m.q(f, "\\exists", co)
}

func (m *latexConv) And(f Formula, co *cnv_out) {
	infix(f, m, " \\land ", co)
}

func (m *latexConv) Or(f Formula, co *cnv_out) {
	infix(f, m, " \\lor ", co)
}

func (m *latexConv) Not(f Formula, co *cnv_out) {
	prefix(f, m, "\\neg (", ")", co)
}

func (m *latexConv) Impl(f Formula, co *cnv_out) {
	infix(f, m, " \\rightarrow ", co)
}

func (m *latexConv) Equiv(f Formula, co *cnv_out) {
	infix(f, m, " \\leftrightarrow ", co)
}

func (m *latexConv) Abs(f Formula, co *cnv_out) {
	prefix(f, m, "|", "|", co)
}

func (m *latexConv) Leop(f Formula, co *cnv_out) {
	infix(f, m, " \\leq ", co)
}

func (m *latexConv) Ltop(f Formula, co *cnv_out) {
	infix(f, m, " < ", co)
}

func (m *latexConv) Eqop(f Formula, co *cnv_out) {
	infix(f, m, " = ", co)
}

func (m *latexConv) Neop(f Formula, co *cnv_out) {
	infix(f, m, " \\neq ", co)
}

func (m *latexConv) List(f Formula, co *cnv_out) {
	prefix(f, m, "[", "]", co)
}

func (m *latexConv) Ftrue() string {
	return "\\top"
}

func (m *latexConv) Ffalse() string {
	return "\\bot"
}
func (m *latexConv) Comment(str string) string {
	return "% " + str
}

func (m *latexConv) Plus(fml Formula, co *cnv_out) {
	mop(fml, m, "+", co)
}
func (m *latexConv) Minus(fml Formula, co *cnv_out) {
	mop(fml, m, "-", co)
}
func (m *latexConv) Mult(fml Formula, co *cnv_out) {
	mop(fml, m, " ", co)
}
func (m *latexConv) Div(fml Formula, co *cnv_out) {
	mop(fml, m, "/", co)
}
func (m *latexConv) Pow(fml Formula, co *cnv_out) {
	for i := 0; i < len(fml.args); i++ {
		if i != 0 {
			co.append("^{")
			conv2(fml.args[i], m, co)
			co.append("}")
		} else {
			conv2(fml.args[i], m, co)
		}
	}
}

func (m *latexConv) uniop(fml Formula, ope string, co *cnv_out) {
	uniop(fml, m, ope, co)
}

func (m *latexConv) Convert(fml Formula, co *cnv_out) (string, error) {
	conv2(fml, m, co)
	return co.str, nil
}

func (m *latexConv) Sep() string {
	return "\\\\"
}
