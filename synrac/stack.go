package qeconv

import (
	"errors"
)

type synStack struct {
	v []synNode
}

func (s *synStack) pop() (synNode, error) {
	if len(s.v) <= 0 {
		return synNode{}, errors.New("empty stack")
	}
	v := s.v[len(s.v)-1]
	s.v = s.v[:len(s.v)-1]
	return v, nil
}

func (s *synStack) push(v synNode) {
	s.v = append(s.v, v)
}

func (s *synStack) empty() bool {
	return len(s.v) == 0
}
