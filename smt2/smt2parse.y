
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
var assert_cnt int
var decfun_cnt int

type smt2node struct {
	lno, col int
	typ int
	str string
}

%}


%union{
	node smt2node
	num int
}

%token number symbol keyword string_
%token forall exists let as theory par
%token assert check_sat declare_fun set_info set_logic exit
%token set_option
%token ltop gtop leop geop eqop
%token plus minus mult div
%token not and or

%type <node> number symbol keyword id string_ spec_const
%type <node> plus minus mult div
%type <node> ltop gtop leop geop eqop
%type <node> and or not
%type <node> check_sat
%type <num> s_exp0 term1 var_bind1 sorted_var1

%left impl repl equiv
%left or
%left and
%left not
%left ltop gtop leop geop eqop
%left plus minus
%left mult div
%left unaryminus unaryplus
%right pow


%%

script : commands { trace("eof") }


commands
	: command			{ trace("command") }
	| commands command	{ trace("commands") }
	;


s_exp : spec_const
	  | symbol
	  | keyword
	  | '(' s_exp0 ')'
	 ;


s_exp0 :              { $$ = 0 }
	   | s_exp0 s_exp { $$ = $1 + 1 }
	   ;

spec_const
	: number { $$ = $1 }
	| string_ { $$ = $1 }
	;

//index
//	: number
//	| symbol
//	;

id : symbol { $$ = $1 }

sort : id	{ if $1.str != "Real" {yylex.Error("unknown sort") }}
//	 | '(' id sort1 ')' 
	 ;

sort1 : sort
	  | sort1 sort
	  ;

attribute_value
	: spec_const
	| symbol
	| '(' s_exp0 ')'
	;

attribute
	: keyword
	| keyword attribute_value ;

attribute1
	: attribute
	| attribute1 attribute
	;

term : spec_const { stack.Push(NewQeNodeNum($1.str, $1.lno)) }
	 | qual_id
	 | '(' qual_id term1 ')'
	 | '(' let '(' var_bind1 ')' term ')'
	 | '(' forall '(' sorted_var1 ')' term ')'
	 | '(' exists '(' sorted_var1 ')' term ')'
//	 | '(' '!' term attribute1 ')'
	  | '(' plus term1 ')'	{
	  	if $3 > 1 {
	  		stack.Push(NewQeNodeStrVal($2.str, $3, $2.lno))
		}}
	  | '(' minus term1 ')'{
	  	if $3 == 1 {
	  		stack.Push(NewQeNodeStr("-.", $2.lno))
		} else {
	  		stack.Push(NewQeNodeStrVal($2.str, $3, $2.lno))
		}}
	  | '(' mult term1 ')'	{ stack.Push(NewQeNodeStrVal($2.str, $3, $2.lno)) }
	  | '(' div term1 ')'	{ stack.Push(NewQeNodeStrVal($2.str, $3, $2.lno)) }
	  | '(' geop term term ')' { stack.Push(NewQeNodeStr($2.str, $2.lno)) }
	  | '(' gtop term term ')' { stack.Push(NewQeNodeStr($2.str, $2.lno)) }
	  | '(' leop term term ')' { stack.Push(NewQeNodeStr($2.str, $2.lno)) }
	  | '(' ltop term term ')' { stack.Push(NewQeNodeStr($2.str, $2.lno)) }
	  | '(' eqop term term ')' { stack.Push(NewQeNodeStr($2.str, $2.lno)) }
	  | '(' not term ')' { stack.Push(NewQeNodeStr("Not", $2.lno)) }
	  | '(' and term1 ')' { stack.Push(NewQeNodeStrVal("And", $3, $2.lno)) }
	  | '(' or term1 ')' { stack.Push(NewQeNodeStrVal("Or", $3, $2.lno)) }
	 ;

term1: term { $$ = 1}
	 | term1 term { $$ = $1 + 1}
	 ;

sorted_var1
	: sorted_var { $$ = 1 }
	| sorted_var1 sorted_var { $$ = $1 + 1 }
	;

var_bind1
	: var_bind { $$ = 1 }
	| var_bind1 var_bind { $$ = $1 + 1 }
	;

sorted_var : '(' symbol sort ')' ;

var_bind   : '(' symbol term ')' ;

qual_id
	: id { stack.Push(NewQeNodeStr($1.str, $1.lno)) }
//	| '(' as id sort ')'
	;


//sort_sym_decl
//	: '(' id number ')'
//	| '(' id number attribute1 ')'
//	;

//meta_spec_const : NUM | DEC | STR ;

//fun_sym_decl
//	: '(' spec_const sort ')'
//	| '(' spec_const sort attribute1 ')'
//	| '(' meta_spec_const sort ')'
//	| '(' meta_spec_const sort attribute1 ')'
//	| '(' id sort1 ')'
//	| '(' id sort1 attribute1 ')'
//	;

//par_fun_sym_decl
//	: fun_sym_decl
//	| '(' par '(' symbol1 ')' '(' id sort1 ')' ')'
//	| '(' par '(' symbol1 ')' '(' id sort1 attribute1 ')' ')'
//				 ;
//
//symbol1
//	: symbol
//	| symbol1 symbol
//	;
//
//sort1
//	: sort
//	| sort1 sort
//	;

//theory_attribute : :sort '(' sort_sym_decl+ ')'
//				 | :funs '(' par_fun_sym_decl+ ')'
//				 | :sorts-description str_ing
//				 | :funs-description str_ing
//				 | :definition str_ing
//				 | :values str_ing
//				 | :notes str_ing
//				 | attribute
//				 ;

//theory_decl : '(' theory symbol theory_attribute+ ')'

//fun_dec
//	: '(' symbol '(' ')' sort ')'
//	| '(' symbol '(' sorted_var1 ')' sort ')'
//	;
//
//fun_def
//	: symbol '(' ')' sort term
//	| symbol '(' sorted_var1 ')' sort term
//	;

//prop_literal : symbol | '('not symbol')' ;

command : '(' assert term ')' { assert_cnt += 1 }
		| '(' check_sat ')' { trace("go check-sat"); 
			stack.Push(NewQeNodeStrVal("And", assert_cnt, 0)) 
			for i := 0; i < decfun_cnt; i++ {
				stack.Push(NewQeNodeStr("Ex", 0))
			}
			
		}
		| '(' exit ')'
		| '(' declare_fun symbol '(' ')' sort ')' {
			stack.Push(NewQeNodeStr($3.str, $3.lno))
			stack.Push(NewQeNodeList(1, $3.lno))
			decfun_cnt += 1
		}
		| '(' declare_fun symbol '(' sort1 ')' sort ')' { yylex.Error("unknown declare") }
//		| '(' define_fun fun_def ')'
		| '(' set_info attribute ')'
		| '(' set_logic symbol ')' { if $3.str != "QF_NRA" && $3.str != "NRA" { yylex.Error("unknown logic") }}
//		| '(' set_option option ')'
		;


%%      /*  start  of  programs  */

type synLex struct {
	scanner.Scanner
	s string
	comment []Comment
	err error
}

type smt2_lext struct {
	val string
	label int
}

var keywords_tbl = []smt2_lext {
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
		var ret[] rune
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

		lval.node = smt2node{lno,col,number, string(ret)}
		return number
	}

	if issimplsym(l.Peek()) {
		var ret[] rune
		for issimplsym(l.Peek()) {
			ret = append(ret, l.Next())
		}
		str := string(ret)

		for i := 0; i < len(keywords_tbl); i++ {
			if keywords_tbl[i].val == str {
				lval.node = smt2node{lno,col, keywords_tbl[i].label, str}
				return keywords_tbl[i].label
			}
		}
		str = strings.Replace(str, "?", "_q_", -1)
		str = strings.Replace(str, "!", "_e_", -1)
		lval.node = smt2node{lno,col, symbol, str}
		return symbol
	}

	if c == ':' {
		var ret[] rune
		ret = append(ret, l.Next())
		for issimplsym(l.Peek()) {
			ret = append(ret, l.Next())
		}
		lval.node = smt2node{lno,col, keyword, string(ret)}
		return keyword
	}
	if c == '|' || c == '"' {
		var ret[] rune
		corg := l.Peek()
		ret = append(ret, l.Next())
		for l.Peek() != corg {
			ret = append(ret, l.Next())
		}
		l.Next()
		if c == '|' {
			lval.node = smt2node{lno,col, symbol, string(ret)}
		} else {
			lval.node = smt2node{lno,col, string_, string(ret)}
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


