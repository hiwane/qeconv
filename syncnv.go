package qeconv

import (
	"errors"
)

type cnv_out struct {
	str string
	lno int
	comment []Node
}

type CnvInf interface {

	Comment(str string) string;

	/* quantifier */
	All(f Formula, co *cnv_out)
	Ex(f Formula, co *cnv_out)

	/* logical operator
	 * Mathematica: Implies/Equivalent
	 * redlog     : impl/repl/equiv
	 * qepcad     : =>/<=/<=>
	 */
	And(f Formula, co *cnv_out)
	Or(f Formula, co *cnv_out)
	Not(f Formula, co *cnv_out)
	Impl(f Formula, co *cnv_out)
	Equiv(f Formula, co *cnv_out)

	/* comparator */
	Leop(f Formula, co *cnv_out)
	Ltop(f Formula, co *cnv_out)
	Eqop(f Formula, co *cnv_out)
	Neop(f Formula, co *cnv_out)

	/* atom */
	Ftrue() string
	Ffalse() string

	/* math op */

	List(f Formula, co *cnv_out)
}

type Formula struct {
	cmd      int
	str      string
	args     []Formula
	priority int
	lineno   int
}

func (c *cnv_out) append(s string) {
	//	fmt.Printf("append [%s]\n", s)
	c.str += s
}

func tofml(s *Stack) Formula {
	n := s.pop()
	fml := Formula{}
	fml.cmd = n.cmd
	fml.str = n.str
	fml.lineno = n.lineno
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
			fml.args[n.val-i-1] = tofml(s)
		}
	}
	return fml
}

func mop(fml Formula, cinf CnvInf, op string, co *cnv_out) {
	for i := 0; i < len(fml.args); i++ {
		if i != 0 {
			co.append(op)
		}
		if fml.priority > 0 && fml.priority < fml.args[i].priority {
			co.append("(")
			conv2(fml.args[i], cinf, co)
			co.append(")")
		} else {
			conv2(fml.args[i], cinf, co)
		}
	}
}
func uniop(fml Formula, cinf CnvInf, op string, co *cnv_out) {
	co.append(op)
	conv2(fml.args[0], cinf, co)
}

func conv(fml Formula, cinf CnvInf, comment []Node) string {
	var co *cnv_out
	co = &cnv_out{str: "", lno: 1, comment: comment}
	conv2(fml, cinf, co)
	return co.str
}

func conv2(fml Formula, cinf CnvInf, co *cnv_out) {
	//	fmt.Printf("fml.cmd=%d,lineno=%d/%d str=%s\n", fml.cmd, fml.lineno, co.lno, fml.str)
	for co.lno < fml.lineno {
		if len(co.comment) > 0 && co.comment[0].lineno == co.lno {
			co.append(cinf.Comment(co.comment[0].str))
			co.comment = co.comment[1: len(co.comment)]
		}
		co.append("\n")
		co.lno++
	}

	switch fml.cmd {
	case ALL:
		cinf.All(fml, co)
	case EX:
		cinf.Ex(fml, co)
	case AND:
		cinf.And(fml, co)
	case OR:
		cinf.Or(fml, co)
	case NOT:
		cinf.Not(fml, co)
	case IMPL:
		cinf.Impl(fml, co)
	case EQUIV:
		cinf.Equiv(fml, co)
	case LEOP:
		cinf.Leop(fml, co)
	case LTOP:
		cinf.Ltop(fml, co)
	case EQOP:
		cinf.Eqop(fml, co)
	case NEOP:
		cinf.Neop(fml, co)
	case LIST:
		cinf.List(fml, co)
	case PLUS:
		mop(fml, cinf, "+", co)
	case MINUS:
		mop(fml, cinf, "-", co)
	case MULT:
		mop(fml, cinf, "*", co)
	case DIV:
		mop(fml, cinf, "/", co)
	case POW:
		mop(fml, cinf, "^", co)
	case NAME, NUMBER:
		co.append(fml.str)
	case F_TRUE:
		co.append(cinf.Ftrue())
	case F_FALSE:
		co.append(cinf.Ffalse())
	case UNARYMINUS:
		uniop(fml, cinf, "-", co)
	case UNARYPLUS:
		uniop(fml, cinf, "+", co)
	default:
		errors.New("unknown type")
	}
}

func prefix(fml Formula, cinf CnvInf, left, right string, co *cnv_out) {
	co.append(left)
	sep := ""
	for i := 0; i < len(fml.args); i++ {
		co.append(sep)
		conv2(fml.args[i], cinf, co)
		sep = ","
	}
	co.append(right)
}

func infix(fml Formula, cinf CnvInf, op string, co *cnv_out) {
	sep := ""
	for i := 0; i < len(fml.args); i++ {
		co.append(sep)
		if fml.priority > 0 && fml.priority < fml.args[i].priority {
			co.append("(")
			conv2(fml.args[i], cinf, co)
			co.append(")")
		} else {
			conv2(fml.args[i], cinf, co)
		}
		sep = op
	}
}
