//line synparse.y:3
package qeconv

import __yyfmt__ "fmt"

//line synparse.y:4
import (
	"fmt"
	"text/scanner"
)

var stack *Stack

type Node struct {
	cmd      int
	val      int
	str      string
	rev      bool
	priority int
	lineno   int
}

//line synparse.y:26
type yySymType struct {
	yys  int
	node Node
	num  int
}

const NAME = 57346
const NUMBER = 57347
const F_TRUE = 57348
const F_FALSE = 57349
const ALL = 57350
const EX = 57351
const AND = 57352
const OR = 57353
const NOT = 57354
const PLUS = 57355
const MINUS = 57356
const COMMA = 57357
const MULT = 57358
const DIV = 57359
const POW = 57360
const EOL = 57361
const LB = 57362
const RB = 57363
const LP = 57364
const RP = 57365
const LC = 57366
const RC = 57367
const INDEXED = 57368
const LIST = 57369
const IMPL = 57370
const REPL = 57371
const EQUIV = 57372
const COMMENT = 57373
const LTOP = 57374
const GTOP = 57375
const LEOP = 57376
const GEOP = 57377
const NEOP = 57378
const EQOP = 57379
const UNARYMINUS = 57380
const UNARYPLUS = 57381

var yyToknames = []string{
	"NAME",
	"NUMBER",
	"F_TRUE",
	"F_FALSE",
	"ALL",
	"EX",
	"AND",
	"OR",
	"NOT",
	"PLUS",
	"MINUS",
	"COMMA",
	"MULT",
	"DIV",
	"POW",
	"EOL",
	"LB",
	"RB",
	"LP",
	"RP",
	"LC",
	"RC",
	"INDEXED",
	"LIST",
	"IMPL",
	"REPL",
	"EQUIV",
	"COMMENT",
	"LTOP",
	"GTOP",
	"LEOP",
	"GEOP",
	"NEOP",
	"EQOP",
	"UNARYMINUS",
	"UNARYPLUS",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line synparse.y:144

/*  start  of  programs  */

type SynLex struct {
	scanner.Scanner
	s       string
	comment []Node
}

type SynLex1 struct {
	val      string
	label    int
	v        int
	argn     int // 引数の数
	priority int // 要素に () が必要かを判定するためのフラグ
}

var sones = []SynLex1{
	{"+", PLUS, '+', 2, 4},
	{"-", MINUS, '-', 2, 4},
	{"*", MULT, '*', 2, 3},
	{"/", DIV, '/', 2, 3},
	{"^", POW, '^', 2, 1},
	{"[", LB, '[', 0, 0},
	{"]", RB, ']', 0, 0},
	{"{", LC, '{', 0, 0},
	{"}", RC, '{', 0, 0},
	{"(", LP, '(', 0, 0},
	{")", RP, ')', 0, 0},
	{",", COMMA, ',', 0, 0},
	{":", EOL, ':', 0, 0},
	{"=", EQOP, '=', 0, 0},
}

var sfuns = []SynLex1{
	{"And", AND, 0, 0, 1},
	{"Or", OR, 0, 0, 2},
	{"Impl", IMPL, 0, 2, 0},
	{"Repl", REPL, 0, 2, 0},
	{"Equiv", EQUIV, 0, 2, 0},
	{"Not", NOT, 0, 1, 0},
	{"All", ALL, 0, 2, 0},
	{"Ex", EX, 0, 2, 0},
	{"true", F_TRUE, 0, 0, 0},
	{"false", F_FALSE, 0, 0, 0},
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

func (l *SynLex) Lex(lval *yySymType) int {

	// skip space
	for {
		for isspace(l.Peek()) {
			l.Next()
		}
		lno := l.Pos().Line
		if l.Peek() != '#' {
			break
		}
		l.Next()
		str := ""
		for l.Peek() != '\n' {
			str += string(l.Next())
		}
		if str != "" {
			fmt.Printf("comment %d: %s\n", lno, str)
			l.comment = append(l.comment, Node{str: str, lineno: lno})
		}
	}

	lno := l.Pos().Line
	c := int(l.Peek())
	for i := 0; i < len(sones); i++ {
		if sones[i].v == c {
			l.Next()
			lval.node = Node{cmd: sones[i].label, val: sones[i].argn, str: sones[i].val, priority: sones[i].priority, lineno: lno}
			return sones[i].label
		}
	}
	if c == '>' {
		l.Next()
		if l.Peek() == '=' {
			l.Next()
			return GEOP
		} else {
			return GTOP
		}
	} else if c == '<' {
		l.Next()
		if l.Peek() == '=' {
			l.Next()
			return LEOP
		} else if l.Peek() == '>' {
			l.Next()
			return NEOP
		} else {
			return LTOP
		}
	}

	if isdigit(l.Peek()) {
		var ret []rune
		for isdigit(l.Peek()) {
			ret = append(ret, l.Next())
		}
		lval.node = Node{cmd: NUMBER, val: 0, str: string(ret), lineno: lno}
		return NUMBER
	}

	if isalnum(l.Peek()) || c == '_' {
		var ret []rune
		for isdigit(l.Peek()) || isletter(l.Peek()) {
			ret = append(ret, l.Next())
		}
		lval.node = Node{cmd: NAME, val: 0, str: string(ret), lineno: lno}
		for i := 0; i < len(sfuns); i++ {
			if lval.node.str == sfuns[i].val {
				lval.node = Node{cmd: sfuns[i].label, val: sfuns[i].argn, str: sfuns[i].val, priority: sfuns[i].priority, lineno: lno}
				return sfuns[i].label
			}
		}

		return NAME
	}

	return int(c)
}

func (l *SynLex) Error(s string) {
	fmt.Printf("syntax error: %s\n", s)
	panic(s)
}

func parse(l *SynLex) *Stack {
	stack = new(Stack)
	yyParse(l)
	return stack
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

const yyNprod = 42
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 179

var yyAct = []int{

	37, 2, 19, 17, 81, 94, 35, 105, 104, 29,
	103, 66, 102, 33, 100, 95, 34, 44, 45, 86,
	46, 47, 48, 63, 50, 52, 53, 64, 54, 54,
	32, 60, 61, 62, 58, 59, 38, 39, 40, 41,
	43, 42, 67, 68, 69, 70, 71, 72, 73, 74,
	75, 76, 66, 57, 31, 79, 22, 18, 82, 82,
	85, 83, 22, 94, 66, 21, 20, 90, 30, 93,
	65, 91, 44, 45, 51, 46, 47, 48, 55, 28,
	49, 92, 56, 27, 26, 96, 25, 24, 97, 98,
	99, 38, 39, 40, 41, 43, 42, 101, 22, 18,
	15, 16, 4, 5, 6, 7, 8, 21, 20, 44,
	45, 23, 46, 47, 48, 36, 12, 48, 89, 64,
	88, 87, 9, 10, 11, 22, 18, 15, 16, 4,
	5, 6, 7, 8, 21, 20, 46, 47, 48, 84,
	80, 14, 78, 12, 77, 22, 13, 3, 1, 9,
	10, 11, 22, 18, 15, 16, 4, 5, 6, 7,
	8, 21, 20, 44, 45, 0, 46, 47, 48, 0,
	12, 0, 0, 0, 0, 0, 9, 10, 11,
}
var yyPact = []int{

	121, -1000, 92, 68, 64, 62, 61, 57, 148, 46,
	32, 8, 148, -1000, 94, -1000, -1000, 59, -1000, 60,
	52, 52, -1000, -1000, -1000, 58, 58, 148, 148, -1000,
	148, 148, 148, 0, 4, 49, -1000, -1000, 52, 52,
	52, 52, 52, 52, 52, 52, 52, 52, 139, 137,
	99, 52, 99, 125, 60, 141, 141, 124, 37, -4,
	106, 105, 103, -1000, -1000, -1000, 148, 150, 150, 150,
	150, 150, 150, 120, 120, 99, 99, -1000, 50, 96,
	148, 48, 60, -10, 148, -1000, -1000, 148, 148, 148,
	-1000, -1000, -9, -1000, 141, -1000, -11, -13, -15, -16,
	-1000, 60, -1000, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 4, 6, 2, 148, 0, 147, 26, 146, 3,
}
var yyR1 = []int{

	0, 4, 4, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 6, 6, 2, 2, 7, 7, 7,
	1, 1, 3, 3, 8, 8, 8, 8, 8, 8,
	8, 8, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9,
}
var yyR2 = []int{

	0, 2, 2, 6, 6, 4, 4, 2, 6, 6,
	6, 3, 1, 3, 2, 1, 3, 1, 3, 3,
	1, 3, 1, 4, 1, 1, 3, 3, 3, 3,
	3, 3, 3, 1, 1, 3, 3, 3, 3, 3,
	2, 2,
}
var yyChk = []int{

	-1000, -4, -5, -6, 8, 9, 10, 11, 12, 28,
	29, 30, 22, -8, 20, 6, 7, -9, 5, -3,
	14, 13, 4, 19, 19, 22, 22, 22, 22, -5,
	22, 22, 22, -5, -9, -2, 21, -5, 32, 33,
	34, 35, 37, 36, 13, 14, 16, 17, 18, 20,
	-9, 22, -9, -7, -3, 20, 24, -7, -2, -2,
	-5, -5, -5, 23, 23, 21, 15, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, 5, 5, -9,
	15, -1, -3, -1, 15, 23, 23, 15, 15, 15,
	-5, 21, -5, 21, 15, 25, -5, -5, -5, -5,
	23, -3, 23, 23, 23, 23,
}
var yyDef = []int{

	0, -2, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 12, 0, 24, 25, 0, 33, 34,
	0, 0, 22, 1, 2, 0, 0, 0, 0, 7,
	0, 0, 0, 0, 0, 0, 14, 15, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	40, 0, 41, 0, 17, 0, 0, 0, 0, 0,
	0, 0, 0, 11, 32, 13, 0, 26, 27, 28,
	29, 30, 31, 35, 36, 37, 38, 39, 0, 0,
	0, 0, 20, 0, 0, 5, 6, 0, 0, 0,
	16, 23, 0, 18, 0, 19, 0, 0, 0, 0,
	3, 21, 4, 8, 9, 10,
}
var yyTok1 = []int{

	1,
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

	case 3:
		//line synparse.y:64
		{
			trace("ALL")
			stack.push(yyS[yypt-5].node)
		}
	case 4:
		//line synparse.y:65
		{
			trace("EX")
			stack.push(yyS[yypt-5].node)
		}
	case 5:
		//line synparse.y:66
		{
			trace("and")
			yyS[yypt-3].node.val = yyS[yypt-1].num
			stack.push(yyS[yypt-3].node)
		}
	case 6:
		//line synparse.y:67
		{
			trace("or")
			yyS[yypt-3].node.val = yyS[yypt-1].num
			stack.push(yyS[yypt-3].node)
		}
	case 7:
		//line synparse.y:68
		{
			trace("not")
			stack.push(yyS[yypt-1].node)
		}
	case 8:
		//line synparse.y:69
		{
			trace("IMPL")
			stack.push(yyS[yypt-5].node)
		}
	case 9:
		//line synparse.y:70
		{
			trace("REPL")
			stack.push(yyS[yypt-5].node)
		}
	case 10:
		//line synparse.y:71
		{
			trace("EQUIV")
			stack.push(yyS[yypt-5].node)
		}
	case 13:
		//line synparse.y:77
		{
			trace("list")
			stack.push(Node{cmd: LIST, val: yyS[yypt-1].num})
		}
	case 14:
		//line synparse.y:81
		{
			trace("empty-list")
			stack.push(Node{cmd: LIST, val: 0})
		}
	case 15:
		//line synparse.y:88
		{
			yyVAL.num = 1
		}
	case 16:
		//line synparse.y:89
		{
			yyVAL.num = yyS[yypt-2].num + 1
		}
	case 17:
		//line synparse.y:93
		{
			trace("list")
			stack.push(Node{cmd: LIST, val: 1})
		}
	case 18:
		//line synparse.y:97
		{
			trace("list")
			stack.push(Node{cmd: LIST, val: yyS[yypt-1].num})
		}
	case 19:
		//line synparse.y:101
		{
			trace("set")
			stack.push(Node{cmd: LIST, val: yyS[yypt-1].num})
		}
	case 20:
		//line synparse.y:108
		{
			yyVAL.num = 1
		}
	case 21:
		//line synparse.y:109
		{
			yyVAL.num = yyS[yypt-2].num + 1
		}
	case 22:
		//line synparse.y:113
		{
			trace("name")
			stack.push(yyS[yypt-0].node)
		}
	case 23:
		//line synparse.y:114
		{
			trace("index")
			stack.push(Node{cmd: INDEXED, val: 2})
		}
	case 24:
		//line synparse.y:118
		{
			trace("true")
			stack.push(Node{cmd: F_TRUE, val: 0})
		}
	case 25:
		//line synparse.y:119
		{
			trace("false")
			stack.push(Node{cmd: F_FALSE, val: 0})
		}
	case 26:
		//line synparse.y:120
		{
			trace("<")
			stack.push(Node{cmd: LTOP, str: "<", val: 2})
		}
	case 27:
		//line synparse.y:121
		{
			trace(">")
			stack.push(Node{cmd: LTOP, str: ">", val: 2, rev: true})
		}
	case 28:
		//line synparse.y:122
		{
			trace("<=")
			stack.push(Node{cmd: LEOP, str: "<=", val: 2})
		}
	case 29:
		//line synparse.y:123
		{
			trace(">=")
			stack.push(Node{cmd: LEOP, str: ">=", val: 2, rev: true})
		}
	case 30:
		//line synparse.y:124
		{
			trace("=")
			stack.push(Node{cmd: EQOP, str: "=", val: 2})
		}
	case 31:
		//line synparse.y:125
		{
			trace("<>")
			stack.push(Node{cmd: NEOP, str: "<>", val: 2})
		}
	case 33:
		//line synparse.y:131
		{
			trace("num")
			stack.push(yyS[yypt-0].node)
		}
	case 35:
		//line synparse.y:133
		{
			trace("+")
			stack.push(yyS[yypt-1].node)
		}
	case 36:
		//line synparse.y:134
		{
			trace("-")
			stack.push(yyS[yypt-1].node)
		}
	case 37:
		//line synparse.y:135
		{
			trace("*")
			stack.push(yyS[yypt-1].node)
		}
	case 38:
		//line synparse.y:136
		{
			trace("/")
			stack.push(yyS[yypt-1].node)
		}
	case 39:
		//line synparse.y:137
		{
			trace("^")
			stack.push(yyS[yypt-0].node)
			stack.push(yyS[yypt-1].node)
		}
	case 40:
		//line synparse.y:138
		{
			trace("-")
			yyS[yypt-1].node.cmd = UNARYMINUS
			yyS[yypt-1].node.val = 1
			yyS[yypt-1].node.priority = 2
			stack.push(yyS[yypt-1].node)
		}
	case 41:
		//line synparse.y:140
		{
			trace("+")
			yyS[yypt-1].node.cmd = UNARYPLUS
			yyS[yypt-1].node.val = 1
			yyS[yypt-1].node.priority = 2
			stack.push(yyS[yypt-1].node)
		}
	}
	goto yystack /* stack new state and value */
}
