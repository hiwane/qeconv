
%{

package qeconv

import (
	. "github.com/hiwane/qeconv/def"
)

var stack *QeStack
var assert_stk ex_andStack

var letmap smt2letdat

type smt2node struct {
	lno, col int
	typ int
	str, org_str string
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
		if l, ok := yylex.(commentI); ok {
			l.push_symbol($2)
		}
	};

var_bind   : lp symbol term rp {
		   letmap.update_letmap(stack, $1, $2)
	};

qual_id
	: id {
		if l, ok := yylex.(commentI); ok {
			l.push_symbol($1)
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

command : lp assert term rp {
			assert_stk.assert()
		}
		| lp check_sat rp { trace("go check-sat"); 
			assert_stk.check_sat()
		}
		| lp exit rp
		| lp declare_fun symbol lp rp sort rp {
			assert_stk.declare_sym($3, yylex)
		}
		| lp declare_fun symbol lp sort1 rp sort rp { yylex.Error("unknown declare") }
		| lp declare_const symbol sort rp {
			assert_stk.declare_sym($3, yylex)
		}
//		| lp define_fun fun_def rp
		| lp set_info attribute rp
		| lp set_logic symbol rp { if $3.str != "QF_NRA" && $3.str != "NRA" { yylex.Error("unknown logic: " + $3.str) }}
//		| lp set_option option rp
		;


%%      /*  start  of  programs  */

