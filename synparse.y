
%{

package qeconv

import (
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

%token NAME NUMBER F_TRUE F_FALSE
%token ALL EX AND OR NOT
%token PLUS MINUS COMMA MULT DIV POW
%token EOL LB RB LP RP LC RC
%token INDEXED LIST
%token IMPL REPL EQUIV
%token COMMENT

%type <num> seq_var seq_fof
%type <node> NAME NUMBER var
%type <node> ALL EX AND OR NOT IMPL REPL EQUIV
%type <node> PLUS MINUS MULT DIV POW

%left IMPL REPL EQUIV
%left OR
%left AND
%left NOT
%left LTOP GTOP LEOP GEOP NEOP EQOP
%left PLUS MINUS
%left MULT DIV
%left UNARYMINUS UNARYPLUS
%right POW


%%

expr
	: fof EOL
	| list_of_fof EOL
	;


fof
	: ALL LP quantifiers COMMA fof RP { trace("ALL"); stack.push($1)}
	| EX  LP quantifiers COMMA fof RP { trace("EX");  stack.push($1)}
	| AND LP seq_fof RP { trace("and"); $1.val=$3; stack.push($1)}
	| OR  LP seq_fof RP { trace("or"); $1.val=$3; stack.push($1)}
	| NOT fof { trace("not"); stack.push($1)}
	| IMPL LP fof COMMA fof RP { trace("IMPL"); stack.push($1)}
	| REPL LP fof COMMA fof RP { trace("REPL"); stack.push($1)}
	| EQUIV LP fof COMMA fof RP { trace("EQUIV"); stack.push($1)}
	| LP fof RP
	| atom
	;

list_of_fof
	: LB seq_fof RB {
		trace("list")
		stack.push(Node{cmd: LIST, val: $2})
	}
	| LB RB {
		trace("empty-list")
		stack.push(Node{cmd: LIST, val: 0})
	}
	;

seq_fof
	: fof	{ $$ = 1 }
	| seq_fof COMMA fof { $$ = $1 + 1 }
	;

quantifiers
	: var {
		trace("list")
		stack.push(Node{cmd: LIST, val: 1})
	}
	| LB seq_var RB  /* [x,y,z] */ {
		trace("list")
		stack.push(Node{cmd: LIST, val: $2})
	}
	| LC seq_var RC  /* {x,y,z} */ {
		trace("set")
		stack.push(Node{cmd: LIST, val: $2})
	}
	;

seq_var
	: var		{ $$ = 1 }
	| seq_var COMMA var { $$ = $1 + 1 }
	;

var
	: NAME  { trace("name");  stack.push($1)}
	| var LB NUMBER RB  { trace("index");  stack.push(Node{cmd: INDEXED, val:2})}
	;

atom
	: F_TRUE  { trace("true");  stack.push(Node{cmd: F_TRUE, val:0})}
	| F_FALSE { trace("false"); stack.push(Node{cmd: F_FALSE, val:0})}
	| poly LTOP poly { trace("<");  stack.push(Node{cmd: LTOP, str: "<", val:2})}
	| poly GTOP poly { trace(">");  stack.push(Node{cmd: LTOP, str: ">", val:2, rev:true})}
	| poly LEOP poly { trace("<="); stack.push(Node{cmd: LEOP, str: "<=", val:2})}
	| poly GEOP poly { trace(">="); stack.push(Node{cmd: LEOP, str: ">=", val:2, rev:true})}
	| poly EQOP poly { trace("=");  stack.push(Node{cmd: EQOP, str: "=", val:2})}
	| poly NEOP poly { trace("<>"); stack.push(Node{cmd: NEOP, str: "<>", val:2})}
	;


poly
	: LP poly RP
	| NUMBER	{ trace("num"); stack.push($1) }
	| var
	| poly PLUS poly	{ trace("+"); stack.push($2)}
	| poly MINUS poly	{ trace("-"); stack.push($2)}
	| poly MULT poly	{ trace("*"); stack.push($2)}
	| poly DIV poly	{ trace("/"); stack.push($2)}
	| poly POW NUMBER { trace("^"); stack.push($3); stack.push($2)}
	| MINUS poly %prec UNARYMINUS	{ trace("-");
		$1.cmd = UNARYMINUS; $1.val = 1; $1.priority = 2; stack.push($1) }
	| PLUS poly %prec UNARYPLUS	{ trace("+");
		$1.cmd = UNARYPLUS; $1.val = 1; $1.priority = 2; stack.push($1) }
	;

%%      /*  start  of  programs  */


type SynLex struct {
	scanner.Scanner
	s string
	comment []Node
}

type SynLex1 struct {
	val string
	label int
	v int
	argn int // 引数の数
	priority int // 要素に () が必要かを判定するためのフラグ
}

var sones = []SynLex1 {
	{"+", PLUS , '+', 2, 4},
	{"-", MINUS, '-', 2, 4},
	{"*", MULT , '*', 2, 3},
	{"/", DIV  , '/', 2, 3},
	{"^", POW  , '^', 2, 1},
	{"[", LB   , '[', 0, 0},
	{"]", RB   , ']', 0, 0},
	{"{", LC   , '{', 0, 0},
	{"}", RC   , '{', 0, 0},
	{"(", LP   , '(', 0, 0},
	{")", RP   , ')', 0, 0},
	{",", COMMA, ',', 0, 0},
	{":", EOL  , ':', 0, 0},
	{"=", EQOP , '=', 0, 0},
}

var sfuns = []SynLex1 {
	{"And"  , AND    , 0, 0, 1},
	{"Or"   , OR     , 0, 0, 2},
	{"Impl" , IMPL   , 0, 2, 0},
	{"Repl" , REPL   , 0, 2, 0},
	{"Equiv", EQUIV  , 0, 2, 0},
	{"Not"  , NOT    , 0, 1, 0},
	{"All"  , ALL    , 0, 2, 0},
	{"Ex"   , EX     , 0, 2, 0},
	{"true" , F_TRUE , 0, 0, 0},
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
			l.comment = append(l.comment, Node{str:str, lineno: lno})
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


