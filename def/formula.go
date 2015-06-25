package qeconv

type Formula struct {
	cmd      int
	str      string
	args     []Formula
	priority int
	lineno   int
}

func NewFormula(cmd int, str string, lno int, priority int) Formula {
	return Formula{cmd: cmd, str: str, lineno: lno, priority: priority}
}

func (self *Formula) Cmd() int {
	return self.cmd
}
func (self *Formula) String() string {
	return self.str
}
func (self *Formula) Priority() int {
	return self.priority
}
func (self *Formula) IsList() bool {
	return self.cmd == LIST
}
func (self *Formula) IsQuantifier() bool {
	return self.cmd == ALL || self.cmd == EX
}
func (self *Formula) IsAtom() bool {
	return self.cmd == LEOP || self.cmd == LTOP ||
		self.cmd == EQOP || self.cmd == NEOP
}
func (self *Formula) IsBool() bool {
	return self.cmd == F_TRUE || self.cmd == F_FALSE
}

func (self *Formula) Args() []Formula {
	return self.args
}
func (self *Formula) Arg(id int) Formula {
	return self.args[id]
}

func (self *Formula) SetArgLen(n int) {
	self.args = make([]Formula, n)
}

func (self *Formula) SetArg(k int, f Formula) {
	self.args[k] = f
}

func (self *Formula) IsQff() bool {
	if self.IsQuantifier() {
		return false
	} else if self.IsAtom() || self.IsBool() {
		return true
	}
	for _, v := range self.Args() {
		if !v.IsQff() {
			return false
		}
	}
	return true
}

func (self *Formula) freeVarsAtom() varSet {
	var vs varSet

	if self.cmd == NAME {
		vs.append(self.str)
		return vs
	} else if self.cmd == NUMBER {
		return vs
	}

	for _, v := range self.Args() {
		vs.union(v.freeVarsAtom())
	}

	return vs
}

func (self *Formula) FreeVars() varSet {

	var vs varSet
	if self.IsBool() {
		return vs
	} else if self.IsAtom() {
		return self.freeVarsAtom()
	} else if self.IsQuantifier() {
		fv := self.Args()[1].FreeVars()
		qv := self.Args()[0].FreeVars()
		// arg[0] を削除.
		fv.setminus(qv)
		return fv
	} else if self.cmd == NAME {
		vs.append(self.str)
	} else {
		for _, v := range self.Args() {
			vs.union(v.FreeVars())
		}
	}
	return vs
}
