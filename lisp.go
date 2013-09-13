package main

import(
    "fmt"
    "regexp"
    "strings"
    "strconv"
)

var leftParen  = regexp.MustCompile("\\(")
var rightParen = regexp.MustCompile("\\)")
var whitespace = regexp.MustCompile("\\s+")

func main() {
	program := "(+ 1 1 (+ 2 2))"
	tokens := tokenize(program)
	tree := readTokens(&tokens)
    fmt.Println(tree)
}

func tokenize(str string) []string {
	str = leftParen.ReplaceAllString(str, " ( ")
	str = rightParen.ReplaceAllString(str, " ) ")
    return strings.Fields(str)
}

func readTokens(tokens *[]string) interface{} {
    t := pop(tokens)

    if t == "(" {
        list := make([]interface{}, 0)
        for (*tokens)[0] != ")" {
            list = append(list, readTokens(tokens))
        }

        pop(tokens) // pop off ")"
        return list

    } else {
        floatVal, err := strconv.ParseFloat(t, 64)
        if err == nil {
            return floatVal
        }

        intVal, err := strconv.Atoi(t)
        if err == nil {
            return intVal
        }

        return t
    }
}

func pop(s *[]string) string {
    ret := (*s)[0]
    *s = (*s)[1:]
    return ret
}
