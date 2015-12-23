package qeconv

import (
	"errors"
	"fmt"
	. "github.com/hiwane/qeconv/def"
	"strconv"
	"strings"
	"text/scanner"
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

type commentI interface {
	append_comment(comm string, lno int)
}

type synLex struct {
	scanner.Scanner
	s       string
	comment []Comment
	err     error
}

type smt2_lext struct {
	val   string
	label int
}

var keywords_tbl = []smt2_lext{
	{"exists", exists},
	{"forall", forall},
	{"as", as},
	{"let", let},
	{"theory", theory},
	{"par", par},
	{"assert", assert},
	{"check-sat", check_sat},
	{"exit", exit},
	{"declare-fun", declare_fun},
	{"declare-const", declare_const},
	{"set-info", set_info},
	{"set-logic", set_logic},
	{"set-option", set_option},
	{"+", plus},
	{"-", minus},
	{"*", mult},
	{"/", div},
	{">=", geop},
	{">", gtop},
	{"<=", leop},
	{"<", ltop},
	{"=", eqop},
	{"not", not},
	{"and", and},
	{"or", or},
	{"implies", implies},
	{"=>", implies},
}

func isupper(ch rune) bool {
	return 'A' <= ch && ch <= 'Z'
}
func islower(ch rune) bool {
	return 'a' <= ch && ch <= 'z'
}
func isalpha(ch rune) bool {
	return isupper(ch) || islower(ch)
}

func issimplsym(ch rune) bool {
	return isalnum(ch) ||
		ch == '+' || ch == '-' || ch == '/' || ch == '*' ||
		ch == '=' || ch == '%' || ch == '?' || ch == '!' ||
		ch == '.' || ch == '$' || ch == '_' || ch == '~' ||
		ch == '&' || ch == '^' || ch == '<' || ch == '>'
}

func isdigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
func isalnum(ch rune) bool {
	return isalpha(ch) || isdigit(ch)
}
func isletter(ch rune) bool {
	return isalpha(ch) || ch == '_'
}
func isspace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (l *synLex) append_comment(comm string, lno int) {
	l.comment = append(l.comment, NewComment(comm, lno))
}

func (l *synLex) Lex(lval *yySymType) int {

	for {
		// skip space
		for isspace(l.Peek()) {
			l.Next()
		}
		if scanner.EOF == l.Peek() {
			trace("Lex! eof " + l.Pos().String())
			l.Next()
			return 0
		}
		c := int(l.Peek())
		if c != ';' {
			break
		}
		// comment
		l.Next()
		str := ""
		for l.Peek() != '\n' {
			str += string(l.Next())
		}
		if str != "" {
			lno := l.Pos().Line
			l.append_comment(str, lno)
		}
	}
	trace("Lex! " + string(l.Peek()) + l.Pos().String())

	c := int(l.Peek())
	lno := l.Pos().Line
	col := l.Pos().Column
	if c == '(' || c == ')' {
		l.Next()
		lval.num = stack.Length()
		if c == '(' {
			return lp
		} else {
			return rp
		}
	}

	if isdigit(l.Peek()) {
		var ret []rune
		for isdigit(l.Peek()) {
			ret = append(ret, l.Next())
		}
		if l.Peek() == '.' {
			l.Next()
			if l.Peek() == '0' {
				l.Next()
			}
			if isdigit(l.Peek()) {
				l.Error("decimal number found ")
			}
		}

		lval.node = smt2node{lno, col, number, string(ret), string(ret)}
		return number
	}

	if issimplsym(l.Peek()) {
		var ret []rune
		for issimplsym(l.Peek()) {
			ret = append(ret, l.Next())
		}
		str := string(ret)

		for i := 0; i < len(keywords_tbl); i++ {
			if keywords_tbl[i].val == str {
				lval.node = smt2node{lno, col, keywords_tbl[i].label, str, str}
				return keywords_tbl[i].label
			}
		}
		org := str
		str = strings.Replace(str, "%", "_P_", -1)
		str = strings.Replace(str, "?", "_q_", -1)
		str = strings.Replace(str, "!", "_e_", -1)
		str = strings.Replace(str, ".", "_d_", -1)
		str = strings.Replace(str, "$", "_D_", -1)
		str = strings.Replace(str, "~", "_t_", -1)
		str = strings.Replace(str, "&", "_s_", -1)
		str = strings.Replace(str, "^", "_h_", -1)
		str = strings.Replace(str, "@", "_a_", -1)
		lval.node = smt2node{lno, col, symbol, str, org}
		return symbol
	}

	if c == ':' {
		var ret []rune
		ret = append(ret, l.Next())
		for issimplsym(l.Peek()) {
			ret = append(ret, l.Next())
		}
		str := string(ret)
		if str == ":status" {
			return kw_status
		}
		lval.node = smt2node{lno, col, keyword, str, str}
		return keyword
	}
	if c == '|' || c == '"' {
		var ret []rune
		corg := l.Peek()
		ret = append(ret, l.Next())
		for l.Peek() != corg {
			ret = append(ret, l.Next())
		}
		l.Next()

		str := string(ret)
		if c == '"' {
			lval.node = smt2node{lno, col, string_, str, str}
			return string_
		}

		org := str

		// white space まで許可されるので適当に名前をつけるしかないだろう.
		if s, ok := symbol_map[org]; ok {
			str = s
		} else {
			symbol_cnt += 1
			str = "___BAR_" + strconv.Itoa(symbol_cnt) + "__"
			// str += strings.TrimFunc(org, func(c rune) bool {
			// 	return !isletter(c)
			// })
			symbol_map[org] = str
		}
		lval.node = smt2node{lno, col, symbol, str, org}

		return symbol
	}

	return int(c)
}

func (l *synLex) Error(s string) {
	pos := l.Pos()
	if l.err == nil {
		l.err = errors.New(fmt.Sprintf("%s:%s:Error:%s \n", pos.String(), string(l.Peek()), s))
	}
}

func parse(str string) (*QeStack, []Comment, error) {
	l := new(synLex)
	l.Init(strings.NewReader(str))
	stack = new(QeStack)
	assert_stk = make([]int, 1)
	decfun_cnt = 0
	symbol_cnt = 0
	letmap.reset()
	symbol_map = make(map[string]string)
	yyParse(l)
	return stack, l.comment, l.err
}

func trace(s string) {
	//	fmt.Printf(s + "\n")
}

func print_ints(a []int, label string) {
	fmt.Printf("%s: ", label)
	for i := 0; i < len(a); i++ {
		fmt.Printf("%d ", a[i])
	}
	fmt.Printf("\n")
}

func update_assert_stk(b bool) {
	// b=true なら, assert 追加
	// b=false なら, declare_fun 追加

	var sgn int
	if b {
		sgn = 1
	} else {
		sgn = -1
	}
	if sgn*assert_stk[len(assert_stk)-1] < 0 {
		assert_stk = append(assert_stk, 0)
	}
	assert_stk[len(assert_stk)-1] += sgn
}


func declare_sym(sym smt2node) {
	if assert_stk[len(assert_stk)-1] > 0 {
		assert_stk = append(assert_stk, -1)
	} else {
		if assert_stk[len(assert_stk)-1] < 0 {
			fmt.Printf("pop dec %v\n", assert_stk)
			stack.Pop()
		}
		assert_stk[len(assert_stk)-1] -= 1
	}
	stack.Push(NewQeNodeStr(sym.str, sym.lno))
	stack.Push(NewQeNodeList(-assert_stk[len(assert_stk)-1], sym.lno))
}

