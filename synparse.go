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
const ABS = 57355
const PLUS = 57356
const MINUS = 57357
const COMMA = 57358
const MULT = 57359
const DIV = 57360
const POW = 57361
const EOL = 57362
const LB = 57363
const RB = 57364
const LP = 57365
const RP = 57366
const LC = 57367
const RC = 57368
const INDEXED = 57369
const LIST = 57370
const IMPL = 57371
const REPL = 57372
const EQUIV = 57373
const COMMENT = 57374
const LTOP = 57375
const GTOP = 57376
const LEOP = 57377
const GEOP = 57378
const NEOP = 57379
const EQOP = 57380
const UNARYMINUS = 57381
const UNARYPLUS = 57382

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
	"ABS",
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

//line synparse.y:169

/*  start  of  programs  */

type Comment struct {
	lineno int
	str    string
}

type SynLex struct {
	scanner.Scanner
	s       string
	comment []Comment
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
	{"}", RC, '}', 0, 0},
	{"(", LP, '(', 0, 0},
	{")", RP, ')', 0, 0},
	{",", COMMA, ',', 0, 0},
	{":", EOL, ':', 0, 0},
	{"=", EQOP, '=', 0, 0},
}

var sfuns = []SynLex1{
	{"And", AND, 0, 0, 1},
	{"Or", OR, 0, 0, 2},
	{"Impl", IMPL, 0, 2, 3},
	{"Repl", REPL, 0, 2, 3},
	{"Equiv", EQUIV, 0, 2, 3},
	{"Not", NOT, 0, 1, 0},
	{"All", ALL, 0, 2, 0},
	{"Ex", EX, 0, 2, 0},
	{"true", F_TRUE, 0, 0, 0},
	{"false", F_FALSE, 0, 0, 0},
	{"abs", ABS, 0, 0, 0},
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
			l.comment = append(l.comment, Comment{str: str, lineno: lno})
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

				// Repl は Impl に変換する.
				if lval.node.str == "Repl" {
					lval.node.rev = true
					lval.node.cmd = IMPL
				}
				return sfuns[i].label
			}
		}

		return NAME
	}

	return int(c)
}

func (l *SynLex) Error(s string) {
	pos := l.Pos()
	fmt.Printf("%s:Error:%s \n", pos.String(), s)
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

const yyNprod = 56
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 281

var yyAct = []int{

	20, 3, 82, 2, 55, 130, 113, 60, 83, 129,
	128, 30, 127, 91, 51, 34, 114, 86, 85, 38,
	45, 46, 18, 47, 48, 49, 84, 56, 56, 125,
	62, 62, 59, 65, 66, 67, 35, 63, 68, 39,
	40, 41, 42, 44, 43, 52, 54, 45, 46, 33,
	47, 48, 49, 96, 113, 32, 24, 69, 92, 92,
	112, 97, 72, 73, 74, 75, 76, 77, 78, 79,
	80, 81, 93, 57, 88, 101, 89, 58, 96, 31,
	45, 46, 29, 47, 48, 49, 95, 106, 107, 108,
	110, 71, 111, 28, 27, 26, 115, 70, 116, 109,
	117, 118, 119, 50, 25, 120, 121, 122, 123, 47,
	48, 49, 45, 46, 126, 47, 48, 49, 104, 105,
	102, 103, 69, 104, 105, 100, 49, 99, 98, 94,
	124, 39, 40, 41, 42, 44, 43, 24, 19, 16,
	17, 5, 6, 7, 8, 9, 21, 23, 22, 87,
	90, 24, 14, 4, 15, 37, 13, 1, 36, 0,
	0, 0, 10, 11, 12, 24, 19, 16, 17, 5,
	6, 7, 8, 9, 21, 23, 22, 0, 0, 0,
	0, 0, 15, 0, 13, 0, 0, 0, 0, 0,
	10, 11, 12, 24, 19, 16, 17, 5, 6, 7,
	8, 9, 21, 23, 22, 45, 46, 0, 47, 48,
	49, 0, 13, 64, 0, 0, 0, 0, 10, 11,
	12, 24, 19, 16, 17, 5, 6, 7, 8, 9,
	21, 23, 22, 0, 0, 0, 0, 0, 0, 0,
	13, 61, 0, 0, 0, 0, 10, 11, 12, 24,
	19, 16, 17, 5, 6, 7, 8, 9, 21, 23,
	22, 24, 19, 0, 0, 0, 0, 0, 13, 0,
	21, 23, 22, 0, 10, 11, 12, 0, 0, 0,
	53,
}
var yyPact = []int{

	161, -1000, 84, -1000, -1000, 72, 71, 70, 59, 245,
	56, 32, 26, 245, -1000, 133, -1000, -1000, 6, -1000,
	82, -9, 257, 257, -1000, -1000, 52, 52, 217, 189,
	-1000, 245, 245, 245, 14, 98, 75, -1000, -1000, 257,
	257, 257, 257, 257, 257, 257, 257, 257, 257, 3,
	144, 257, 107, 257, 107, 134, 82, 147, 147, 113,
	62, -1000, -1000, 37, -1000, 112, 111, 109, -1000, -1000,
	-1000, 161, 191, 191, 191, 191, 191, 191, 92, 92,
	107, 107, -1000, -1000, 3, 3, 3, 77, 66, 33,
	245, 38, 82, -10, 245, -1000, 245, -1000, 245, 245,
	245, -1000, 3, 3, 3, 3, 106, -1000, -1000, -1000,
	-1000, 5, -1000, 147, -1000, -12, -1000, -14, -15, -19,
	101, 101, -1000, -1000, -1000, -1000, 82, -1000, -1000, -1000,
	-1000,
}
var yyPgo = []int{

	0, 13, 7, 158, 2, 0, 157, 3, 1, 153,
	4, 152, 22,
}
var yyR1 = []int{

	0, 6, 7, 7, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 9, 9, 3, 3,
	2, 2, 10, 10, 10, 1, 1, 5, 5, 11,
	11, 11, 11, 11, 11, 11, 11, 4, 4, 4,
	4, 4, 4, 4, 4, 12, 12, 12, 12, 12,
	12, 12, 12, 12, 12, 12,
}
var yyR2 = []int{

	0, 2, 1, 1, 6, 6, 4, 3, 4, 3,
	2, 6, 6, 6, 3, 1, 3, 2, 1, 3,
	1, 3, 1, 3, 3, 1, 3, 1, 4, 1,
	1, 3, 3, 3, 3, 3, 3, 1, 3, 3,
	3, 3, 3, 2, 2, 3, 1, 1, 4, 3,
	3, 3, 3, 3, 2, 2,
}
var yyChk = []int{

	-1000, -6, -7, -8, -9, 8, 9, 10, 11, 12,
	29, 30, 31, 23, -11, 21, 6, 7, -12, 5,
	-5, 13, 15, 14, 4, 20, 23, 23, 23, 23,
	-8, 23, 23, 23, -8, -12, -3, 22, -7, 33,
	34, 35, 36, 38, 37, 14, 15, 17, 18, 19,
	21, 23, -12, 23, -12, -10, -5, 21, 25, -10,
	-2, 24, -8, -2, 24, -8, -8, -8, 24, 24,
	22, 16, -12, -12, -12, -12, -12, -12, -12, -12,
	-12, -12, -4, 5, 23, 15, 14, 5, -12, -12,
	16, -1, -5, -1, 16, 24, 16, 24, 16, 16,
	16, -7, 14, 15, 17, 18, -4, -4, -4, 22,
	24, -8, 22, 16, 26, -8, -8, -8, -8, -8,
	-4, -4, -4, -4, 24, 24, -5, 24, 24, 24,
	24,
}
var yyDef = []int{

	0, -2, 0, 2, 3, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 15, 0, 29, 30, 0, 46,
	47, 0, 0, 0, 27, 1, 0, 0, 0, 0,
	10, 0, 0, 0, 0, 0, 0, 17, 18, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 54, 0, 55, 0, 22, 0, 0, 0,
	0, 7, 20, 0, 9, 0, 0, 0, 14, 45,
	16, 0, 31, 32, 33, 34, 35, 36, 49, 50,
	51, 52, 53, 37, 0, 0, 0, 0, 0, 0,
	0, 0, 25, 0, 0, 6, 0, 8, 0, 0,
	0, 19, 0, 0, 0, 0, 0, 43, 44, 28,
	48, 0, 23, 0, 24, 0, 21, 0, 0, 0,
	39, 40, 41, 42, 38, 4, 26, 5, 11, 12,
	13,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40,
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

	case 4:
		//line synparse.y:67
		{
			trace("ALL")
			stack.push(yyS[yypt-5].node)
		}
	case 5:
		//line synparse.y:68
		{
			trace("EX")
			stack.push(yyS[yypt-5].node)
		}
	case 6:
		//line synparse.y:69
		{
			trace("and")
			yyS[yypt-3].node.val = yyS[yypt-1].num
			stack.push(yyS[yypt-3].node)
		}
	case 7:
		//line synparse.y:70
		{
			trace("and()")
			stack.push(Node{cmd: F_TRUE, val: 0})
		}
	case 8:
		//line synparse.y:71
		{
			trace("or")
			yyS[yypt-3].node.val = yyS[yypt-1].num
			stack.push(yyS[yypt-3].node)
		}
	case 9:
		//line synparse.y:72
		{
			trace("or()")
			stack.push(Node{cmd: F_FALSE, val: 0})
		}
	case 10:
		//line synparse.y:73
		{
			trace("not")
			stack.push(yyS[yypt-1].node)
		}
	case 11:
		//line synparse.y:74
		{
			trace("IMPL")
			stack.push(yyS[yypt-5].node)
		}
	case 12:
		//line synparse.y:75
		{
			trace("REPL")
			stack.push(yyS[yypt-5].node)
		}
	case 13:
		//line synparse.y:76
		{
			trace("EQUIV")
			stack.push(yyS[yypt-5].node)
		}
	case 16:
		//line synparse.y:82
		{
			trace("list")
			stack.push(Node{cmd: LIST, val: yyS[yypt-1].num, lineno: yyS[yypt-2].node.lineno})
		}
	case 17:
		//line synparse.y:86
		{
			trace("empty-list")
			stack.push(Node{cmd: LIST, val: 0, lineno: yyS[yypt-1].node.lineno})
		}
	case 18:
		//line synparse.y:93
		{
			yyVAL.num = 1
		}
	case 19:
		//line synparse.y:94
		{
			yyVAL.num = yyS[yypt-2].num + 1
		}
	case 20:
		//line synparse.y:98
		{
			yyVAL.num = 1
		}
	case 21:
		//line synparse.y:99
		{
			yyVAL.num = yyS[yypt-2].num + 1
		}
	case 22:
		//line synparse.y:103
		{
			trace("var")
			stack.push(Node{cmd: LIST, val: 1})
		}
	case 23:
		//line synparse.y:107
		{
			trace("list")
			stack.push(Node{cmd: LIST, val: yyS[yypt-1].num, lineno: yyS[yypt-2].node.lineno})
		}
	case 24:
		//line synparse.y:111
		{
			trace("set")
			stack.push(Node{cmd: LIST, val: yyS[yypt-1].num, lineno: yyS[yypt-2].node.lineno})
		}
	case 25:
		//line synparse.y:118
		{
			yyVAL.num = 1
		}
	case 26:
		//line synparse.y:119
		{
			yyVAL.num = yyS[yypt-2].num + 1
		}
	case 27:
		//line synparse.y:123
		{
			trace("name")
			stack.push(yyS[yypt-0].node)
		}
	case 28:
		//line synparse.y:124
		{
			trace("index")
			stack.push(Node{cmd: INDEXED, val: 2})
		}
	case 29:
		//line synparse.y:128
		{
			trace("true")
			stack.push(Node{cmd: F_TRUE, val: 0})
		}
	case 30:
		//line synparse.y:129
		{
			trace("false")
			stack.push(Node{cmd: F_FALSE, val: 0})
		}
	case 31:
		//line synparse.y:130
		{
			trace("<")
			stack.push(Node{cmd: LTOP, str: "<", val: 2})
		}
	case 32:
		//line synparse.y:131
		{
			trace(">")
			stack.push(Node{cmd: LTOP, str: ">", val: 2, rev: true})
		}
	case 33:
		//line synparse.y:132
		{
			trace("<=")
			stack.push(Node{cmd: LEOP, str: "<=", val: 2})
		}
	case 34:
		//line synparse.y:133
		{
			trace(">=")
			stack.push(Node{cmd: LEOP, str: ">=", val: 2, rev: true})
		}
	case 35:
		//line synparse.y:134
		{
			trace("=")
			stack.push(Node{cmd: EQOP, str: "=", val: 2})
		}
	case 36:
		//line synparse.y:135
		{
			trace("<>")
			stack.push(Node{cmd: NEOP, str: "<>", val: 2})
		}
	case 37:
		//line synparse.y:140
		{
			trace("num")
			stack.push(yyS[yypt-0].node)
		}
	case 38:
		//line synparse.y:141
		{
		}
	case 39:
		//line synparse.y:142
		{
			trace("+")
			stack.push(yyS[yypt-1].node)
		}
	case 40:
		//line synparse.y:143
		{
			trace("-")
			stack.push(yyS[yypt-1].node)
		}
	case 41:
		//line synparse.y:144
		{
			trace("*")
			stack.push(yyS[yypt-1].node)
		}
	case 42:
		//line synparse.y:145
		{
			trace("/")
			stack.push(yyS[yypt-1].node)
		}
	case 43:
		//line synparse.y:146
		{
			trace("-")
			yyS[yypt-1].node.cmd = UNARYMINUS
			yyS[yypt-1].node.val = 1
			yyS[yypt-1].node.priority = 2
			stack.push(yyS[yypt-1].node)
		}
	case 44:
		//line synparse.y:148
		{
			trace("+")
			yyS[yypt-1].node.cmd = UNARYPLUS
			yyS[yypt-1].node.val = 1
			yyS[yypt-1].node.priority = 2
			stack.push(yyS[yypt-1].node)
		}
	case 46:
		//line synparse.y:154
		{
			trace("num")
			stack.push(yyS[yypt-0].node)
		}
	case 48:
		//line synparse.y:156
		{
			trace("abs")
			yyS[yypt-3].node.val = 1
			stack.push(yyS[yypt-3].node)
		}
	case 49:
		//line synparse.y:157
		{
			trace("+")
			stack.push(yyS[yypt-1].node)
		}
	case 50:
		//line synparse.y:158
		{
			trace("-")
			stack.push(yyS[yypt-1].node)
		}
	case 51:
		//line synparse.y:159
		{
			trace("*")
			stack.push(yyS[yypt-1].node)
		}
	case 52:
		//line synparse.y:160
		{
			trace("/")
			stack.push(yyS[yypt-1].node)
		}
	case 53:
		//line synparse.y:161
		{
			trace("^")
			stack.push(yyS[yypt-1].node)
		}
	case 54:
		//line synparse.y:162
		{
			trace("-")
			yyS[yypt-1].node.cmd = UNARYMINUS
			yyS[yypt-1].node.val = 1
			yyS[yypt-1].node.priority = 2
			stack.push(yyS[yypt-1].node)
		}
	case 55:
		//line synparse.y:164
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
