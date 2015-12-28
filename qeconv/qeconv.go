package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hiwane/qeconv"
)

func main() {

	var (
		filename string
		output   string
		from     string
		to       string
		dup      bool
		cnv      bool
		idx      int
	)

	flag.StringVar(&from, "f", "syn", "from {syn|smt2}")
	flag.StringVar(&to, "t", "syn", "to {tex|math|qep|red|rc|syn|smt2}")
	flag.StringVar(&filename, "i", "", "input file")
	flag.StringVar(&output, "o", "", "output file")
	flag.BoolVar(&dup, "s", false, "remove duplicate formulas")
	flag.IntVar(&idx, "n", 0, "index of fof")
	flag.BoolVar(&cnv, "X", false, "convert symbol name")
	flag.Parse()
	if flag.NArg() > 0 {
		flag.PrintDefaults()
		os.Exit(2)
	}
	var err error
	var b []byte
	if filename == "" {
		b, err = ioutil.ReadAll(os.Stdin)
	} else {
		b, err = ioutil.ReadFile(filename)
	}

	if err == nil {
		cinf, err := qeconv.Str2cinf(to)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		parser, err := qeconv.Str2Parser(from)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		str, err := qeconv.Convert(parser, cinf, string(b), dup, idx, cnv)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if output == "" {
			fmt.Printf("%s\n", str)
		} else {
			ioutil.WriteFile(output, []byte(str), 0644)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
