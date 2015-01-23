package qeconv

import (
	"testing"
)

func TestToVarSet(t *testing.T) {
	var vs varSet

	vs.remove("0")
	if !vs.valid() || vs.Len() != 0 {
		t.Errorf("remove0")
	}
	vs.append("2")
	if !vs.valid() || vs.Len() != 1 {
		t.Errorf("append0")
	}
	vs.append("1")
	if !vs.valid() || vs.Len() != 2 {
		t.Errorf("append1")
	}
	vs.append("3")
	if !vs.valid() || vs.Len() != 3 {
		t.Errorf("append2")
	}
	vs.append("2")
	if !vs.valid() || vs.Len() != 3 {
		t.Errorf("append3")
	}
	vs.append("3")
	if !vs.valid() || vs.Len() != 3 {
		t.Errorf("append4")
	}
	vs.append("4")
	if !vs.valid() || vs.Len() != 4 {
		t.Errorf("append5")
	}
	vs.append("4")
	if !vs.valid() || vs.Len() != 4 {
		t.Errorf("append6")
	}
	vs.append("5")
	if !vs.valid() || vs.Len() != 5 {
		t.Errorf("append7")
	}
	vs.remove("5")
	if !vs.valid() || vs.Len() != 4 {
		t.Errorf("remove8")
	}
	vs.remove("1")
	if !vs.valid() || vs.Len() != 3 {
		t.Errorf("remove9")
	}
	vs.remove("9")
	if !vs.valid() || vs.Len() != 3 {
		t.Errorf("remove10")
	}
	vs.remove("1")
	if !vs.valid() || vs.Len() != 3 {
		t.Errorf("remove11")
	}
	vs.remove("3")
	if !vs.valid() || vs.Len() != 2 {
		t.Errorf("remove12")
	}
}
