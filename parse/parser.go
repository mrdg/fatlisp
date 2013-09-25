package parse

import (
	"fmt"
	"strconv"
	"os"
	"strings"
)

type Tree struct {
	name string
	input string
	lex *lexer
}

type Type int

const (
	intType Type = iota
	floatType
	stringType
	listType
	idType
	fnType
	nilType
)

type Value struct {
	typ  Type
	data interface{}
}

type List struct {
	values *[]Value
}

func NewTree(name, input string) Tree {
	return Tree{
		name: name,
		input: input,
		lex: Lex(name, input),
	}
}

func (tree Tree) Parse(s string) Value {
	root := newList()
	stack := []*Value{&root}

	item := tree.lex.NextToken()
	for item.typ != itemEOF {
		switch item.typ {
		case itemStartList:
			current := stack[len(stack)-1]
			list := newList()
			current.push(list)
			stack = append(stack, &list)

		case itemCloseList:
			stack = stack[:len(stack)-1]

		case itemIdentifier:
			current := stack[len(stack)-1]
			current.push(Value{typ: idType, data: item.val})

		case itemNumber:
			current := stack[len(stack)-1]
			current.push(parseNumber(item.val))

		case itemError:
			fmt.Printf("Error: %s - %s\n", tree.errorPos(item), item.val)
			os.Exit(-1)
		}

		item = tree.lex.NextToken()
	}
	return *(stack[0])
}

func (list *Value) push(val Value) {
	l := (*list).data.(List)
	*l.values = append(*l.values, val)
}

func newList(vals ...Value) Value {
	return Value{typ: listType, data: List{values: &vals}}
}

type Fn func(args ...Value) Value

func newFn(fn Fn) Value {
	return Value{typ: fnType, data: fn}
}

func newInt(i int64) Value {
	return Value{typ: intType, data: i}
}

func newFloat(f float64) Value {
	return Value{typ: floatType, data: f}
}

func vtos(v Value) []Value {
	list := v.data.(List)
	return *list.values
}

func vtoi(v Value) int64 {
	return v.data.(int64)
}

func vtof(v Value) float64 {
	return v.data.(float64)
}

func (tree Tree) errorPos(i item) string {
	pos := i.pos + 1 // i.pos is zero indexed, so add 1
	str := tree.input[:pos]
	lines := 1 + strings.Count(str, "\n")
	lastLine := strings.LastIndex(str, "\n")
	var col int
	if lastLine == -1 {
		col = pos
	} else {
		col = pos - (lastLine + 1)
	}

	return fmt.Sprintf("%s:%d:%d", tree.name, lines, col)
}

func parseNumber(s string) Value {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return newFloat(f)
		} else {
			panic("Invalid number syntax")
		}
	}
	return newInt(i)
}

func (v Value) String() string {
	switch v.typ {
	case intType:
		return fmt.Sprintf("%d", vtoi(v))
	case floatType:
		return fmt.Sprintf("%v", vtof(v))
	case idType:
		return v.data.(string)
	case listType:
		str := "("
		list := v.data.(List)
		for i, val := range *list.values {
			str += val.String()
			if i != len(*list.values)-1 {
				str += " "
			}
		}
		str += ")"
		return str
	case nilType:
		return "nil"
	default:
		return v.String()
	}
}
