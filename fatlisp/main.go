package main

import (
	"fmt"
	"github.com/mrdg/fatlisp"
	"io/ioutil"
)

func main() {
	src, err := ioutil.ReadFile("test.clj")
	if err != nil {
		fmt.Println(err)
	} else {
		tree := fatlisp.NewTree("test.clj", string(src))
		nodes := tree.Parse()
		fatlisp.Eval(nodes)
	}
}
