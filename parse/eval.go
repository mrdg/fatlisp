package parse

import "fmt"

type Env struct {
	parent *Env
	defs   map[string]Value
}

func newEnv() *Env {
	return &Env{defs: make(map[string]Value, 1)}
}

// Construct a new function scope based on parent scope.
func newFunctionEnv(parent *Env, params Value, args []Value) *Env {
	env := newEnv()
	env.parent = parent

	for i, p := range *(params.data.(List).values) {
		name := p.data.(string)
		value := args[i]
		env.set(name, value)
	}
	return env
}

func (e *Env) set(key string, v Value) {
	e.defs[key] = v
}

func (e Env) get(id string) Value {
	v, ok := e.defs[id]
	if !ok {
		if e.parent != nil {
			parent := *e.parent
			return parent.get(id)
		}
	}
	return v
}

func Eval(root Value) []Value {
	global := newEnv()
	global.set("+", newFn(add))
	global.set("puts", newFn(puts))

	results := []Value{}
	for _, v := range vtos(root) {
		results = append(results, eval(v, global))
	}
	return results
}

func eval(v Value, e *Env) Value {
	switch v.typ {
	case idType:
		id := v.data.(string)
		return e.get(id)
	case listType:
		list := vtos(v)
		args := make([]Value, len(list))

		if list[0].typ == idType {
			form, ok := specialForm(list[0])
			if ok {
				return form(e, list...)
			}
		}

		for i, c := range list {
			args[i] = eval(c, e)
		}
		fn := args[0].data.(Fn)
		args = args[1:] // Pop of the function

		return fn(args...)
	default:
		return v
	}
}

func specialForm(v Value) (func(env *Env, args ...Value) Value, bool) {
	name := v.data.(string)
	switch name {
	case "def":
		return def, true
	case "fn":
		return fn, true
	case "if":
		return ifForm, true
	case "quote":
		return quote, true
	default:
		return nil, false
	}
}

func quote(e *Env, vals ...Value) Value {
	return vals[1]
}

func fn(e *Env, vals ...Value) Value {
	vals = vals[1:] // Pop off fn keyword
	expectArgCount("fn", vals, 2)

	params := vals[0]
	body := vals[1]

	expectArg("fn", vals, 0, listType)

	return newFn(func(args ...Value) Value {
		argc := len(vtos(params))
		expectArgCount("function", args, argc)
		return eval(body, newFunctionEnv(e, params, args))
	})
}

func def(e *Env, args ...Value) Value {
	name := args[1].data.(string)
	val := eval(args[2], e)
	e.set(name, val)
	return args[1]
}

func ifForm(e *Env, args ...Value) Value {
	args = args[1:]
	expectArgCount("if", args, 3)

	cond := args[0]
	ifClause := args[1]
	elseClause := args[2]

	if truthy(eval(cond, e)) {
		return eval(ifClause, e)
	} else {
		return eval(elseClause, e)
	}
}

func truthy(v Value) bool {
	if v.typ == nilType {
		return false
	}
	if v.typ == boolType && !vtob(v) {
		return false
	}

	return true
}

func puts(vals ...Value) Value {
	for _, v := range vals {
		fmt.Print(v)
		fmt.Print(" ")
	}
	fmt.Print("\n")
	return Value{typ: nilType}
}

func add(vals ...Value) Value {
	var sum int64 = 0
	for _, v := range vals {
		switch v.typ {
		case floatType:
			return sumFloats(vals...)
		case intType:
			sum += vtoi(v)
		default:
			panic(fmt.Sprintf("+: Unexpected %s", v))
		}
	}
	return newInt(sum)
}

func sumFloats(vals ...Value) Value {
	sum := 0.0
	for _, v := range vals {
		switch v.typ {
		case intType:
			sum += float64(vtoi(v))
		case floatType:
			sum += vtof(v)
		default:
			panic(fmt.Sprintf("+: Unexpected  %s", v))
		}
	}
	return newFloat(sum)
}

func expectArgCount(name string, args []Value, expect int) {
	if len(args) != expect {
		arguments := "argument"
		if expect != 1 {
			arguments += "s"
		}

		err := fmt.Sprintf("Error: %s expects %d %s. Got %d.",
			name, expect, arguments, len(args))
		panic(err)
	}
}

func expectArg(name string, args []Value, index int, expect Type) {
	if args[index].typ != expect {
		err := fmt.Sprintf("Error: argument %d of %s should be of type %s.",
			index+1, name, expect)
		panic(err)
	}
}
