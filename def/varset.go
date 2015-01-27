package qeconv

import (
	"sort"
)

type varSet struct {
	v []string
}

func (self *varSet) Len() int {
	return len(self.v)
}
func (self *varSet) Get(i int) string {
	return self.v[i]
}

func (self *varSet) valid() bool {
	for i := 1; i < len(self.v); i++ {
		if self.v[i-1] >= self.v[i] {
			return false
		}
	}
	return true
}

func (self *varSet) append(str string) {
	idx := sort.SearchStrings(self.v, str)
	if idx >= self.Len() {
		self.v = append(self.v, str)
	} else if self.v[idx] != str {
		var v []string
		v = append(v, self.v[:idx]...)
		v = append(v, str)
		v = append(v, self.v[idx:]...)
		self.v = v
	}
}

func (self *varSet) remove(str string) {
	idx := sort.SearchStrings(self.v, str)
	if idx < self.Len() && self.v[idx] == str {
		if idx+1 == self.Len() {
			self.v = self.v[:idx]
		} else if idx == 0 {
			self.v = self.v[idx+1:]
		} else {
			self.v = append(self.v[:idx], self.v[idx+1:]...)
		}
	}
}

func (self *varSet) setminus(vs varSet) {
	for _, v := range vs.v {
		self.remove(v)
	}
}

func (self *varSet) union(vs varSet) {
	for _, v := range vs.v {
		self.append(v)
	}
}
