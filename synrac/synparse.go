//line synparse.y:3
package qeconv

import __yyfmt__ "fmt"

//line synparse.y:4
import (
	"errors"
	"fmt"
	. "github.com/hiwane/qeconv/def"
	"strings"
	"text/scanner"
)

var stack *QeStack

//line synparse.y:20
type yySymType struct {
	yys  int
	node QeNode
	num  int
}

const name = 57346
const number = 57347
const f_true = 57348
const f_false = 57349
const all = 57350
const ex = 57351
const and = 57352
const or = 57353
const not = 57354
const abs = 57355
const plus = 57356
const minus = 57357
const comma = 57358
const mult = 57359
const div = 57360
const pow = 57361
const ltop = 57362
const gtop = 57363
const leop = 57364
const geop = 57365
const neop = 57366
const eqop = 57367
const eol = 57368
const lb = 57369
const rb = 57370
const lp = 57371
const rp = 57372
const lc = 57373
const rc = 57374
const indexed = 57375
const list = 57376
const impl = 57377
const repl = 57378
const equiv = 57379
const comment = 57380
const unaryminus = 57381
const unaryplus = 57382

var yyToknames = []string{
	"name",
	"number",
	"f_true",
	"f_false",
	"all",
	"ex",
	"and",
	"or",
	"not",
	"abs",
	"plus",
	"minus",
	"comma",
	"mult",
	"div",
	"pow",
	"ltop",
	"gtop",
	"leop",
	"geop",
	"neop",
	"eqop",
	"eol",
	"lb",
	"rb",
	"lp",
	"rp",
	"lc",
	"rc",
	"indexed",
	"list",
	"impl",
	"repl",
	"equiv",
	"comment",
	"unaryminus",
	"unaryplus",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line synparse.y:162

/*  start  of  programs  */

type SynLex struct {
	scanner.Scanner
	s       string
	comment []Comment
	err     error
}

type SynLex1 struct {
	val      string
	label    int
	v        int
	argn     int // $B0z?t$N?t(B
	priority int // $BMWAG$K(B () $B$,I,MW$+$rH=Dj$9$k$?$a$N%U%i%0(B
}

var sones = []SynLex1{
	{"+", plus, '+', 2, 4},
	{"-", minus, '-', 2, 4},
	{"*", mult, '*', 2, 3},
	{"/", div, '/', 2, 3},
	{"^", pow, '^', 2, 1},
	{"[", lb, '[', 0, 0},
	{"]", rb, ']', 0, 0},
	{"{", lc, '{', 0, 0},
	{"}", rc, '}', 0, 0},
	{"(", lp, '(', 0, 0},
	{")", rp, ')', 0, 0},
	{",", comma, ',', 0, 0},
	{":", eol, ':', 0, 0},
	{"=", eqop, '=', 0, 0},
}

var sfuns = []SynLex1{
	{"And", and, 0, 0, 1},
	{"Or", or, 0, 0, 2},
	{"Impl", impl, 0, 2, 3},
	{"Repl", repl, 0, 2, 3},
	{"Equiv", equiv, 0, 2, 3},
	{"Not", not, 0, 1, 0},
	{"All", all, 0, 2, 0},
	{"Ex", ex, 0, 2, 0},
	{"true", f_true, 0, 0, 0},
	{"false", f_false, 0, 0, 0},
	{"abs", abs, 0, 0, 0},
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
			l.comment = append(l.comment, NewComment(str, lno))
		}
	}

	lno := l.Pos().Line
	c := int(l.Peek())
	for i := 0; i < len(sones); i++ {
		if sones[i].v == c {
			l.Next()
			lval.node = NewQeNodeStr(sones[i].val, lno)
			return sones[i].label
		}
	}
	if c == '>' {
		l.Next()
		if l.Peek() == '=' {
			l.Next()
			lval.node = NewQeNodeStr(">=", lno)
			return geop
		} else {
			lval.node = NewQeNodeStr(">", lno)
			return gtop
		}
	} else if c == '<' {
		l.Next()
		if l.Peek() == '=' {
			l.Next()
			lval.node = NewQeNodeStr("<=", lno)
			return leop
		} else if l.Peek() == '>' {
			l.Next()
			lval.node = NewQeNodeStr("<>", lno)
			return neop
		} else {
			lval.node = NewQeNodeStr("<", lno)
			return ltop
		}
	}

	if isdigit(l.Peek()) {
		var ret []rune
		for isdigit(l.Peek()) {
			ret = append(ret, l.Next())
		}
		lval.node = NewQeNodeNum(string(ret), lno)
		return number
	}

	if isalnum(l.Peek()) || c == '_' {
		var ret []rune
		for isdigit(l.Peek()) || isletter(l.Peek()) {
			ret = append(ret, l.Next())
		}
		str := string(ret)
		for i := 0; i < len(sfuns); i++ {
			if str == sfuns[i].val {
				lval.node = NewQeNodeStr(str, lno)
				return sfuns[i].label
			}
		}

		lval.node = NewQeNodeStr(str, lno)
		return name
	}

	return int(c)
}

func (l *SynLex) Error(s string) {
	pos := l.Pos()
	if l.err == nil {
		l.err = errors.New(fmt.Sprintf("%s:Error:%s \n", pos.String(), s))
	}
}

func parse(str string) (*QeStack, []Comment, error) {
	l := new(SynLex)
	l.Init(strings.NewReader(str))
	stack = new(QeStack)
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

const yyNprod = 56
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 301

var yyAct = []int{

	20, 3, 82, 2, 45, 46, 130, 47, 48, 49,
	113, 30, 129, 91, 128, 34, 24, 127, 96, 38,
	69, 125, 18, 68, 60, 51, 114, 56, 56, 33,
	62, 62, 97, 65, 66, 67, 35, 45, 46, 57,
	47, 48, 49, 58, 113, 52, 54, 96, 71, 32,
	31, 29, 109, 110, 63, 28, 112, 55, 92, 92,
	70, 95, 72, 73, 74, 75, 76, 77, 78, 79,
	80, 81, 93, 83, 88, 101, 89, 102, 103, 27,
	104, 105, 86, 85, 26, 59, 50, 106, 107, 108,
	25, 49, 111, 124, 104, 105, 115, 84, 116, 100,
	117, 118, 119, 99, 98, 120, 121, 122, 123, 47,
	48, 49, 94, 90, 126, 24, 19, 16, 17, 5,
	6, 7, 8, 9, 21, 23, 22, 45, 46, 87,
	47, 48, 49, 24, 14, 4, 1, 36, 15, 37,
	13, 0, 0, 0, 0, 0, 10, 11, 12, 24,
	19, 16, 17, 5, 6, 7, 8, 9, 21, 23,
	22, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 15, 0, 13, 0, 0, 0, 0, 0,
	10, 11, 12, 24, 19, 16, 17, 5, 6, 7,
	8, 9, 21, 23, 22, 24, 19, 16, 17, 5,
	6, 7, 8, 9, 21, 23, 22, 0, 13, 64,
	0, 0, 0, 0, 10, 11, 12, 0, 0, 0,
	13, 61, 0, 0, 0, 0, 10, 11, 12, 24,
	19, 16, 17, 5, 6, 7, 8, 9, 21, 23,
	22, 45, 46, 0, 47, 48, 49, 39, 40, 41,
	42, 44, 43, 0, 13, 0, 0, 69, 0, 0,
	10, 11, 12, 45, 46, 0, 47, 48, 49, 39,
	40, 41, 42, 44, 43, 24, 19, 0, 0, 0,
	0, 0, 0, 0, 21, 23, 22, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	53,
}
var yyPact = []int{

	145, -1000, 64, -1000, -1000, 55, 50, 26, 22, 225,
	21, 20, 0, 225, -1000, 111, -1000, -1000, 249, -1000,
	59, -4, 271, 271, -1000, -1000, 12, 12, 191, 179,
	-1000, 225, 225, 225, -7, 227, 32, -1000, -1000, 271,
	271, 271, 271, 271, 271, 271, 271, 271, 271, 68,
	124, 271, 72, 271, 72, 97, 59, 129, 129, 96,
	31, -1000, -1000, 2, -1000, 88, 87, 83, -1000, -1000,
	-1000, 145, 113, 113, 113, 113, 113, 113, 92, 92,
	72, 72, -1000, -1000, 68, 68, 68, 24, 23, -10,
	225, 28, 59, -6, 225, -1000, 225, -1000, 225, 225,
	225, -1000, 68, 68, 68, 68, 63, -1000, -1000, -1000,
	-1000, -9, -1000, 129, -1000, -13, -1000, -16, -18, -24,
	77, 77, -1000, -1000, -1000, -1000, 59, -1000, -1000, -1000,
	-1000,
}
var yyPgo = []int{

	0, 13, 24, 137, 2, 0, 136, 3, 1, 135,
	57, 134, 22,
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
	35, 36, 37, 29, -11, 27, 6, 7, -12, 5,
	-5, 13, 15, 14, 4, 26, 29, 29, 29, 29,
	-8, 29, 29, 29, -8, -12, -3, 28, -7, 20,
	21, 22, 23, 25, 24, 14, 15, 17, 18, 19,
	27, 29, -12, 29, -12, -10, -5, 27, 31, -10,
	-2, 30, -8, -2, 30, -8, -8, -8, 30, 30,
	28, 16, -12, -12, -12, -12, -12, -12, -12, -12,
	-12, -12, -4, 5, 29, 15, 14, 5, -12, -12,
	16, -1, -5, -1, 16, 30, 16, 30, 16, 16,
	16, -7, 14, 15, 17, 18, -4, -4, -4, 28,
	30, -8, 28, 16, 32, -8, -8, -8, -8, -8,
	-4, -4, -4, -4, 30, 30, -5, 30, 30, 30,
	30,
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
		//line synparse.y:64
		{
			trace("ALL")
			stack.Push(yyS[yypt-5].node)
		}
	case 5:
		//line synparse.y:65
		{
			trace("EX")
			stack.Push(yyS[yypt-5].node)
		}
	case 6:
		//line synparse.y:66
		{
			trace("and")
			yyS[yypt-3].node.SetVal(yyS[yypt-1].num)
			stack.Push(yyS[yypt-3].node)
		}
	case 7:
		//line synparse.y:67
		{
			trace("and()")
			stack.Push(NewQeNodeBool(true, yyS[yypt-2].node.GetLno()))
		}
	case 8:
		//line synparse.y:68
		{
			trace("or")
			yyS[yypt-3].node.SetVal(yyS[yypt-1].num)
			stack.Push(yyS[yypt-3].node)
		}
	case 9:
		//line synparse.y:69
		{
			trace("or()")
			stack.Push(NewQeNodeBool(false, yyS[yypt-2].node.GetLno()))
		}
	case 10:
		//line synparse.y:70
		{
			trace("not")
			stack.Push(yyS[yypt-1].node)
		}
	case 11:
		//line synparse.y:71
		{
			trace("IMPL")
			stack.Push(yyS[yypt-5].node)
		}
	case 12:
		//line synparse.y:72
		{
			trace("REPL")
			stack.Push(yyS[yypt-5].node)
		}
	case 13:
		//line synparse.y:73
		{
			trace("EQUIV")
			stack.Push(yyS[yypt-5].node)
		}
	case 16:
		//line synparse.y:79
		{
			trace("list")
			stack.Push(NewQeNodeList(yyS[yypt-1].num, yyS[yypt-2].node.GetLno()))
		}
	case 17:
		//line synparse.y:83
		{
			trace("empty-list")
			stack.Push(NewQeNodeList(0, yyS[yypt-1].node.GetLno()))
		}
	case 18:
		//line synparse.y:90
		{
			yyVAL.num = 1
		}
	case 19:
		//line synparse.y:91
		{
			yyVAL.num = yyS[yypt-2].num + 1
		}
	case 20:
		//line synparse.y:95
		{
			yyVAL.num = 1
		}
	case 21:
		//line synparse.y:96
		{
			yyVAL.num = yyS[yypt-2].num + 1
		}
	case 22:
		//line synparse.y:100
		{
			trace("var")
			stack.Push(NewQeNodeList(1, -1))
		}
	case 23:
		//line synparse.y:104
		{
			trace("list")
			stack.Push(NewQeNodeList(yyS[yypt-1].num, yyS[yypt-2].node.GetLno()))
		}
	case 24:
		//line synparse.y:108
		{
			trace("set")
			stack.Push(NewQeNodeList(yyS[yypt-1].num, yyS[yypt-2].node.GetLno()))
		}
	case 25:
		//line synparse.y:115
		{
			yyVAL.num = 1
		}
	case 26:
		//line synparse.y:116
		{
			yyVAL.num = yyS[yypt-2].num + 1
		}
	case 27:
		//line synparse.y:120
		{
			trace("name")
			stack.Push(yyS[yypt-0].node)
		}
	case 28:
		//line synparse.y:121
		{
			trace("index")
			stack.Push(NewQeNode(INDEXED, 2, -1))
		}
	case 29:
		//line synparse.y:125
		{
			trace("true")
			stack.Push(yyS[yypt-0].node)
		}
	case 30:
		//line synparse.y:126
		{
			trace("false")
			stack.Push(yyS[yypt-0].node)
		}
	case 31:
		//line synparse.y:127
		{
			trace("<")
			stack.Push(yyS[yypt-1].node)
		}
	case 32:
		//line synparse.y:128
		{
			trace(">")
			stack.Push(yyS[yypt-1].node)
		}
	case 33:
		//line synparse.y:129
		{
			trace("<=")
			stack.Push(yyS[yypt-1].node)
		}
	case 34:
		//line synparse.y:130
		{
			trace(">=")
			stack.Push(yyS[yypt-1].node)
		}
	case 35:
		//line synparse.y:131
		{
			trace("=")
			stack.Push(yyS[yypt-1].node)
		}
	case 36:
		//line synparse.y:132
		{
			trace("<>")
			stack.Push(yyS[yypt-1].node)
		}
	case 37:
		//line synparse.y:137
		{
			trace("num")
			stack.Push(yyS[yypt-0].node)
		}
	case 38:
		//line synparse.y:138
		{
		}
	case 39:
		//line synparse.y:139
		{
			trace("+")
			stack.Push(yyS[yypt-1].node)
		}
	case 40:
		//line synparse.y:140
		{
			trace("-")
			stack.Push(yyS[yypt-1].node)
		}
	case 41:
		//line synparse.y:141
		{
			trace("*")
			stack.Push(yyS[yypt-1].node)
		}
	case 42:
		//line synparse.y:142
		{
			trace("/")
			stack.Push(yyS[yypt-1].node)
		}
	case 43:
		//line synparse.y:143
		{
			trace("-")
			NewQeNodeStr("-.", yyS[yypt-1].node.GetLno())
		}
	case 44:
		//line synparse.y:144
		{
			trace("+")
			NewQeNodeStr("+.", yyS[yypt-1].node.GetLno())
		}
	case 46:
		//line synparse.y:149
		{
			trace("num")
			stack.Push(yyS[yypt-0].node)
		}
	case 48:
		//line synparse.y:151
		{
			trace("abs")
			stack.Push(yyS[yypt-3].node)
		}
	case 49:
		//line synparse.y:152
		{
			trace("+")
			stack.Push(yyS[yypt-1].node)
		}
	case 50:
		//line synparse.y:153
		{
			trace("-")
			stack.Push(yyS[yypt-1].node)
		}
	case 51:
		//line synparse.y:154
		{
			trace("*")
			stack.Push(yyS[yypt-1].node)
		}
	case 52:
		//line synparse.y:155
		{
			trace("/")
			stack.Push(yyS[yypt-1].node)
		}
	case 53:
		//line synparse.y:156
		{
			trace("^")
			stack.Push(yyS[yypt-1].node)
		}
	case 54:
		//line synparse.y:157
		{
			trace("-")
			stack.Push(NewQeNodeStr("-.", yyS[yypt-1].node.GetLno()))
		}
	case 55:
		//line synparse.y:158
		{
			trace("+")
			stack.Push(NewQeNodeStr("+.", yyS[yypt-1].node.GetLno()))
		}
	}
	goto yystack /* stack new state and value */
}
