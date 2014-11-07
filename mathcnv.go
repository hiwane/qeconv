package synparse

type MathConv struct {
}

func (m *MathConv) All(f Formula) string {
	return prefix(f, m, "ForAll[", "]")
}

func (m *MathConv) Ex(f Formula) string {
	return prefix(f, m, "Exists[", "]")
}

func (m *MathConv) And(f Formula) string {
	return infix(f, m, " && ")
}

func (m *MathConv) Or(f Formula) string {
	return infix(f, m, " || ")
}

func (m *MathConv) Not(f Formula) string {
	return prefix(f, m, "Not[", "]")
}

func (m *MathConv) Impl(f Formula) string {
	return prefix(f, m, "Implies[", "]")
}

func (m *MathConv) Equiv(f Formula) string {
	return prefix(f, m, "Equivalent[", "]")
}

func (m *MathConv) Leop(f Formula) string {
	return infix(f, m, " <= ")
}

func (m *MathConv) Ltop(f Formula) string {
	return infix(f, m, " < ")
}

func (m *MathConv) Eqop(f Formula) string {
	return infix(f, m, " == ")
}

func (m *MathConv) Neop(f Formula) string {
	return infix(f, m, " != ")
}

func (m *MathConv) List(f Formula) string {
	return prefix(f, m, "{", "}")
}

func (m *MathConv) Ftrue() string {
	return "True"
}

func (m *MathConv) Ffalse() string {
	return "False"
}

func ToMath(l *SynLex) string {
	stack = new(Stack)
	yyParse(l)
	fml := tofml(stack)
	trace("go conv")
	str := conv(fml, new(MathConv))
	trace("end: " + str)
	return str
}

