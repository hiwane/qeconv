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
	)

	flag.StringVar(&from, "f", "syn", "from {syn}")
	flag.StringVar(&to, "t", "syn", "to {math|tex|qep|red|syn}")
	flag.StringVar(&filename, "i", "", "input file")
	flag.StringVar(&output, "o", "", "output file")
	flag.Parse()
	var err error
	var b []byte
	if filename == "" {
		b, err = ioutil.ReadAll(os.Stdin)
	} else {
		b, err = ioutil.ReadFile(filename)
	}

	if err == nil {
		var str string
		if to == "math" {
			str = qeconv.ToMath(string(b))
		} else if to == "tex" {
			str = qeconv.ToLaTeX(string(b))
		} else if to == "syn" {
			str = qeconv.ToSyn(string(b))
		} else if to == "red" {
			str = qeconv.ToRedlog(string(b))
		} else if to == "qep" {
			str, err = qeconv.ToQepcad(string(b))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		} else if to == "rc" {
			str, err = qeconv.ToRegularChains(string(b))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(os.Stderr, "unsupported -t "+to)
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
