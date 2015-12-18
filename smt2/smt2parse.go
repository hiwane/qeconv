//line smt2parse.y:3
package qeconv

import __yyfmt__ "fmt"

//line smt2parse.y:4
import (
	"errors"
	"fmt"
	. "github.com/hiwane/qeconv/def"
	"strings"
	"text/scanner"
)

var stack *QeStack
var assert_cnt int
var decfun_cnt int

type smt2node struct {
	lno, col int
	typ      int
	str      string
}

//line smt2parse.y:27
type yySymType struct {
	yys  int
	node smt2node
	num  int
}

const number = 57346
const symbol = 57347
const keyword = 57348
const string_ = 57349
const forall = 57350
const exists = 57351
const let = 57352
const as = 57353
const theory = 57354
const par = 57355
const assert = 57356
const check_sat = 57357
const declare_fun = 57358
const set_info = 57359
const set_logic = 57360
const exit = 57361
const set_option = 57362
const ltop = 57363
const gtop = 57364
const leop = 57365
const geop = 57366
const eqop = 57367
const plus = 57368
const minus = 57369
const mult = 57370
const div = 57371
const not = 57372
const and = 57373
const or = 57374
const implies = 57375
const impl = 57376
const repl = 57377
const equiv = 57378
const unaryminus = 57379
const unaryplus = 57380
const pow = 57381

var yyToknames = []string{
	"number",
	"symbol",
	"keyword",
	"string_",
	"forall",
	"exists",
	"let",
	"as",
	"theory",
	"par",
	"assert",
	"check_sat",
	"declare_fun",
	"set_info",
	"set_logic",
	"exit",
	"set_option",
	"ltop",
	"gtop",
	"leop",
	"geop",
	"eqop",
	"plus",
	"minus",
	"mult",
	"div",
	"not",
	"and",
	"or",
	"implies",
	"impl",
	"repl",
	"equiv",
	"unaryminus",
	"unaryplus",
	"pow",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line smt2parse.y:247

/*  start  of  programs  */

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

func (l *synLex) Lex(lval *yySymType) int {

	// skip space
	for isspace(l.Peek()) {
		l.Next()
	}
	if scanner.EOF == l.Peek() {
		trace("Lex! eof " + l.Pos().String())
		l.Next()
		return 0
	}
	trace("Lex! " + string(l.Peek()) + l.Pos().String())

	c := int(l.Peek())
	lno := l.Pos().Line
	col := l.Pos().Column
	if c == '(' || c == ')' {
		l.Next()
		return c
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

		lval.node = smt2node{lno, col, number, string(ret)}
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
				lval.node = smt2node{lno, col, keywords_tbl[i].label, str}
				return keywords_tbl[i].label
			}
		}
		str = strings.Replace(str, "?", "_q_", -1)
		str = strings.Replace(str, "!", "_e_", -1)
		lval.node = smt2node{lno, col, symbol, str}
		return symbol
	}

	if c == ':' {
		var ret []rune
		ret = append(ret, l.Next())
		for issimplsym(l.Peek()) {
			ret = append(ret, l.Next())
		}
		lval.node = smt2node{lno, col, keyword, string(ret)}
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
		if c == '|' {
			lval.node = smt2node{lno, col, symbol, string(ret)}
		} else {
			lval.node = smt2node{lno, col, string_, string(ret)}
		}
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
	assert_cnt = 0
	decfun_cnt = 0
	yyParse(l)
	return stack, l.comment, l.err
}

func trace(s string) {
	//	fmt.Printf(s + "\n")
}

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 58
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 206

var yyAct = []int{

	52, 13, 18, 80, 77, 73, 71, 12, 81, 111,
	79, 132, 16, 102, 103, 17, 131, 16, 102, 103,
	17, 81, 108, 78, 105, 130, 47, 129, 128, 126,
	118, 117, 116, 115, 19, 19, 60, 61, 62, 63,
	64, 65, 16, 19, 68, 17, 114, 72, 104, 127,
	113, 112, 75, 104, 100, 51, 92, 75, 75, 75,
	75, 87, 88, 89, 90, 91, 82, 75, 75, 95,
	98, 69, 72, 72, 50, 101, 96, 97, 15, 94,
	45, 106, 26, 109, 21, 20, 109, 56, 57, 58,
	59, 81, 78, 55, 54, 53, 44, 66, 67, 4,
	24, 72, 19, 110, 14, 119, 121, 107, 122, 123,
	120, 19, 125, 72, 29, 30, 28, 124, 16, 19,
	27, 17, 101, 16, 19, 25, 17, 38, 36, 37,
	35, 39, 31, 32, 33, 34, 40, 41, 42, 43,
	16, 19, 22, 17, 16, 19, 23, 17, 16, 19,
	3, 17, 46, 5, 15, 93, 70, 99, 2, 15,
	86, 16, 19, 1, 17, 16, 19, 76, 17, 16,
	48, 0, 17, 0, 0, 0, 15, 85, 0, 0,
	15, 84, 0, 0, 15, 83, 6, 7, 9, 10,
	11, 8, 0, 0, 0, 0, 0, 15, 74, 0,
	0, 15, 0, 0, 0, 49,
}
var yyPact = []int{

	59, -1000, 59, -1000, 172, -1000, 161, 44, 43, 137,
	94, 120, 41, -1000, -1000, 106, -1000, -1000, -1000, -1000,
	-1000, -1000, 56, 39, 165, 33, -1000, 161, 55, 54,
	53, 161, 161, 161, 161, 161, 161, 161, 161, 161,
	161, 161, 161, 161, 30, -1000, -1000, -1000, -1000, -1000,
	-1000, 157, -1000, 52, 51, 51, 144, 140, 136, 119,
	161, 161, 161, 161, 161, 15, 114, 38, 161, 97,
	29, -1000, -1000, 13, -1000, -1000, -17, -1000, 102, -19,
	-1000, 98, -32, -1000, -1000, -1000, -1000, 10, 9, 5,
	-8, -9, -1000, -1000, -1000, -10, -11, -1000, 97, -1000,
	-1000, -1000, -1000, -1000, -1000, 161, -1000, 161, 161, -1000,
	97, 161, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -12,
	8, -13, -14, -16, -25, -30, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000,
}
var yyPgo = []int{

	0, 2, 1, 5, 55, 167, 10, 163, 158, 150,
	157, 6, 156, 152, 146, 146, 0, 104, 3, 4,
}
var yyR1 = []int{

	0, 7, 8, 8, 10, 10, 10, 10, 3, 3,
	2, 2, 1, 11, 12, 12, 13, 13, 13, 14,
	14, 15, 15, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 4, 4, 6, 6, 5, 5, 18, 19,
	17, 9, 9, 9, 9, 9, 9, 9,
}
var yyR2 = []int{

	0, 1, 1, 2, 1, 1, 1, 3, 0, 2,
	1, 1, 1, 1, 1, 2, 1, 1, 3, 1,
	2, 1, 2, 1, 1, 4, 7, 7, 7, 4,
	4, 4, 4, 5, 5, 5, 5, 5, 4, 4,
	4, 5, 1, 2, 1, 2, 1, 2, 4, 4,
	1, 4, 3, 3, 7, 8, 4, 4,
}
var yyChk = []int{

	-1000, -7, -8, -9, 40, -9, 14, 15, 19, 16,
	17, 18, -16, -2, -17, 40, 4, 7, -1, 5,
	41, 41, 5, -14, 6, 5, 41, -17, 10, 8,
	9, 26, 27, 28, 29, 24, 22, 23, 21, 25,
	30, 31, 32, 33, 40, 41, -13, -2, 5, 40,
	41, -4, -16, 40, 40, 40, -4, -4, -4, -4,
	-16, -16, -16, -16, -16, -16, -4, -4, -16, 41,
	-12, -11, -1, -3, 41, -16, -5, -19, 40, -6,
	-18, 40, -6, 41, 41, 41, 41, -16, -16, -16,
	-16, -16, 41, 41, 41, -16, -11, -11, 41, -10,
	41, -2, 5, 6, 40, 41, -19, 5, 41, -18,
	5, 41, 41, 41, 41, 41, 41, 41, 41, -11,
	-3, -16, -16, -16, -11, -16, 41, 41, 41, 41,
	41, 41, 41,
}
var yyDef = []int{

	0, -2, 1, 2, 0, 3, 0, 0, 0, 0,
	0, 0, 0, 23, 24, 0, 10, 11, 50, 12,
	52, 53, 0, 0, 19, 0, 51, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 56, 20, 16, 17, 8,
	57, 0, 42, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 14, 13, 0, 25, 43, 0, 46, 0, 0,
	44, 0, 0, 29, 30, 31, 32, 0, 0, 0,
	0, 0, 38, 39, 40, 0, 0, 15, 0, 9,
	18, 4, 5, 6, 8, 0, 47, 0, 0, 45,
	0, 0, 33, 34, 35, 36, 37, 41, 54, 0,
	0, 0, 0, 0, 0, 0, 55, 7, 26, 49,
	27, 48, 28,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	40, 41,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(c), uint(char))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line smt2parse.y:60
		{
			trace("eof")
		}
	case 2:
		//line smt2parse.y:64
		{
			trace("command")
		}
	case 3:
		//line smt2parse.y:65
		{
			trace("commands")
		}
	case 8:
		//line smt2parse.y:76
		{
			yyVAL.num = 0
		}
	case 9:
		//line smt2parse.y:77
		{
			yyVAL.num = yyS[yypt-1].num + 1
		}
	case 10:
		//line smt2parse.y:81
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 11:
		//line smt2parse.y:82
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 12:
		//line smt2parse.y:90
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 13:
		//line smt2parse.y:92
		{
			if yyS[yypt-0].node.str != "Real" {
				yylex.Error("unknown sort")
			}
		}
	case 23:
		//line smt2parse.y:115
		{
			stack.Push(NewQeNodeNum(yyS[yypt-0].node.str, yyS[yypt-0].node.lno))
		}
	case 29:
		//line smt2parse.y:122
		{
			if yyS[yypt-1].num > 1 {
				stack.Push(NewQeNodeStrVal(yyS[yypt-2].node.str, yyS[yypt-1].num, yyS[yypt-2].node.lno))
			}
		}
	case 30:
		//line smt2parse.y:126
		{
			if yyS[yypt-1].num == 1 {
				stack.Push(NewQeNodeStr("-.", yyS[yypt-2].node.lno))
			} else {
				stack.Push(NewQeNodeStrVal(yyS[yypt-2].node.str, yyS[yypt-1].num, yyS[yypt-2].node.lno))
			}
		}
	case 31:
		//line smt2parse.y:132
		{
			stack.Push(NewQeNodeStrVal(yyS[yypt-2].node.str, yyS[yypt-1].num, yyS[yypt-2].node.lno))
		}
	case 32:
		//line smt2parse.y:133
		{
			stack.Push(NewQeNodeStrVal(yyS[yypt-2].node.str, yyS[yypt-1].num, yyS[yypt-2].node.lno))
		}
	case 33:
		//line smt2parse.y:134
		{
			stack.Push(NewQeNodeStr(yyS[yypt-3].node.str, yyS[yypt-3].node.lno))
		}
	case 34:
		//line smt2parse.y:135
		{
			stack.Push(NewQeNodeStr(yyS[yypt-3].node.str, yyS[yypt-3].node.lno))
		}
	case 35:
		//line smt2parse.y:136
		{
			stack.Push(NewQeNodeStr(yyS[yypt-3].node.str, yyS[yypt-3].node.lno))
		}
	case 36:
		//line smt2parse.y:137
		{
			stack.Push(NewQeNodeStr(yyS[yypt-3].node.str, yyS[yypt-3].node.lno))
		}
	case 37:
		//line smt2parse.y:138
		{
			stack.Push(NewQeNodeStr(yyS[yypt-3].node.str, yyS[yypt-3].node.lno))
		}
	case 38:
		//line smt2parse.y:139
		{
			stack.Push(NewQeNodeStr("Not", yyS[yypt-2].node.lno))
		}
	case 39:
		//line smt2parse.y:140
		{
			stack.Push(NewQeNodeStrVal("And", yyS[yypt-1].num, yyS[yypt-2].node.lno))
		}
	case 40:
		//line smt2parse.y:141
		{
			stack.Push(NewQeNodeStrVal("Or", yyS[yypt-1].num, yyS[yypt-2].node.lno))
		}
	case 41:
		//line smt2parse.y:142
		{
			stack.Push(NewQeNodeStr("Impl", yyS[yypt-3].node.lno))
		}
	case 42:
		//line smt2parse.y:145
		{
			yyVAL.num = 1
		}
	case 43:
		//line smt2parse.y:146
		{
			yyVAL.num = yyS[yypt-1].num + 1
		}
	case 44:
		//line smt2parse.y:150
		{
			yyVAL.num = 1
		}
	case 45:
		//line smt2parse.y:151
		{
			yyVAL.num = yyS[yypt-1].num + 1
		}
	case 46:
		//line smt2parse.y:155
		{
			yyVAL.num = 1
		}
	case 47:
		//line smt2parse.y:156
		{
			yyVAL.num = yyS[yypt-1].num + 1
		}
	case 50:
		//line smt2parse.y:164
		{
			stack.Push(NewQeNodeStr(yyS[yypt-0].node.str, yyS[yypt-0].node.lno))
		}
	case 51:
		//line smt2parse.y:225
		{
			assert_cnt += 1
		}
	case 52:
		//line smt2parse.y:226
		{
			trace("go check-sat")
			stack.Push(NewQeNodeStrVal("And", assert_cnt, 0))
			for i := 0; i < decfun_cnt; i++ {
				stack.Push(NewQeNodeStr("Ex", 0))
			}

		}
	case 54:
		//line smt2parse.y:234
		{
			stack.Push(NewQeNodeStr(yyS[yypt-4].node.str, yyS[yypt-4].node.lno))
			stack.Push(NewQeNodeList(1, yyS[yypt-4].node.lno))
			decfun_cnt += 1
		}
	case 55:
		//line smt2parse.y:239
		{
			yylex.Error("unknown declare")
		}
	case 57:
		//line smt2parse.y:242
		{
			if yyS[yypt-1].node.str != "QF_NRA" && yyS[yypt-1].node.str != "NRA" {
				yylex.Error("unknown logic")
			}
		}
	}
	goto yystack /* stack new state and value */
}
