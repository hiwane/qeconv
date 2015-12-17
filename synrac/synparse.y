
%{

package qeconv

import (
	. "github.com/hiwane/qeconv/def"
	"text/scanner"
	"fmt"
	"errors"
	"strings"
)

var stack *QeStack


%}


%union{
	node QeNode
	num int
}

%token name number f_true f_false
%token all ex and or not abs
%token plus minus comma mult div pow
%token ltop gtop leop geop neop eqop
%token eol lb rb lp rp lc rc
%token indexed list
%token impl repl equiv
%token comment

%type <num> seq_var seq_fof seq_mobj
%type <node> name number rational var
%type <node> all ex and or not impl repl equiv abs
%type <node> ltop gtop leop geop neop eqop
%type <node> plus minus mult div pow
%type <node> lb lc
%type <node> f_true f_false

%left impl repl equiv
%left or
%left and
%left not
%left ltop gtop leop geop neop eqop
%left plus minus
%left mult div
%left unaryminus unaryplus
%right pow


%%

expr
	: mobj eol
	;

mobj : fof
	 | list_of_mobj
	 ;

fof
	: all lp quantifiers comma fof rp { trace("ALL"); stack.Push($1)}
	| ex  lp quantifiers comma fof rp { trace("EX");  stack.Push($1)}
	| and lp seq_fof rp { trace("and"); $1.SetVal($3); stack.Push($1)}
	| and lp rp { trace("and()"); stack.Push(NewQeNodeBool(true, $1.GetLno()))}
	| or  lp seq_fof rp { trace("or"); $1.SetVal($3); stack.Push($1)}
	| or lp rp { trace("or()"); stack.Push(NewQeNodeBool(false, $1.GetLno()))}
	| not fof { trace("not"); stack.Push($1)}
	| impl lp fof comma fof rp { trace("IMPL"); stack.Push($1)}
	| repl lp fof comma fof rp { trace("REPL"); stack.Push($1)}
	| equiv lp fof comma fof rp { trace("EQUIV"); stack.Push($1)}
	| lp fof rp
	| atom
	;

list_of_mobj
	: lb seq_mobj rb {
		trace("list")
		stack.Push(NewQeNodeList($2, $1.GetLno()))
	}
	| lb rb {
		trace("empty-list")
		stack.Push(NewQeNodeList(0, $1.GetLno()))
	}
	;

seq_mobj
	: mobj { $$ = 1}
	| seq_mobj comma mobj { $$ = $1 + 1 }
	;

seq_fof
	: fof	{ $$ = 1 }
	| seq_fof comma fof { $$ = $1 + 1 }
	;

quantifiers
	: var {
		trace("var")
		stack.Push(NewQeNodeList(1,-1))
	}
	| lb seq_var rb  /* [x,y,z] */ {
		trace("list")
		stack.Push(NewQeNodeList($2, $1.GetLno()))
	}
	| lc seq_var rc  /* {x,y,z} */ {
		trace("set")
		stack.Push(NewQeNodeList($2, $1.GetLno()))
	}
	;

seq_var
	: var		{ $$ = 1 }
	| seq_var comma var { $$ = $1 + 1 }
	;

var
	: name  { trace("name");  stack.Push($1)}
	| var lb number rb  { trace("index");  stack.Push(NewQeNode(INDEXED, 2, -1))}
	;

atom
	: f_true  { trace("true");  stack.Push($1)}
	| f_false { trace("false"); stack.Push($1)}
	| poly ltop poly { trace("<");  stack.Push($2)}
	| poly gtop poly { trace(">");  stack.Push($2)}
	| poly leop poly { trace("<="); stack.Push($2)}
	| poly geop poly { trace(">="); stack.Push($2)}
	| poly eqop poly { trace("=");  stack.Push($2)}
	| poly neop poly { trace("<>"); stack.Push($2)}
	;


rational
	: number	{ trace("num"); stack.Push($1) }
	| lp rational rp	{}
	| rational plus rational	{ trace("+"); stack.Push($2)}
	| rational minus rational	{ trace("-"); stack.Push($2)}
	| rational mult rational	{ trace("*"); stack.Push($2)}
	| rational div rational	{ trace("/"); stack.Push($2)}
	| minus rational %prec unaryminus	{ trace("-"); NewQeNodeStr("-.", $1.GetLno()) }
	| plus rational %prec unaryplus	{ trace("+"); NewQeNodeStr("+.", $1.GetLno()) }
	;

poly
	: lp poly rp
	| number	{ trace("num"); stack.Push($1) }
	| var
	| abs lp poly rp    { trace("abs"); stack.Push($1); }
	| poly plus poly	{ trace("+"); stack.Push($2)}
	| poly minus poly	{ trace("-"); stack.Push($2)}
	| poly mult poly	{ trace("*"); stack.Push($2)}
	| poly div poly	{ trace("/"); stack.Push($2)}
	| poly pow rational { trace("^"); stack.Push($2)}
	| minus poly %prec unaryminus	{ trace("-"); stack.Push(NewQeNodeStr("-.", $1.GetLno())) }
	| plus poly %prec unaryplus	{ trace("+"); stack.Push(NewQeNodeStr("+.", $1.GetLno())) }
	;


%%      /*  start  of  programs  */

type SynLex struct {
	scanner.Scanner
	s string
	comment []Comment
	err error
}

type SynLex1 struct {
	val string
	label int
	v int
	argn int // 引数の数
	priority int // 要素に () が必要かを判定するためのフラグ
}

var sones = []SynLex1 {
	{"+", plus , '+', 2, 4},
	{"-", minus, '-', 2, 4},
	{"*", mult , '*', 2, 3},
	{"/", div  , '/', 2, 3},
	{"^", pow  , '^', 2, 1},
	{"[", lb   , '[', 0, 0},
	{"]", rb   , ']', 0, 0},
	{"{", lc   , '{', 0, 0},
	{"}", rc   , '}', 0, 0},
	{"(", lp   , '(', 0, 0},
	{")", rp   , ')', 0, 0},
	{",", comma, ',', 0, 0},
	{":", eol  , ':', 0, 0},
	{"=", eqop , '=', 0, 0},
}

var sfuns = []SynLex1 {
	{"And"  , and    , 0, 0, 1},
	{"Or"   , or     , 0, 0, 2},
	{"Impl" , impl   , 0, 2, 3},
	{"Repl" , repl   , 0, 2, 3},
	{"Equiv", equiv  , 0, 2, 3},
	{"Not"  , not    , 0, 1, 0},
	{"All"  , all    , 0, 2, 0},
	{"Ex"   , ex     , 0, 2, 0},
	{"true" , f_true , 0, 0, 0},
	{"false", f_false, 0, 0, 0},
	{"abs",   abs    , 0, 0, 0},
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
		var ret[] rune
		for isdigit(l.Peek()) {
			ret = append(ret, l.Next())
		}
		lval.node = NewQeNodeNum(string(ret), lno)
		return number
	}

	if isalnum(l.Peek()) || c == '_' {
		var ret[] rune
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


