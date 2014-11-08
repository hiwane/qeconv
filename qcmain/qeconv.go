package main

import (
	"flag"
	"fmt"
	"os"
	"io/ioutil"

	"github.com/hiwane/qeconv"
)

func main() {

	var (
		filename string
		output string
		from string
		to string)

	flag.StringVar(&from,     "f" , "syn", "from {syn}")
	flag.StringVar(&to,       "t" , "math", "to {math}")
	flag.StringVar(&filename, "i" , "", "input file")
	flag.StringVar(&output,   "o" , "", "output file")
	flag.Parse()
	var err error
	var b []byte
	if filename == "" {
		b, err = ioutil.ReadAll(os.Stdin)
	} else {
		b, err = ioutil.ReadFile(filename)
	}

	if err == nil {
		str := qeconv.ToMath(string(b))
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
