package qeconv

import (
	"bytes"
	"errors"
)

const NAME = 57346
const NUMBER = 57347
const F_TRUE = 57348
const F_FALSE = 57349
const ALL = 57350
const EX = 57351
const AND = 57352
const OR = 57353
const NOT = 57354
const ABS = 57355
const PLUS = 57356
const MINUS = 57357
const COMMA = 57358
const MULT = 57359
const DIV = 57360
const POW = 57361
const EOL = 57362
const LB = 57363
const RB = 57364
const LP = 57365
const RP = 57366
const LC = 57367
const RC = 57368
const INDEXED = 57369
const LIST = 57370
const IMPL = 57371
const REPL = 57372
const EQUIV = 57373
const COMMENT = 57374
const LTOP = 57375
const GTOP = 57376
const LEOP = 57377
const GEOP = 57378
const NEOP = 57379
const EQOP = 57380
const UNARYMINUS = 57381
const UNARYPLUS = 57382

type Comment struct {
	lineno int
	str    string
}

type CnvOut struct {
	strbuf  *bytes.Buffer
	lno     int
	comment []Comment

	// input info
	//	input string
	//	index int
}

type CnvInfMathOp interface {
	/* math op */
	Plus(f Formula, co *CnvOut)
	Minus(f Formula, co *CnvOut)
	Mult(f Formula, co *CnvOut)
	Div(f Formula, co *CnvOut)
	Pow(f Formula, co *CnvOut)
	Comment(str string) string
	Abs(f Formula, co *CnvOut)
	Uniop(f Formula, ope string, co *CnvOut)

	/* atom */
	Ftrue() string
	Ffalse() string
	LineAlign() bool
}

type CnvInf interface {
	CnvInfMathOp
	Convert(fml Formula, co *CnvOut) (string, error)
	Sep() string

	/* quantifier */
	All(f Formula, co *CnvOut)
	Ex(f Formula, co *CnvOut)

	/* logical operator
	 * Mathematica: Implies/Equivalent
	 * redlog     : impl/repl/equiv
	 * qepcad     : =>/<=/<=>
	 */
	And(f Formula, co *CnvOut)
	Or(f Formula, co *CnvOut)
	Not(f Formula, co *CnvOut)
	Impl(f Formula, co *CnvOut)
	Equiv(f Formula, co *CnvOut)

	/* comparator */
	Leop(f Formula, co *CnvOut)
	Ltop(f Formula, co *CnvOut)
	Eqop(f Formula, co *CnvOut)
	Neop(f Formula, co *CnvOut)

	List(f Formula, co *CnvOut)
}

type CnvInfStrstruct struct {
}

func (self *CnvInfStrstruct) LineAlign() bool {
	return true
}

type Parser interface {
	Parse(str string, varcnv bool) (Formula, []Comment, error)

	// str $B$r2r@O$7!"<!$KFI$`$H$3$m$^$G$N(B index $B$rJV$9(B
	Next(str string) int
}

func NewCnvOut(comment []Comment) *CnvOut {
	co := new(CnvOut)
	co.Reset()
	co.comment = comment
	return co
}

func (c *CnvOut) Append(s string) {
	//	fmt.Printf("append [%s]\n", s)
	c.strbuf.WriteString(s)
}

func (c *CnvOut) String() string {
	//	fmt.Printf("append [%s]\n", s)
	return c.strbuf.String()
}

func (c *CnvOut) Reset() {
	c.strbuf = new(bytes.Buffer)
	c.lno = 1
}

func NewComment(str string, lno int) Comment {
	return Comment{str: str, lineno: lno}
}

func skipcomment(fml Formula, cinf CnvInfMathOp, co *CnvOut) {
	if cinf.LineAlign() {
		for co.lno < fml.lineno {
			if len(co.comment) > 0 && co.comment[0].lineno == co.lno {
				co.Append(cinf.Comment(co.comment[0].str))
				co.comment = co.comment[1:len(co.comment)]
			}
			co.Append("\n")
			co.lno++
		}
	}
}

func Convm(fml Formula, cinf CnvInfMathOp, co *CnvOut) {

	skipcomment(fml, cinf, co)

	switch fml.cmd {
	case PLUS:
		cinf.Plus(fml, co)
	case MINUS:
		cinf.Minus(fml, co)
	case MULT:
		cinf.Mult(fml, co)
	case DIV:
		cinf.Div(fml, co)
	case POW:
		cinf.Pow(fml, co)
	case ABS:
		cinf.Abs(fml, co)
	case NAME, NUMBER:
		co.Append(fml.str)
	case F_TRUE:
		co.Append(cinf.Ftrue())
	case F_FALSE:
		co.Append(cinf.Ffalse())
	case UNARYMINUS:
		cinf.Uniop(fml, "-", co)
	case UNARYPLUS:
		cinf.Uniop(fml, "+", co)
	default:
		errors.New("unknown type")
	}
}

// $B$[$\(B Infix(). call Convm() instead of Conv2()
func Mop(fml Formula, cinf CnvInfMathOp, op string, co *CnvOut) {
	for i := 0; i < len(fml.Args()); i++ {
		if i != 0 {
			co.Append(op)
		}
		if fml.priority > 0 && fml.priority < fml.Args()[i].priority {
			co.Append("(")
			Convm(fml.Args()[i], cinf, co)
			co.Append(")")
		} else {
			Convm(fml.Args()[i], cinf, co)
		}
	}
}

func Uniop(fml Formula, cinf CnvInfMathOp, op string, co *CnvOut) {
	co.Append("(")
	co.Append(op)
	if fml.priority > 0 && fml.priority < fml.Args()[0].priority {
		co.Append("(")
		Convm(fml.Args()[0], cinf, co)
		co.Append(")")
	} else {
		Convm(fml.Args()[0], cinf, co)
	}
	co.Append(")")
}

func Conv2(fml Formula, cinf CnvInf, co *CnvOut) {
	//	fmt.Printf("fml.cmd=%d,lineno=%d/%d str=%s\n", fml.cmd, fml.lineno, co.lno, fml.str)
	skipcomment(fml, cinf, co)

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
	default:
		Convm(fml, cinf, co)
	}
}

func Prefixm(fml Formula, cinf CnvInf, left, mid, right string, co *CnvOut) {
	co.Append(left)
	sep := ""
	for i := 0; i < len(fml.Args()); i++ {
		co.Append(sep)
		Conv2(fml.Args()[i], cinf, co)
		sep = mid
	}
	co.Append(right)
}

func Prefix(fml Formula, cinf CnvInf, left, right string, co *CnvOut) {
	Prefixm(fml, cinf, left, ",", right, co)
}

func Infixm(fml Formula, cinf CnvInf, op string, co *CnvOut, str, end string) {
	sep := ""
	for i := 0; i < len(fml.Args()); i++ {
		co.Append(sep)
		if fml.priority > 0 && fml.priority < fml.Args()[i].priority {
			co.Append(str)
			Conv2(fml.Args()[i], cinf, co)
			co.Append(end)
		} else {
			Conv2(fml.Args()[i], cinf, co)
		}
		sep = op
	}
}

func Infix(fml Formula, cinf CnvInf, op string, co *CnvOut) {
	Infixm(fml, cinf, op, co, "(", ")")
}
