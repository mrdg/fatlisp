package main

import (
	"fmt"
	"github.com/mrdg/lisp/parse"
)

func main() {
	fmt.Println(parse.Parse("+ 1 1"))
}
