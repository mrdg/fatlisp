package main

import (
	"fmt"
	"io/ioutil"
	"github.com/mrdg/lisp/parse"
)

func main() {
	src, err := ioutil.ReadFile("test.clj")
	if err != nil {
		fmt.Println(err)
	} else {
		tree := parse.NewTree("test.clj", string(src))
		nodes := tree.Parse(string(src))
		parse.Eval(nodes)
	}
}
