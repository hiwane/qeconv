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
	)

	flag.StringVar(&from, "f", "syn", "from {syn}")
	flag.StringVar(&to, "t", "syn", "to {tex|math|qep|red|rc|syn|smt2}")
	flag.StringVar(&filename, "i", "", "input file")
	flag.StringVar(&output, "o", "", "output file")
	flag.BoolVar(&dup, "s", false, "dup")
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
		str,err := qeconv.Convert(string(b), to, dup)
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
