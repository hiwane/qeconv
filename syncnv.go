
package qeconv

import (
	"errors"
)

type CnvInf interface {

	/* quantifier */
	All(f Formula) string
	Ex(f Formula) string

	/* logical operator
	 * Mathematica: Implies/Equivalent
	 * redlog     : impl/repl/equiv
	 * qepcad     : =>/<=/<=>
	 */
	And(f Formula) string
	Or(f Formula) string
	Not(f Formula) string
	Impl(f Formula) string
	Equiv(f Formula) string

	/* comparator */
	Leop(f Formula) string
	Ltop(f Formula) string
	Eqop(f Formula) string
	Neop(f Formula) string

	/* atom */
	Ftrue() string
	Ffalse() string

	/* math op */

	List(f Formula) string
}


type Formula struct {
	cmd int
	str string
	args []Formula
	priority int
}

func tofml(s *Stack) Formula {
	n := s.pop()
	fml := Formula{}
	fml.cmd = n.cmd
	fml.str = n.str
	fml.priority = n.priority
	fml.args = make([]Formula, n.val)
	if (fml.cmd == OR || fml.cmd == AND) && n.val == 1 {
		return tofml(s)
	}
	if n.rev {
		for i := 0; i < n.val; i++ {
			fml.args[i] = tofml(s)
		}
	} else {
		for i := 0; i < n.val; i++ {
			fml.args[n.val - i - 1] = tofml(s)
		}
	}
	return fml
}

func mop(fml Formula, cinf CnvInf, op string) string {
	ret := ""
	for i := 0; i < len(fml.args); i++ {
		if i != 0 {
			ret += op
		}
		if fml.priority > 0 && fml.priority < fml.args[i].priority {
			ret += "("
			ret += conv(fml.args[i], cinf)
			ret += ")"
		} else {
			ret += conv(fml.args[i], cinf)
		}
	}
	return ret
}
func uniop(fml Formula, cinf CnvInf, op string) string {
	ret := op
	ret += conv(fml.args[0], cinf)
	return ret
}

func conv(fml Formula, cinf CnvInf) string {
	ret := ""
	switch fml.cmd {
	case ALL:
		ret += cinf.All(fml)
	case EX:
		ret += cinf.Ex(fml)
	case AND:
		ret += cinf.And(fml)
	case OR:
		ret += cinf.Or(fml)
	case NOT:
		ret += cinf.Not(fml)
	case IMPL:
		ret += cinf.Impl(fml)
	case EQUIV:
		ret += cinf.Equiv(fml)
	case LEOP:
		ret += cinf.Leop(fml)
	case LTOP:
		ret += cinf.Ltop(fml)
	case EQOP:
		ret += cinf.Eqop(fml)
	case NEOP:
		ret += cinf.Neop(fml)
	case PLUS:
		ret += mop(fml, cinf, "+")
	case MINUS:
		ret += mop(fml, cinf, "-")
	case MULT:
		ret += mop(fml, cinf, "*")
	case DIV:
		ret += mop(fml, cinf, "/")
	case POW:
		ret += mop(fml, cinf, "^")
	case NAME:
		ret += fml.str
	case NUMBER:
		ret += fml.str
	case F_TRUE:
		ret += cinf.Ftrue()
	case F_FALSE:
		ret += cinf.Ffalse()
	case UNARYMINUS:
		ret += uniop(fml, cinf, "-")
	case UNARYPLUS:
		ret += uniop(fml, cinf, "+")
	default:
		errors.New("unknown type")
	}

	return ret
}

func prefix(fml Formula, cinf CnvInf, left, right string) string {
	ret := left
	sep := ""
	for i := 0; i < len(fml.args); i++ {
		ret += sep
		ret += conv(fml.args[i], cinf)
		sep = ","
	}
	ret += right
	return ret
}

func infix(fml Formula, cinf CnvInf, op string) string {
	ret := ""
	sep := ""
	for i := 0; i < len(fml.args); i++ {
		ret += sep
		if fml.priority > 0 && fml.priority < fml.args[i].priority {
			ret += "("
			ret += conv(fml.args[i], cinf)
			ret += ")"
		} else {
			ret += conv(fml.args[i], cinf)
		}
		sep = op
	}
	return ret
}

