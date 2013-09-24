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
        v := parse.Parse(string(src))
        parse.Eval(v)
    }
}
