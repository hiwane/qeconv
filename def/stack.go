package qeconv

import (
	"errors"
)

type QeStack struct {
	v [] QeNode
}

func (s *QeStack) pop() (QeNode, error) {
	if len(s.v) <= 0 {
		return QeNode{}, errors.New("empty stack")
	}
	v := s.v[len(s.v)-1]
	s.v = s.v[:len(s.v)-1]
	return v, nil
}

func (s *QeStack) Push(v QeNode) {
	s.v = append(s.v, v)
}

func (s *QeStack) empty() bool {
	return len(s.v) == 0
}
