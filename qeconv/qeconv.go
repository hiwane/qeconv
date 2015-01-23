package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hiwane/qeconv"
)

func str2cinf(to string) (ConvInf, error) {
	if to == "math" {
		return new(mathConv)
	} else if to == "tex" {
		return new(latexConv)
	} else if to == "syn" {
		return new(synConv)
	} else if to == "red" {
		return new(redConv)
	} else if to == "qep" {
		return new(qepConv)
	} else if to == "smt2" {
		return new(smt2Conv)
	} else if to == "rc" {
		return new(smt2Conv)
	} else {
		return nil, errors.New("unknown converter")
	}
}

func main() {

	var (
		filename string
		output   string
		from     string
		to       string
		dup      bool
		idx      int
	)

	flag.StringVar(&from, "f", "syn", "from {syn}")
	flag.StringVar(&to, "t", "syn", "to {tex|math|qep|red|rc|syn|smt2}")
	flag.StringVar(&filename, "i", "", "input file")
	flag.StringVar(&output, "o", "", "output file")
	flag.BoolVar(&dup, "s", false, "remove duplicate formulas")
	flag.IntVar(&idx, "n", 0, "index of fof")
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
		cinf, err := str2cinf(to)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		str, err := qeconv.Convert(cinf, string(b), dup, idx)
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
