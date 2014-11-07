
%{

package synparse

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
}

%}


%union{
	node Node
	num int
}

%token NAME NUMBER F_TRUE F_FALSE
%token ALL EX AND OR
%token PLUS MINUS COMMA MULT DIV
%token EOL LB RB LP RP LC RC
%token INDEXED LIST

%type <num> seq_var seq_fof
%type <node> NAME NUMBER var

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
	: ALL LP quantifiers COMMA fof RP { trace("ALL");  stack.push(Node{cmd:ALL,val:2})}
	| EX  LP quantifiers COMMA fof RP { trace("EX");  stack.push(Node{cmd:EX,val:2})}
	| AND LP seq_fof RP { trace("and"); stack.push(Node{cmd:AND, val:$3})}
	| OR  LP seq_fof RP { trace("or");  stack.push(Node{cmd:OR, val:$3})}
	| NOT fof { trace("not");  stack.push(Node{cmd:NOT,val:1})}
	| IMPL LP fof COMMA fof RP { trace("IMPL");  stack.push(Node{cmd:IMPL,val:2})}
	| REPL LP fof COMMA fof RP { trace("REPL");  stack.push(Node{cmd:REPL,val:2,rev:true})}
	| EQUIV LP fof COMMA fof RP { trace("EQUIV");  stack.push(Node{cmd:EQUIV,val:2})}
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
	| poly PLUS poly	{ trace("+"); stack.push(Node{cmd: PLUS, str: "+", val:2}) }
	| poly MINUS poly	{ trace("-"); stack.push(Node{cmd: MINUS, str: "-", val:2}) }
	| poly MULT poly	{ trace("*"); stack.push(Node{cmd: MULT, str: "*", val:2}) }
	| poly DIV poly	{ trace("/"); stack.push(Node{cmd: DIV, str: "/", val:2}) }
	| poly POW NUMBER { trace("^"); stack.push(Node{cmd: POW, str: "^", val:2}) }
	| MINUS poly %prec UNARYMINUS	{ trace("-"); stack.push(Node{cmd: UNARYMINUS, str: "-", val:1}) }
	| PLUS poly %prec UNARYPLUS	{ trace("+"); stack.push(Node{cmd: UNARYPLUS, str: "+", val:1}) }
	;

%%      /*  start  of  programs  */


type SynLex struct {
	scanner.Scanner
	s string
}

type SynLex1 struct {
	val string
	label int
	v int
}

var sones = []SynLex1 {
	{"+", PLUS, '+'},
	{"-", MINUS, '-'},
	{"*", MULT, '*'},
	{"/", DIV, '/'},
	{"^", POW, '^'},
	{"[", LB, '['},
	{"]", RB, ']'},
	{"{", LC, '{'},
	{"}", RC, '{'},
	{"(", LP, '('},
	{")", RP, ')'},
	{",", COMMA, ','},
	{":", EOL, ':'},
	{"=", EQOP, '='},
}


var sfuns = []SynLex1 {
	{"And", AND, 0},
	{"Or", OR, 0},
	{"Impl", IMPL, 0},
	{"Repl", REPL, 0},
	{"Equiv", EQUIV, 0},
	{"Not", NOT, 0},
	{"All", ALL, 0},
	{"Ex", EX, 0},
	{"true", F_TRUE, 0},
	{"false", F_FALSE, 0},
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
	for isspace(l.Peek()) {
		l.Next()
	}

	c := int(l.Peek())
	for i := 0; i < len(sones); i++ {
		if sones[i].v == c {
			l.Next()
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
		lval.node = Node{cmd: NUMBER, val: 0, str: string(ret)}
		return NUMBER
	}

	if isalnum(l.Peek()) || c == '_' {
		var ret[] rune
		for isdigit(l.Peek()) || isletter(l.Peek()) {
			ret = append(ret, l.Next())
		}
		lval.node = Node{cmd: NAME, val: 0, str: string(ret)}
		for i := 0; i < len(sfuns); i++ {
			if lval.node.str == sfuns[i].val {
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
	fmt.Printf(s + "\n")
}


