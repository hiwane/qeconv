package qeconv

import (
	"strconv"
)

type QeNode struct {
	cmd      int
	val      int
	str      string
	rev      bool
	priority int
	lineno   int
}

func NewQeNodeBool(b bool, lno int) QeNode {
	if b {
		return QeNode{cmd: F_TRUE, val: 0, lineno: lno, str: "true"}
	} else {
		return QeNode{cmd: F_FALSE, val: 0, lineno: lno, str: "false"}
	}
}

func NewQeNodeStrVal(str string, val, lno int) QeNode {
	q := NewQeNodeStr(str, lno)
	q.SetVal(val)
	return q
}

func NewQeNodeStr(str string, lno int) QeNode {
	if str == "<" {
		return QeNode{cmd: LTOP, str: "<", val: 2}
	} else if str == ">" {
		return QeNode{cmd: LTOP, str: ">", val: 2, rev: true}
	} else if str == "<=" {
		return QeNode{cmd: LEOP, str: "<=", val: 2}
	} else if str == ">=" {
		return QeNode{cmd: LEOP, str: ">=", val: 2, rev: true}
	} else if str == "=" {
		return QeNode{cmd: EQOP, str: ">=", val: 2}
	} else if str == "<>" {
		return QeNode{cmd: NEOP, str: "<>", val: 2}
	} else if str == "+" {
		return QeNode{cmd: PLUS, val: 2, str: str, priority: 4, lineno: lno}
	} else if str == "-" {
		return QeNode{cmd: MINUS, val: 2, str: str, priority: 4, lineno: lno}
	} else if str == "+." {
		return QeNode{cmd: UNARYPLUS, val: 1, str: str, priority: 2, lineno: lno}
	} else if str == "-." {
		return QeNode{cmd: UNARYMINUS, val: 1, str: str, priority: 2, lineno: lno}
	} else if str == "*" {
		return QeNode{cmd: MULT, val: 2, str: str, priority: 3, lineno: lno}
	} else if str == "/" {
		return QeNode{cmd: DIV, val: 2, str: str, priority: 3, lineno: lno}
	} else if str == "^" {
		return QeNode{cmd: POW, val: 2, str: str, priority: 1, lineno: lno}
	} else if str == "[" {
		return QeNode{cmd: LB, str: str, lineno: lno}
	} else if str == "]" {
		return QeNode{cmd: RB, str: str, lineno: lno}
	} else if str == "{" {
		return QeNode{cmd: LC, str: str, lineno: lno}
	} else if str == "}" {
		return QeNode{cmd: RC, str: str, lineno: lno}
	} else if str == "(" {
		return QeNode{cmd: LP, str: str, lineno: lno}
	} else if str == ")" {
		return QeNode{cmd: RP, str: str, lineno: lno}
	} else if str == "," {
		return QeNode{cmd: COMMA, str: str, lineno: lno}
	} else if str == ":" {
		return QeNode{cmd: EOL, str: str, lineno: lno}
	} else if str == "=" {
		return QeNode{cmd: EQOP, str: str, lineno: lno}
	} else if str == "And" {
		return QeNode{cmd: AND, str: str, lineno: lno, priority: 1}
	} else if str == "Or" {
		return QeNode{cmd: OR, str: str, lineno: lno, priority: 2}
	} else if str == "Impl" {
		return QeNode{cmd: IMPL, val: 2, str: str, lineno: lno, priority: 3}
	} else if str == "Repl" {
		return QeNode{cmd: IMPL, val: 2, str: str, lineno: lno, priority: 3, rev: true}
	} else if str == "Equiv" {
		return QeNode{cmd: EQUIV, val: 2, str: str, lineno: lno, priority: 3}
	} else if str == "Not" {
		return QeNode{cmd: NOT, val: 1, str: str, lineno: lno, priority: 0}
	} else if str == "All" {
		return QeNode{cmd: ALL, val: 2, str: str, lineno: lno, priority: 0}
	} else if str == "Ex" {
		return QeNode{cmd: EX, val: 2, str: str, lineno: lno, priority: 0}
	} else if str == "true" {
		return NewQeNodeBool(true, lno)
	} else if str == "false" {
		return NewQeNodeBool(false, lno)
	} else if str == "abs" {
		return QeNode{cmd: ABS, str: str, lineno: lno, val: 1}
	} else {
		return QeNode{cmd: NAME, str: str, lineno: lno}
	}
}

func NewQeNodeList(val, lno int) QeNode {
	if val <= 0 {
		return QeNode{cmd: LIST, val: val, str: "LIST0"}
	} else {
		return QeNode{cmd: LIST, val: val, lineno: lno, str: "LIST" + strconv.Itoa(val)}
	}
}

func NewQeNodeNum(str string, lno int) QeNode {
	// big integer 対応が面倒なので文字列のまま
	return QeNode{cmd: NUMBER, str: str, lineno: lno}
}

func NewQeNode(cmd, val, lno int) QeNode {
	if lno <= 0 {
		return QeNode{cmd: cmd, val: val}
	} else {
		return QeNode{cmd: cmd, val: val, lineno: lno}
	}
}

func (n *QeNode) SetVal(v int) {
	n.val = v
}

func (n *QeNode) GetLno() int {
	return n.lineno
}

func (n *QeNode) String() string {
	return n.str + ":" + strconv.Itoa(n.cmd)
}

func ToFml(s *QeStack) Formula {
	n, _ := s.pop()
	fml := NewFormula(n)
	fml.SetArgLen(n.val)
	if (n.cmd == OR || n.cmd == AND) && n.val == 1 {
		return ToFml(s)
	}
	if n.rev {
		for i := 0; i < n.val; i++ {
			fml.SetArg(i, ToFml(s))
		}
	} else {
		for i := 0; i < n.val; i++ {
			fml.SetArg(n.val-i-1, ToFml(s))
		}
	}
	return fml
}
