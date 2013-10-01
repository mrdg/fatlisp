package fatlisp

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type parser struct {
	name  string // Name of the parsed file.
	input string // Source string
	lex   *lexer

	// Used to keep track of open lists while parsing tokens.
	stack       []*Value
	currentList *Value

	// Contains all quotes encountered during parsing. After
	// parsing quotes will be expanded. e.g. '(1 2 3) -> (quote (1 2 3))
	quotes []Quote
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
	boolType
)

type Value struct {
	typ  Type
	data interface{}
}

type List struct {
	values *[]Value
}

// Represents quotes found in the lexer stream.  Their location in the tree
// (list + index) is stored.  After the entire tree is parsed, the element in
// list on index is replaced with (quote element). 'id' is the identifier that
// becomes the first element in the list the quote expands to.
type Quote struct {
	list  *Value
	index int
	id    string
}

func newParser(name, input string) parser {
	root := newList()

	return parser{
		name:        name,
		input:       input,
		lex:         Lex(name, input),
		stack:       []*Value{&root},
		currentList: &root,
	}
}

func Parse(name, input string) Value {
	p := newParser(name, input)
	return p.parse()
}

func (p parser) parse() Value {
	item := p.lex.NextToken()
	for item.typ != itemEOF {
		switch item.typ {
		case itemStartList:
			list := newList()
			p.currentList.push(list)
			p.pushList(&list)

		case itemCloseList:
			p.popList()

		case itemIdentifier:
			p.currentList.push(parseIdentifier(item.val))

		case itemNumber:
			p.currentList.push(parseNumber(item.val))

		case itemString:
			p.currentList.push(parseString(item.val))

		case itemError:
			fmt.Printf("Error: %s - %s\n", p.errorPos(item), item.val)
			os.Exit(-1)

		case itemQuote:
			i := len(vtos(*p.currentList))
			q := Quote{list: p.currentList, index: i, id: "quote"}
			p.quotes = append(p.quotes, q)
		}

		item = p.lex.NextToken()
	}
	p.expandQuotes()
	return *p.currentList
}

func (p *parser) expandQuotes() {
	for _, q := range p.quotes {
		list := newList()
		list.push(Value{typ: idType, data: q.id})
		list.push(q.list.get(q.index))
		q.list.replace(q.index, list)
	}
}

func (p *parser) pushList(list *Value) {
	p.stack = append(p.stack, list)
	p.currentList = list
}

func (p *parser) popList() {
	p.stack = p.stack[:len(p.stack)-1]
	p.currentList = p.stack[len(p.stack)-1]
}

func (p parser) errorPos(i item) string {
	str := p.input[:i.pos-1]
	lastLine := strings.LastIndex(str, "\n")

	var col int
	if lastLine == -1 {
		col = i.pos
	} else {
		col = i.pos - (lastLine + 1)
	}

	lineNum := 1 + strings.Count(str, "\n")
	return fmt.Sprintf("%s:%d:%d", p.name, lineNum, col)
}

func (list *Value) push(val Value) {
	l := list.data.(List)
	*l.values = append(*l.values, val)
}

func (list *Value) replace(index int, val Value) {
	l := list.data.(List)
	(*l.values)[index] = val
}

func (list Value) get(index int) Value {
	l := list.data.(List)
	return (*l.values)[index]
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

func vtob(v Value) bool {
	return v.data.(bool)
}

func parseString(s string) Value {
	// Strip of quotes that are included in the token.
	s = s[1:]
	s = s[:len(s)-1]

	return Value{typ: stringType, data: s}
}

func parseIdentifier(s string) Value {
	if s == "true" {
		return Value{typ: boolType, data: true}
	}
	if s == "false" {
		return Value{typ: boolType, data: false}
	}
	if s == "nil" {
		return Value{typ: nilType, data: nil}
	}
	return Value{typ: idType, data: s}
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
	case stringType:
		return v.data.(string)
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
	case boolType:
		return fmt.Sprintf("%v", vtob(v))
	case fnType:
		return fmt.Sprintf("<fn>")
	default:
		return v.String()
	}
}

func (t Type) String() string {
	var s string
	switch t {
	case intType:
		s = "Int"
	case floatType:
		s = "Float"
	case stringType:
		s = "String"
	case idType:
		s = "Identifier"
	case listType:
		s = "List"
	case fnType:
		s = "Fn"
	case nilType:
		s = "Nil"
	case boolType:
		s = "Bool"
	}
	return s
}
