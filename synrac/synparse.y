
%{

package qeconv

import (
	. "github.com/hiwane/qeconv/def"
	"text/scanner"
	"fmt"
)

var stack *Stack


type Node struct {
	cmd int
	val int
	str string
	rev bool
	priority int
	lineno int
}

%}


%union{
	node Node
	num int
}

%token name number f_true f_false
%token all ex and or not abs
%token plus minus comma mult div pow
%token eol lb rb lp rp lc rc
%token indexed list
%token impl repl equiv
%token comment

%type <num> seq_var seq_fof seq_mobj
%type <node> name number rational var
%type <node> all ex and or not impl repl equiv abs
%type <node> plus minus mult div pow
%type <node> lb lc

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
	: all lp quantifiers comma fof rp { trace("ALL"); stack.push($1)}
	| ex  lp quantifiers comma fof rp { trace("EX");  stack.push($1)}
	| and lp seq_fof rp { trace("and"); $1.val=$3; stack.push($1)}
	| and lp rp { trace("and()"); stack.push(Node{cmd: F_TRUE, val:0})}
	| or  lp seq_fof rp { trace("or"); $1.val=$3; stack.push($1)}
	| or lp rp { trace("or()"); stack.push(Node{cmd: F_FALSE, val:0})}
	| not fof { trace("not"); stack.push($1)}
	| impl lp fof comma fof rp { trace("IMPL"); stack.push($1)}
	| repl lp fof comma fof rp { trace("REPL"); stack.push($1)}
	| equiv lp fof comma fof rp { trace("EQUIV"); stack.push($1)}
	| lp fof rp
	| atom
	;

list_of_mobj
	: lb seq_mobj rb {
		trace("list")
		stack.push(Node{cmd: LIST, val: $2, lineno: $1.lineno})
	}
	| lb rb {
		trace("empty-list")
		stack.push(Node{cmd: LIST, val: 0, lineno: $1.lineno})
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
		stack.push(Node{cmd: LIST, val: 1})
	}
	| lb seq_var rb  /* [x,y,z] */ {
		trace("list")
		stack.push(Node{cmd: LIST, val: $2, lineno: $1.lineno})
	}
	| lc seq_var rc  /* {x,y,z} */ {
		trace("set")
		stack.push(Node{cmd: LIST, val: $2, lineno: $1.lineno})
	}
	;

seq_var
	: var		{ $$ = 1 }
	| seq_var comma var { $$ = $1 + 1 }
	;

var
	: name  { trace("name");  stack.push($1)}
	| var lb number rb  { trace("index");  stack.push(Node{cmd: INDEXED, val:2})}
	;

atom
	: f_true  { trace("true");  stack.push(Node{cmd: F_TRUE, val:0})}
	| f_false { trace("false"); stack.push(Node{cmd: F_FALSE, val:0})}
	| poly ltop poly { trace("<");  stack.push(Node{cmd: LTOP, str: "<", val:2})}
	| poly gtop poly { trace(">");  stack.push(Node{cmd: LTOP, str: ">", val:2, rev:true})}
	| poly leop poly { trace("<="); stack.push(Node{cmd: LEOP, str: "<=", val:2})}
	| poly geop poly { trace(">="); stack.push(Node{cmd: LEOP, str: ">=", val:2, rev:true})}
	| poly eqop poly { trace("=");  stack.push(Node{cmd: EQOP, str: "=", val:2})}
	| poly neop poly { trace("<>"); stack.push(Node{cmd: NEOP, str: "<>", val:2})}
	;


rational
	: number	{ trace("num"); stack.push($1) }
	| lp rational rp	{}
	| rational plus rational	{ trace("+"); stack.push($2)}
	| rational minus rational	{ trace("-"); stack.push($2)}
	| rational mult rational	{ trace("*"); stack.push($2)}
	| rational div rational	{ trace("/"); stack.push($2)}
	| minus rational %prec unaryminus	{ trace("-");
		$1.cmd = UNARYMINUS; $1.val = 1; $1.priority = 2; stack.push($1) }
	| plus rational %prec unaryplus	{ trace("+");
		$1.cmd = UNARYPLUS; $1.val = 1; $1.priority = 2; stack.push($1) }
	;

poly
	: lp poly rp
	| number	{ trace("num"); stack.push($1) }
	| var
	| abs lp poly rp    { trace("abs"); $1.val = 1; stack.push($1); }
	| poly plus poly	{ trace("+"); stack.push($2)}
	| poly minus poly	{ trace("-"); stack.push($2)}
	| poly mult poly	{ trace("*"); stack.push($2)}
	| poly div poly	{ trace("/"); stack.push($2)}
	| poly pow rational { trace("^"); stack.push($2)}
	| minus poly %prec unaryminus	{ trace("-");
		$1.cmd = UNARYMINUS; $1.val = 1; $1.priority = 2; stack.push($1) }
	| plus poly %prec unaryplus	{ trace("+");
		$1.cmd = UNARYPLUS; $1.val = 1; $1.priority = 2; stack.push($1) }
	;


%%      /*  start  of  programs  */

type SynLex struct {
	scanner.Scanner
	s string
	comment []Comment
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
		var ret[] rune
		for isdigit(l.Peek()) {
			ret = append(ret, l.Next())
		}
		lval.node = Node{cmd: NUMBER, val: 0, str: string(ret), lineno: lno}
		return NUMBER
	}

	if isalnum(l.Peek()) || c == '_' {
		var ret[] rune
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


func tofml(s *Stack) Formula {
	n, _ := s.pop()
	fml := NewFormula(n.cmd, n.str, n.lineno, n.priority)
	fml.SetArgLen(n.val)
	if (n.cmd == or || n.cmd == and) && n.val == 1 {
		return tofml(s)
	}
	if n.rev {
		for i := 0; i < n.val; i++ {
			fml.SetArg(i, tofml(s))
		}
	} else {
		for i := 0; i < n.val; i++ {
			fml.SetArg(n.val-i-1, tofml(s))
		}
	}
	return fml
}

