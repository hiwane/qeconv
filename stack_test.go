package qeconv

import (
	"testing"
)


func TestStack(t *testing.T) {
	var s *Stack
	s = new(Stack)
	if !s.empty() {
		t.Error("empty0")
	}
	s.push(Node{cmd:1})
	if s.empty() {
		t.Error("empty1")
	}
	s.push(Node{cmd:2})
	if s.empty() {
		t.Error("empty2")
	}

	u, err := s.pop()
	if err != nil || u.cmd != 2 {
		t.Error("pop2")
	}
	if s.empty() {
		t.Error("empty2-2")
	}

	s.push(Node{cmd:3})
	if s.empty() {
		t.Error("empty3")
	}

	s.push(Node{cmd:4})
	if s.empty() {
		t.Error("empty4")
	}
	u, err = s.pop()
	if err != nil || u.cmd != 4 {
		t.Error("pop4")
	}
	if s.empty() {
		t.Error("empty4-2")
	}

	u, err = s.pop()
	if err != nil || u.cmd != 3 {
		t.Error("pop3")
	}
	if s.empty() {
		t.Error("empty3-2")
	}

	u, err = s.pop()
	if err != nil || u.cmd != 1 {
		t.Error("pop1")
	}
	if !s.empty() {
		t.Error("empty1-2")
	}

	u, err = s.pop()
	if err == nil {
		t.Error("pop0")
	}
}
