package main

import(
    "fmt"
)

type lispFloat struct {
    value float64
}

type lispInt struct {
    value int
}

type lispValue interface {}

func (l lispInt) String() string {
    return fmt.Sprint(l.value)
}

type Node struct {
    value lispValue
    children NodeList
}

type NodeList []Node

type Env struct {
    parent *Env
    defs map[string]lispValue
}

type lispFunc struct {
    fn func(args... lispValue) lispValue
}

func (l lispFunc) call(args... lispValue) lispValue {
    return l.fn(args...)
}

func main() {
    fmt.Println("Go!")

    a := lispInt{2}
    b := lispInt{3}
    c := lispInt{5}
    d := lispInt{6}

    lispAdd := lispFunc{fn: add}

    args2 := NodeList{Node{value: c}, Node{value: d}}
    branch := Node{value: lispAdd, children: args2}

    args := NodeList{Node{value: a}, Node{value: b}, branch}
    root := Node{value: lispAdd, children: args}

    value := eval(root)
    fmt.Println(value)
}

func eval(node Node) lispValue {
    // TODO: handle special forms

    if len(node.children) > 0 {
        args := make([]lispValue, len(node.children))
        for i, c := range node.children {
            args[i] = eval(c)
        }
        fn, _ := node.value.(lispFunc)
        return fn.call(args...)
    } else {
        return node.value
    }
}

func add(args... lispValue) lispValue {
    returnFloat := false
    fsum := 0.0
    isum := 0

    for _, n := range args {
        switch t := n.(type) {
        case lispFloat:
            returnFloat = true
            fsum += n.(lispFloat).value
        case lispInt:
            isum += n.(lispInt).value
        default:
            panic(fmt.Sprintf("+: Unexpected type %T", t))
        }
    }

    if returnFloat {
        fsum += float64(isum)
        return lispFloat{fsum}
    } else {
        return lispInt{isum}
    }
}
