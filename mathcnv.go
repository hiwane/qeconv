package qeconv

type mathConv struct {
}

func (m *mathConv) All(f Formula, co *cnv_out) {
	prefix(f, m, "ForAll[", "]", co)
}

func (m *mathConv) Ex(f Formula, co *cnv_out) {
	prefix(f, m, "Exists[", "]", co)
}

func (m *mathConv) And(f Formula, co *cnv_out) {
	infix(f, m, " && ", co)
}

func (m *mathConv) Or(f Formula, co *cnv_out) {
	infix(f, m, " || ", co)
}

func (m *mathConv) Not(f Formula, co *cnv_out) {
	prefix(f, m, "Not[", "]", co)
}

func (m *mathConv) Impl(f Formula, co *cnv_out) {
	prefix(f, m, "Implies[", "]", co)
}

func (m *mathConv) Equiv(f Formula, co *cnv_out) {
	prefix(f, m, "Equivalent[", "]", co)
}

func (m *mathConv) Abs(f Formula, co *cnv_out) {
	co.append("Abs[")
	convm(f.args[0], m, co)
	co.append("]")
}

func (m *mathConv) Leop(f Formula, co *cnv_out) {
	infix(f, m, " <= ", co)
}

func (m *mathConv) Ltop(f Formula, co *cnv_out) {
	infix(f, m, " < ", co)
}

func (m *mathConv) Eqop(f Formula, co *cnv_out) {
	infix(f, m, " == ", co)
}

func (m *mathConv) Neop(f Formula, co *cnv_out) {
	infix(f, m, " != ", co)
}

func (m *mathConv) List(f Formula, co *cnv_out) {
	prefix(f, m, "{", "}", co)
}

func (m *mathConv) Plus(fml Formula, co *cnv_out) {
	mop(fml, m, "+", co)
}
func (m *mathConv) Minus(fml Formula, co *cnv_out) {
	mop(fml, m, "-", co)
}
func (m *mathConv) Mult(fml Formula, co *cnv_out) {
	mop(fml, m, "*", co)
}
func (m *mathConv) Div(fml Formula, co *cnv_out) {
	mop(fml, m, "/", co)
}
func (m *mathConv) Pow(fml Formula, co *cnv_out) {
	mop(fml, m, "^", co)
}

func (m *mathConv) uniop(fml Formula, ope string, co *cnv_out) {
	uniop(fml, m, ope, co)
}

func (m *mathConv) Ftrue() string {
	return "True"
}

func (m *mathConv) Ffalse() string {
	return "False"
}
func (m *mathConv) Comment(str string) string {
	return "(*" + str + "*)"
}

func (m *mathConv) Convert(fml Formula, co *cnv_out) (string, error) {
	conv2(fml, m, co)
	return co.str, nil
}

func (m *mathConv) Sep() string {
	return ";"
}
