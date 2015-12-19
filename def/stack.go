package qeconv

import (
	"errors"
)

type QeStack struct {
	v []QeNode
}

func (s *QeStack) pop() (QeNode, error) {
	if len(s.v) <= 0 {
		return QeNode{}, errors.New("empty stack")
	}
	v := s.v[len(s.v)-1]
	s.v = s.v[:len(s.v)-1]
	return v, nil
}

func (s *QeStack) Popn(n int) *QeStack {
	stack := new(QeStack)
	m := len(s.v) - n
	stack.v = append(stack.v, s.v[m:]...)
	s.v = s.v[:m]
	return stack
}

func (s *QeStack) Push(v QeNode) {
	s.v = append(s.v, v)
}

func (s *QeStack) Pushn(v *QeStack) {
	s.v = append(s.v, v.v...)
}

func (s *QeStack) empty() bool {
	return len(s.v) == 0
}

func (s *QeStack) Length() int {
	return len(s.v)
}

func (s *QeStack) String() string {
	ret := ""
	for i := 0; i < len(s.v); i++ {
		if i == 0 {
			ret = s.v[i].String()
		} else {
			ret = ret + " " + s.v[i].String()
		}
	}
	return ret
}
