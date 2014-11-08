package qeconv

import (
	"errors"
)

type Stack struct {
	v []Node
}

func (s *Stack) pop() Node {
	if len(s.v) <= 0 {
		errors.New("empty stack")
	}
	v := s.v[len(s.v)-1]
	s.v = s.v[:len(s.v)-1]
	return v
}

func (s *Stack) push(v Node) {
	s.v = append(s.v, v)
}

func (s *Stack) empty() bool {
	return len(s.v) == 0
}
