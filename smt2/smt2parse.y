
%{

package qeconv

import (
	. "github.com/hiwane/qeconv/def"
	"text/scanner"
	"errors"
	"fmt"
	"strings"
)

var stack *QeStack
var assert_cnt int
var decfun_cnt int


var letmap smt2letdat


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
%token kw_status
%token forall exists let as theory par
%token assert check_sat declare_const declare_fun set_info set_logic exit
%token set_option
%token ltop gtop leop geop eqop
%token plus minus mult div
%token not and or implies
%token lp rp

%type <node> number symbol keyword id string_ spec_const
%type <node> forall exists let
%type <node> plus minus mult div
%type <node> ltop gtop leop geop eqop
%type <node> and or not implies
%type <node> check_sat
%type <num> s_exp0 term1 var_bind1 sorted_var1
%type <num> lp rp

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
	  | lp s_exp0 rp
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

sort : id	{ if $1.str != "Real" {yylex.Error("unknown sort: " + $1.str) }}
//	 | lp id sort1 rp 
	 ;

sort1 : sort
	  | sort1 sort
	  ;

attribute_value
	: spec_const
	| symbol
	| lp s_exp0 rp
	;

attribute
	: keyword
	| kw_status symbol {
		if l, ok := yylex.(commentI); ok {
			l.append_comment(":status " + $2.str, $2.lno)
		} }
	| keyword attribute_value ;

attribute1
	: attribute
	| attribute1 attribute
	;

term : spec_const { stack.Push(NewQeNodeNum($1.str, $1.lno)) }
	 | qual_id
	 | lp qual_id term1 rp
	 | lp let lp var_bind1 rp term rp {
	 	letmap.popn($4)
	 }
	 | lp forall lp quantifiers rp term rp {stack.Push(NewQeNodeStr("All", $2.lno))}
	 | lp exists lp quantifiers rp term rp {stack.Push(NewQeNodeStr("Ex", $2.lno))}
//	 | lp '!' term attribute1 rp {yylex.Error("unsupported !")}
	  | lp plus term1 rp	{
	  	if $3 > 1 {
	  		stack.Push(NewQeNodeStrVal($2.str, $3, $2.lno))
		}}
	  | lp minus term1 rp{
	  	if $3 == 1 {
	  		stack.Push(NewQeNodeStr("-.", $2.lno))
		} else {
	  		stack.Push(NewQeNodeStrVal($2.str, $3, $2.lno))
		}}
	  | lp mult term1 rp	{ stack.Push(NewQeNodeStrVal($2.str, $3, $2.lno)) }
	  | lp div term1 rp	{ stack.Push(NewQeNodeStrVal($2.str, $3, $2.lno)) }
	  | lp geop term term rp { stack.Push(NewQeNodeStr($2.str, $2.lno)) }
	  | lp gtop term term rp { stack.Push(NewQeNodeStr($2.str, $2.lno)) }
	  | lp leop term term rp { stack.Push(NewQeNodeStr($2.str, $2.lno)) }
	  | lp ltop term term rp { stack.Push(NewQeNodeStr($2.str, $2.lno)) }
	  | lp eqop term term rp { stack.Push(NewQeNodeStr($2.str, $2.lno)) }
	  | lp not term rp { stack.Push(NewQeNodeStr("Not", $2.lno)) }
	  | lp and term1 rp { stack.Push(NewQeNodeStrVal("And", $3, $2.lno)) }
	  | lp or term1 rp { stack.Push(NewQeNodeStrVal("Or", $3, $2.lno)) }
	  | lp implies term term rp { stack.Push(NewQeNodeStr("Impl", $2.lno)) }
	 ;

term1: term { $$ = 1}
	 | term1 term { $$ = $1 + 1}
	 ;
quantifiers: sorted_var1 {
		   stack.Push(NewQeNodeList($1, 0))
	}

sorted_var1
	: sorted_var { $$ = 1 }
	| sorted_var1 sorted_var { $$ = $1 + 1 }
	;

var_bind1
	: var_bind { $$ = 1 }
	| var_bind1 var_bind { $$ = $1 + 1 }
	;

sorted_var : lp symbol sort rp {
		stack.Push(NewQeNodeStr($2.str, $2.lno))
	};

var_bind   : lp symbol term rp {
		   letmap.update_letmap(stack, $1, $2)
	};

qual_id
	: id {
		v, ok := letmap.get($1.str)
		if ok {
			// letmap の内容を挿入する.
			stack.Pushn(v)
		} else {
			stack.Push(NewQeNodeStr($1.str, $1.lno))
		}
	}
//	| lp as id sort rp
	;


//sort_sym_decl
//	: lp id number rp
//	| lp id number attribute1 rp
//	;

//meta_spec_const : NUM | DEC | STR ;

//fun_sym_decl
//	: lp spec_const sort rp
//	| lp spec_const sort attribute1 rp
//	| lp meta_spec_const sort rp
//	| lp meta_spec_const sort attribute1 rp
//	| lp id sort1 rp
//	| lp id sort1 attribute1 rp
//	;

//par_fun_sym_decl
//	: fun_sym_decl
//	| lp par lp symbol1 rp lp id sort1 rp rp
//	| lp par lp symbol1 rp lp id sort1 attribute1 rp rp
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

//theory_attribute : :sort lp sort_sym_decl+ rp
//				 | :funs lp par_fun_sym_decl+ rp
//				 | :sorts-description str_ing
//				 | :funs-description str_ing
//				 | :definition str_ing
//				 | :values str_ing
//				 | :notes str_ing
//				 | attribute
//				 ;

//theory_decl : lp theory symbol theory_attribute+ rp

//fun_dec
//	: lp symbol lp rp sort rp
//	| lp symbol lp sorted_var1 rp sort rp
//	;
//
//fun_def
//	: symbol lp rp sort term
//	| symbol lp sorted_var1 rp sort term
//	;

//prop_literal : symbol | lpnot symbolrp ;

command : lp assert term rp { assert_cnt += 1 }
		| lp check_sat rp { trace("go check-sat"); 
			stack.Push(NewQeNodeStrVal("And", assert_cnt, 0)) 
			for i := 0; i < decfun_cnt; i++ {
				stack.Push(NewQeNodeStr("Ex", 0))
			}
		}
		| lp exit rp
		| lp declare_fun symbol lp rp sort rp {
			stack.Push(NewQeNodeStr($3.str, $3.lno))
			stack.Push(NewQeNodeList(1, $3.lno))
			decfun_cnt += 1
		}
		| lp declare_fun symbol lp sort1 rp sort rp { yylex.Error("unknown declare") }
		| lp declare_const symbol sort rp {
			stack.Push(NewQeNodeStr($3.str, $3.lno))
			stack.Push(NewQeNodeList(1, $3.lno))
			decfun_cnt += 1
		}
//		| lp define_fun fun_def rp
		| lp set_info attribute rp
		| lp set_logic symbol rp { if $3.str != "QF_NRA" && $3.str != "NRA" { yylex.Error("unknown logic") }}
//		| lp set_option option rp
		;


%%      /*  start  of  programs  */

type commentI interface {
	append_comment(comm string, lno int)
}


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
	{"declare-const", declare_const},
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
	{"=>", implies},
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

func (l *synLex) append_comment(comm string, lno int) {
	l.comment = append(l.comment, NewComment(comm, lno))
}

func (l *synLex) Lex(lval *yySymType) int {

	for {
		// skip space
		for isspace(l.Peek()) {
			l.Next()
		}
		if scanner.EOF == l.Peek() {
			trace("Lex! eof " + l.Pos().String())
			l.Next()
			return 0
		}
		c := int(l.Peek())
		if c != ';' {
			break
		}
		// comment
		l.Next()
		str := ""
		for l.Peek() != '\n' {
			str += string(l.Next())
		}
		if str != "" {
			lno := l.Pos().Line
			l.append_comment(str, lno)
		}
	}
	trace("Lex! " + string(l.Peek()) + l.Pos().String())

	c := int(l.Peek())
	lno := l.Pos().Line
	col := l.Pos().Column
	if c == '(' || c == ')' {
		l.Next()
		lval.num = stack.Length()
		if c == '(' {
			return lp
		} else {
			return rp
		}
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
		str = strings.Replace(str, ".", "_d_", -1)
		lval.node = smt2node{lno,col, symbol, str}
		return symbol
	}

	if c == ':' {
		var ret[] rune
		ret = append(ret, l.Next())
		for issimplsym(l.Peek()) {
			ret = append(ret, l.Next())
		}
		str := string(ret)
		if str == ":status" {
			return kw_status
		}
		lval.node = smt2node{lno,col, keyword, str}
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
	letmap.reset()
	yyParse(l)
	return stack, l.comment, l.err
}

func trace(s string) {
//	fmt.Printf(s + "\n")
}


