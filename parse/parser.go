package parse

import (
    "fmt"
    "strconv"
)

type Value interface {}

type lispFloat struct {
    value float64
}

type lispInt struct {
    value int64
}

type lispString struct {
    value string
}

type Identifier struct {
    value string
}

type lispFunc struct {
    fn func(args... Value) Value
}

type List struct {
    values []Value
}

func newList(values ...Value) List {
    return List{values: values}
}

func (l *List) push(vals ...Value) {
    l.values = append(l.values, vals...)
}

func last(l List) Value {
    i := len(l.values)
    return l.values[i - 1]
}


func (l lispFunc) call(args... Value) Value {
    return l.fn(args...)
}

func Parse(s string) Value {
    lexer := Lex("test.lisp", "(+ 1 2 3)")
    stack := []List{}
    stack = append(stack, newList{})

    item := lexer.NextToken()
    for item.typ != itemEOF {
        switch item.typ {
        case itemStartList:

        case itemCloseList:

        case itemIdentifier:

        case itemNumber:
        }

        item = lexer.NextToken()
    }
    fmt.Println(stack[0])

    return stack

}


func parseNumber(s string) Value {
    i, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        f, err := strconv.ParseFloat(s, 32)
        if err != nil {
            return lispFloat{f}
        }
    }
    return lispInt{i}
}
func (l List) String() string {
    str := "("
    for i, item := range(l.values) {
            switch item.(type) {
            case int:
                if i == len(l.values) - 1 {
                    str += fmt.Sprintf("%d)", item.(int))
                } else {
                    str += fmt.Sprintf("%d ", item.(int))
                }
            case List:
                if i == len(l.values) - 1 {
                    str += item.(List).String() + ")"
                } else {
                    str += item.(List).String() + " "
                }
            }
    }
    return str
}
