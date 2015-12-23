//line smt2parse.y:3
package qeconv

import __yyfmt__ "fmt"

//line smt2parse.y:4
import (
	. "github.com/hiwane/qeconv/def"
)

var stack *QeStack
var assert_cnt int
var decfun_cnt int
var symbol_cnt int
var symbol_map map[string]string

var letmap smt2letdat

type smt2node struct {
	lno, col     int
	typ          int
	str, org_str string
}

//line smt2parse.y:29
type yySymType struct {
	yys  int
	node smt2node
	num  int
}

const number = 57346
const symbol = 57347
const keyword = 57348
const string_ = 57349
const kw_status = 57350
const forall = 57351
const exists = 57352
const let = 57353
const as = 57354
const theory = 57355
const par = 57356
const assert = 57357
const check_sat = 57358
const declare_const = 57359
const declare_fun = 57360
const set_info = 57361
const set_logic = 57362
const exit = 57363
const set_option = 57364
const ltop = 57365
const gtop = 57366
const leop = 57367
const geop = 57368
const eqop = 57369
const plus = 57370
const minus = 57371
const mult = 57372
const div = 57373
const not = 57374
const and = 57375
const or = 57376
const implies = 57377
const lp = 57378
const rp = 57379
const impl = 57380
const repl = 57381
const equiv = 57382
const unaryminus = 57383
const unaryplus = 57384
const pow = 57385

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"number",
	"symbol",
	"keyword",
	"string_",
	"kw_status",
	"forall",
	"exists",
	"let",
	"as",
	"theory",
	"par",
	"assert",
	"check_sat",
	"declare_const",
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
	"lp",
	"rp",
	"impl",
	"repl",
	"equiv",
	"unaryminus",
	"unaryplus",
	"pow",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line smt2parse.y:278

/*  start  of  programs  */

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 61
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 202

var yyAct = [...]int{

	58, 14, 19, 87, 83, 79, 48, 13, 84, 112,
	85, 139, 138, 17, 109, 110, 18, 20, 137, 17,
	109, 110, 18, 17, 20, 20, 18, 49, 52, 17,
	20, 136, 18, 135, 133, 125, 124, 123, 122, 66,
	67, 68, 69, 70, 71, 111, 134, 74, 121, 105,
	49, 111, 107, 120, 77, 16, 57, 75, 81, 119,
	118, 16, 101, 81, 81, 81, 81, 94, 95, 96,
	97, 98, 89, 81, 81, 102, 115, 99, 49, 49,
	78, 108, 103, 104, 56, 50, 29, 113, 22, 21,
	116, 62, 63, 64, 65, 88, 17, 53, 84, 18,
	61, 72, 73, 17, 20, 60, 18, 59, 49, 47,
	4, 20, 126, 128, 117, 129, 130, 127, 114, 132,
	49, 20, 55, 28, 131, 32, 33, 31, 54, 108,
	17, 20, 24, 18, 23, 16, 100, 15, 25, 41,
	39, 40, 38, 42, 34, 35, 36, 37, 43, 44,
	45, 46, 17, 20, 30, 18, 17, 20, 3, 18,
	51, 5, 16, 93, 17, 20, 76, 18, 17, 20,
	106, 18, 6, 7, 10, 9, 11, 12, 8, 26,
	2, 27, 1, 86, 16, 92, 82, 0, 16, 91,
	0, 0, 0, 0, 0, 0, 16, 90, 0, 0,
	16, 80,
}
var yyPact = [...]int{

	74, -1000, 74, -1000, 157, -1000, 19, 52, 51, 129,
	127, 173, 118, 49, -1000, -1000, 116, -1000, -1000, -1000,
	-1000, -1000, -1000, 73, 106, 48, 92, 117, 47, -1000,
	19, 71, 69, 64, 19, 19, 19, 19, 19, 19,
	19, 19, 19, 19, 19, 19, 19, 20, 43, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 164, -1000, 62,
	59, 59, 160, 152, 148, 126, 19, 19, 19, 19,
	19, 40, 99, 25, 19, 106, 12, -1000, -1000, 15,
	-1000, -1000, -28, -1000, 113, 39, 59, -1000, 109, 23,
	-1000, -1000, -1000, -1000, 22, 16, 11, 1, 0, -1000,
	-1000, -1000, -1, -2, -1000, 106, -1000, -1000, -1000, -1000,
	-1000, -1000, 19, -1000, 19, 19, -1000, 106, 19, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -3, 9, -4, -6,
	-19, -25, -26, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 2, 1, 5, 56, 186, 183, 182, 180, 158,
	170, 6, 166, 160, 138, 138, 0, 137, 10, 3,
	4,
}
var yyR1 = [...]int{

	0, 7, 8, 8, 10, 10, 10, 10, 3, 3,
	2, 2, 1, 11, 12, 12, 13, 13, 13, 14,
	14, 14, 15, 15, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 4, 4, 18, 6, 6, 5, 5,
	19, 20, 17, 9, 9, 9, 9, 9, 9, 9,
	9,
}
var yyR2 = [...]int{

	0, 1, 1, 2, 1, 1, 1, 3, 0, 2,
	1, 1, 1, 1, 1, 2, 1, 1, 3, 1,
	2, 2, 1, 2, 1, 1, 4, 7, 7, 7,
	4, 4, 4, 4, 5, 5, 5, 5, 5, 4,
	4, 4, 5, 1, 2, 1, 1, 2, 1, 2,
	4, 4, 1, 4, 3, 3, 7, 8, 5, 4,
	4,
}
var yyChk = [...]int{

	-1000, -7, -8, -9, 36, -9, 15, 16, 21, 18,
	17, 19, 20, -16, -2, -17, 36, 4, 7, -1,
	5, 37, 37, 5, 5, -14, 6, 8, 5, 37,
	-17, 11, 9, 10, 28, 29, 30, 31, 26, 24,
	25, 23, 27, 32, 33, 34, 35, 36, -11, -1,
	37, -13, -2, 5, 36, 5, 37, -4, -16, 36,
	36, 36, -4, -4, -4, -4, -16, -16, -16, -16,
	-16, -16, -4, -4, -16, 37, -12, -11, 37, -3,
	37, -16, -5, -20, 36, -18, -6, -19, 36, -18,
	37, 37, 37, 37, -16, -16, -16, -16, -16, 37,
	37, 37, -16, -11, -11, 37, -10, 37, -2, 5,
	6, 36, 37, -20, 5, 37, -19, 5, 37, 37,
	37, 37, 37, 37, 37, 37, -11, -3, -16, -16,
	-16, -11, -16, 37, 37, 37, 37, 37, 37, 37,
}
var yyDef = [...]int{

	0, -2, 1, 2, 0, 3, 0, 0, 0, 0,
	0, 0, 0, 0, 24, 25, 0, 10, 11, 52,
	12, 54, 55, 0, 0, 0, 19, 0, 0, 53,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 13,
	59, 21, 16, 17, 8, 20, 60, 0, 43, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 14, 58, 0,
	26, 44, 0, 48, 0, 0, 45, 46, 0, 0,
	30, 31, 32, 33, 0, 0, 0, 0, 0, 39,
	40, 41, 0, 0, 15, 0, 9, 18, 4, 5,
	6, 8, 0, 49, 0, 0, 47, 0, 0, 34,
	35, 36, 37, 38, 42, 56, 0, 0, 0, 0,
	0, 0, 0, 57, 7, 27, 51, 28, 50, 29,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lookahead func() int
}

func (p *yyParserImpl) Lookahead() int {
	return p.lookahead()
}

func yyNewParser() yyParser {
	p := &yyParserImpl{
		lookahead: func() int { return -1 },
	}
	return p
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
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

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yytoken := -1 // yychar translated into internal numbering
	yyrcvr.lookahead = func() int { return yychar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yychar = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
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
		yychar, yytoken = yylex1(yylex, &yylval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yychar = -1
		yytoken = -1
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
			yychar, yytoken = yylex1(yylex, &yylval)
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
			if yyn < 0 || yyn == yytoken {
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
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
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
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yychar = -1
			yytoken = -1
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
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
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
		yyDollar = yyS[yypt-1 : yypt+1]
		//line smt2parse.y:66
		{
			trace("eof")
		}
	case 2:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line smt2parse.y:70
		{
			trace("command")
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line smt2parse.y:71
		{
			trace("commands")
		}
	case 8:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line smt2parse.y:82
		{
			yyVAL.num = 0
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line smt2parse.y:83
		{
			yyVAL.num = yyDollar[1].num + 1
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line smt2parse.y:87
		{
			yyVAL.node = yyDollar[1].node
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line smt2parse.y:88
		{
			yyVAL.node = yyDollar[1].node
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line smt2parse.y:96
		{
			yyVAL.node = yyDollar[1].node
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line smt2parse.y:98
		{
			if yyDollar[1].node.str != "Real" {
				yylex.Error("unknown sort: " + yyDollar[1].node.str)
			}
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line smt2parse.y:114
		{
			if l, ok := yylex.(commentI); ok {
				l.append_comment(":status "+yyDollar[2].node.str, yyDollar[2].node.lno)
			}
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line smt2parse.y:125
		{
			stack.Push(NewQeNodeNum(yyDollar[1].node.str, yyDollar[1].node.lno))
		}
	case 27:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line smt2parse.y:128
		{
			letmap.popn(yyDollar[4].num)
		}
	case 28:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line smt2parse.y:131
		{
			stack.Push(NewQeNodeStr("All", yyDollar[2].node.lno))
		}
	case 29:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line smt2parse.y:132
		{
			stack.Push(NewQeNodeStr("Ex", yyDollar[2].node.lno))
		}
	case 30:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line smt2parse.y:134
		{
			if yyDollar[3].num > 1 {
				stack.Push(NewQeNodeStrVal(yyDollar[2].node.str, yyDollar[3].num, yyDollar[2].node.lno))
			}
		}
	case 31:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line smt2parse.y:138
		{
			if yyDollar[3].num == 1 {
				stack.Push(NewQeNodeStr("-.", yyDollar[2].node.lno))
			} else {
				stack.Push(NewQeNodeStrVal(yyDollar[2].node.str, yyDollar[3].num, yyDollar[2].node.lno))
			}
		}
	case 32:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line smt2parse.y:144
		{
			stack.Push(NewQeNodeStrVal(yyDollar[2].node.str, yyDollar[3].num, yyDollar[2].node.lno))
		}
	case 33:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line smt2parse.y:145
		{
			stack.Push(NewQeNodeStrVal(yyDollar[2].node.str, yyDollar[3].num, yyDollar[2].node.lno))
		}
	case 34:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line smt2parse.y:146
		{
			stack.Push(NewQeNodeStr(yyDollar[2].node.str, yyDollar[2].node.lno))
		}
	case 35:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line smt2parse.y:147
		{
			stack.Push(NewQeNodeStr(yyDollar[2].node.str, yyDollar[2].node.lno))
		}
	case 36:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line smt2parse.y:148
		{
			stack.Push(NewQeNodeStr(yyDollar[2].node.str, yyDollar[2].node.lno))
		}
	case 37:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line smt2parse.y:149
		{
			stack.Push(NewQeNodeStr(yyDollar[2].node.str, yyDollar[2].node.lno))
		}
	case 38:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line smt2parse.y:150
		{
			stack.Push(NewQeNodeStr(yyDollar[2].node.str, yyDollar[2].node.lno))
		}
	case 39:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line smt2parse.y:151
		{
			stack.Push(NewQeNodeStr("Not", yyDollar[2].node.lno))
		}
	case 40:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line smt2parse.y:152
		{
			stack.Push(NewQeNodeStrVal("And", yyDollar[3].num, yyDollar[2].node.lno))
		}
	case 41:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line smt2parse.y:153
		{
			stack.Push(NewQeNodeStrVal("Or", yyDollar[3].num, yyDollar[2].node.lno))
		}
	case 42:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line smt2parse.y:154
		{
			stack.Push(NewQeNodeStr("Impl", yyDollar[2].node.lno))
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line smt2parse.y:157
		{
			yyVAL.num = 1
		}
	case 44:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line smt2parse.y:158
		{
			yyVAL.num = yyDollar[1].num + 1
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line smt2parse.y:160
		{
			stack.Push(NewQeNodeList(yyDollar[1].num, 0))
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line smt2parse.y:165
		{
			yyVAL.num = 1
		}
	case 47:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line smt2parse.y:166
		{
			yyVAL.num = yyDollar[1].num + 1
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line smt2parse.y:170
		{
			yyVAL.num = 1
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line smt2parse.y:171
		{
			yyVAL.num = yyDollar[1].num + 1
		}
	case 50:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line smt2parse.y:174
		{
			stack.Push(NewQeNodeStr(yyDollar[2].node.str, yyDollar[2].node.lno))
		}
	case 51:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line smt2parse.y:178
		{
			letmap.update_letmap(stack, yyDollar[1].num, yyDollar[2].node)
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line smt2parse.y:183
		{
			v, ok := letmap.get(yyDollar[1].node.str)
			if ok {
				// letmap $B$NFbMF$rA^F~$9$k(B.
				stack.Pushn(v)
			} else {
				stack.Push(NewQeNodeStr(yyDollar[1].node.str, yyDollar[1].node.lno))
			}
		}
	case 53:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line smt2parse.y:252
		{
			assert_cnt += 1
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line smt2parse.y:253
		{
			trace("go check-sat")
			stack.Push(NewQeNodeStrVal("And", assert_cnt, 0))
			for i := 0; i < decfun_cnt; i++ {
				stack.Push(NewQeNodeStr("Ex", 0))
			}
		}
	case 56:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line smt2parse.y:260
		{
			stack.Push(NewQeNodeStr(yyDollar[3].node.str, yyDollar[3].node.lno))
			stack.Push(NewQeNodeList(1, yyDollar[3].node.lno))
			decfun_cnt += 1
		}
	case 57:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line smt2parse.y:265
		{
			yylex.Error("unknown declare")
		}
	case 58:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line smt2parse.y:266
		{
			stack.Push(NewQeNodeStr(yyDollar[3].node.str, yyDollar[3].node.lno))
			stack.Push(NewQeNodeList(1, yyDollar[3].node.lno))
			decfun_cnt += 1
		}
	case 60:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line smt2parse.y:273
		{
			if yyDollar[3].node.str != "QF_NRA" && yyDollar[3].node.str != "NRA" {
				yylex.Error("unknown logic: " + yyDollar[3].node.str)
			}
		}
	}
	goto yystack /* stack new state and value */
}
