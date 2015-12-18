package qeconv

import (
	"errors"
	. "github.com/hiwane/qeconv/def"
	tex "github.com/hiwane/qeconv/latex"
	math "github.com/hiwane/qeconv/math"
	qep "github.com/hiwane/qeconv/qepcad"
	red "github.com/hiwane/qeconv/redlog"
	rc "github.com/hiwane/qeconv/regularchains"
	smt2 "github.com/hiwane/qeconv/smt2"
	syn "github.com/hiwane/qeconv/synrac"
	"strconv"
)

func cmpfml(f1, f2 Formula) bool {
	if f1.Cmd() != f2.Cmd() {
		return false
	}

	if len(f1.Args()) != len(f2.Args()) {
		return false
	}

	for i := 0; i < len(f1.Args()); i++ {
		if !cmpfml(f1.Arg(i), f2.Arg(i)) {
			return false
		}
	}

	return true
}

// 重複定義の削除
func rmdup(fml Formula) Formula {
	if !fml.IsList() {
		return fml
	}

	for i := 1; i < len(fml.Args()); i++ {
		for j := 0; j < i; j++ {
			if cmpfml(fml.Arg(i), fml.Arg(j)) {
				args := fml.Args()
				fml.SetArgLen(len(args) - 1)
				for k := 0; k < i; k++ {
					fml.SetArg(k, args[k])
				}
				for k := i + 1; k < len(args); k++ {
					fml.SetArg(k-1, args[k])
				}
				break
			}
		}
	}

	return fml
}

func cntfml(fml Formula) int {
	if !fml.IsList() {
		return 1
	}

	count := 0
	for _, v := range fml.Args() {
		count += cntfml(v)
	}

	return count
}

func getfmlidx(fml Formula, idx int) Formula {

	if !fml.IsList() {
		return fml
	}

	count := 0
	for _, v := range fml.Args() {
		c := cntfml(v)
		if c+count > idx {
			return getfmlidx(v, idx-count)
		}
		count += c
	}

	return fml
}

func Str2Parser(to string) (Parser, error) {
	if to == "syn" {
		return syn.NewSynParse(), nil
	} else if to == "smt2" {
		return smt2.NewSmt2Parse(), nil
	}
	return nil, errors.New("unknown converter")
}

func Str2cinf(to string) (CnvInf, error) {
	if to == "math" {
		return new(math.MathConv), nil
	} else if to == "tex" {
		return new(tex.LatexConv), nil
	} else if to == "syn" {
		return new(syn.SynConv), nil
	} else if to == "red" {
		return new(red.RedConv), nil
	} else if to == "qep" {
		return new(qep.QepConv), nil
	} else if to == "smt2" {
		return new(smt2.Smt2Conv), nil
	} else if to == "rc" {
		return new(rc.RegChainConv), nil
	}
	return nil, errors.New("unknown converter")
}

func Convert(p Parser, cinf CnvInf, str string, dup bool, index int) (string, error) {
	var ret string
	count := 0

	for {
		// コメント行を無視して, separator である : を探索する.
		idx := p.Next(str)
		if idx <= 0 {
			break
		}

		fml, cmts, err := p.Parse(str[:idx])
		if err != nil {
			return "", err
		}
		str = str[idx:]
		if dup {
			fml = rmdup(fml)
		}

		if index != 0 {
			cnt := cntfml(fml)
			if index < 0 {
				count += cnt
				ret = strconv.Itoa(count)
				continue
			} else if count > index-1 || index-1 >= count+cnt {
				count += cnt
				continue
			} else {
				fml = getfmlidx(fml, index-count-1)
				count += cnt
				cmts = make([]Comment, 0)
			}
		}

		var str2 string
		var sep string
		sep = cinf.Sep()
		co := NewCnvOut(cmts)
		str2, err = cinf.Convert(fml, co)
		if err != nil {
			return str2, err
		}
		ret += str2 + sep
	}

	return ret, nil
}
