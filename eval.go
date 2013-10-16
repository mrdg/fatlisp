package fatlisp

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

func (e Env) get(id string) (Value, error) {
	v, ok := e.defs[id]
	if !ok {
		if e.parent != nil {
			parent := *e.parent
			return parent.get(id)
		} else {
			return Value{}, fmt.Errorf("Error: unable to resolve '%s'", id)
		}
	}
	return v, nil
}

func Eval(root Value) ([]Value, error) {
	global := newEnv()
	global.set("+", newFn(add))
	global.set("puts", newFn(puts))

	results := []Value{}
	for _, v := range vtos(root) {
		res, err := eval(v, global)
		if err != nil {
			return results, err
		} else {
			results = append(results, res)
		}
	}
	return results, nil
}

func eval(v Value, e *Env) (Value, error) {
	switch v.typ {
	case idType:
		id := v.data.(string)
		return e.get(id)
	case listType:
		list := vtos(v)
		args := make([]Value, len(list))

		if list[0].typ == idType {
			form, ok := getSpecialForm(list[0])
			if ok {
				return form(e, list...)
			}
		}

		for i, c := range list {
			res, err := eval(c, e)
			if err != nil {
				return Value{}, err
			} else {
				args[i] = res
			}
		}
		fn := args[0].data.(Fn)
		args = args[1:] // Pop of the function

		return fn(args...)
	default:
		return v, nil
	}
}

type specialForm func(env *Env, args ...Value) (Value, error)

func getSpecialForm(v Value) (specialForm, bool) {
	name := v.data.(string)
	switch name {
	case "def":
		return def, true
	case "fn":
		return fn, true
	case "if":
		return _if, true
	case "quote":
		return quote, true
	default:
		return nil, false
	}
}

func quote(e *Env, vals ...Value) (Value, error) {
	return vals[1], nil
}

func fn(e *Env, vals ...Value) (Value, error) {
	vals = vals[1:] // Pop off fn keyword

	if err := expectArgCount("fn", vals, 2); err != nil {
		return Value{}, err
	}
	if err := expectArg("fn", vals, 0, listType); err != nil {
		return Value{}, err
	}

	params := vals[0]
	body := vals[1]

	return newFn(func(args ...Value) (Value, error) {
		argc := len(vtos(params))
		if err := expectArgCount("function", args, argc); err != nil {
			return Value{}, err
		}
		res, err := eval(body, newFunctionEnv(e, params, args))
		if err != nil {
			return Value{}, err
		}
		return res, nil
	}), nil
}

func def(e *Env, args ...Value) (Value, error) {
	name := args[1].data.(string)
	val, err := eval(args[2], e)
	if err != nil {
		return Value{}, err
	}
	e.set(name, val)
	return args[1], nil
}

func _if(env *Env, args ...Value) (Value, error) {
	var val Value
	var err error

	args = args[1:]
	if err = expectArgCount("if", args, 3); err != nil {
		return Value{}, err
	}

	cond := args[0]
	ifClause := args[1]
	elseClause := args[2]

	val, err = eval(cond, env)
	if err != nil {
		return Value{}, err
	}

	if truthy(val) {
		val, err = eval(ifClause, env)
	} else {
		val, err = eval(elseClause, env)
	}
	if err != nil {
		return Value{}, err
	}
	return val, nil
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

func puts(vals ...Value) (Value, error) {
	for _, v := range vals {
		fmt.Print(v)
		fmt.Print(" ")
	}
	fmt.Print("\n")
	return Value{typ: nilType}, nil
}

func add(vals ...Value) (Value, error) {
	var sum int64 = 0
	for _, v := range vals {
		switch v.typ {
		case floatType:
			return sumFloats(vals...)
		case intType:
			sum += vtoi(v)
		default:
			return Value{}, fmt.Errorf("+: Unexpected %s", v)
		}
	}
	return newInt(sum), nil
}

func sumFloats(vals ...Value) (Value, error) {
	sum := 0.0
	for _, v := range vals {
		switch v.typ {
		case intType:
			sum += float64(vtoi(v))
		case floatType:
			sum += vtof(v)
		default:
			return Value{}, fmt.Errorf("+: Unexpected  %s", v)
		}
	}
	return newFloat(sum), nil
}

func expectArgCount(name string, args []Value, expect int) error {
	if len(args) != expect {
		arguments := "argument"
		if expect != 1 {
			arguments += "s"
		}

		err := fmt.Errorf("Error: %s expects %d %s. Got %d.",
			name, expect, arguments, len(args))
		return err
	}
	return nil
}

func expectArg(name string, args []Value, index int, expect Type) error {
	if args[index].typ != expect {
		err := fmt.Errorf("Error: argument %d of %s should be of type %s.",
			index+1, name, expect)
		return err
	}
	return nil
}
