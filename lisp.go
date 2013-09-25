package main

import (
	"fmt"
	"github.com/mrdg/lisp/parse"
	"io/ioutil"
)

func main() {
	src, err := ioutil.ReadFile("test.clj")
	if err != nil {
		fmt.Println(err)
	} else {
		v := parse.Parse(string(src))
		parse.Eval(v)
	}
}
