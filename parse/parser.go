package parse

import (
    "fmt"
    "strconv"
)

type Type int

const (
    intType Type = iota
    floatType
    stringType
    listType
    idType
)

type Value struct {
    typ Type
    data interface{}
}

type List struct {
    values *[]Value
}

func (list *Value) push(val Value) {
    l := (*list).data.(List)
    *l.values = append(*l.values, val)
}

func newList(vals ...Value) Value {
    slice := make([]Value, len(vals))
    for i, v := range(vals) {
        slice[i] = v
    }
    return Value{typ: listType, data: List{values: &slice}}
}

func newInt(i int64) Value {
    return Value{typ: intType, data: i}
}

func newFloat(f float64) Value {
    return Value{typ: floatType, data: f}
}


func Parse(s string) Value {
    lexer := Lex("test.lisp", "(+ 1 2 3 (+ 2 2 (* 3 4)))")
    stack := []*Value{}
    root := newList()
    stack = append(stack, &root)

    item := lexer.NextToken()
    for item.typ != itemEOF {
        switch item.typ {
        case itemStartList:
            current := stack[len(stack) - 1]
            list := newList()
            current.push(list)
            stack = append(stack, &list)

        case itemCloseList:
            stack = stack[:len(stack) - 1]

        case itemIdentifier:
            current := stack[len(stack) - 1]
            current.push(Value{typ: idType, data: item.val})

        case itemNumber:
            current := stack[len(stack) - 1]
            current.push(parseNumber(item.val))
        }

        item = lexer.NextToken()
    }
    return *(stack[0])
}

func parseNumber(s string) Value {
    i, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        f, err := strconv.ParseFloat(s, 32)
        if err != nil {
            return newFloat(f)
        }
    }
    return newInt(i)
}

func (v Value) String() string {
    switch v.typ {
    case intType:
        return fmt.Sprintf("%d", v.data.(int64))
    case idType:
        return v.data.(string)
    case listType:
        str := "("
        list := v.data.(List)
        for i, val := range(*list.values) {
            str += val.String()
            if i != len(*list.values) - 1 {
                str += " "
            }
        }
        str += ")"
        return str
    default:
        return v.String()
    }
}
