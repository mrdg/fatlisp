package main

import (
	"flag"
	"fmt"
	"github.com/mrdg/fatlisp"
	"io/ioutil"
	"os"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "usage: %s <file>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	file := flag.Arg(0)

	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	tree, err := fatlisp.Parse(file, string(src))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx := fatlisp.NewContext()

	_, err = ctx.Eval(tree)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
