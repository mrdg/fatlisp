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

func (e Env) get(val Value) (Value, error) {
	id := val.data.(string)
	v, ok := e.defs[id]
	if !ok {
		if e.parent != nil {
			parent := *e.parent
			return parent.get(val)
		} else {
			err := newError(val.origin, "unable to resolve %s", id)
			return Value{}, err
		}
	}
	return v, nil
}

func Eval(root Value) ([]Value, error) {
	global := newEnv()
	global.set("+", newFn(add, -1))
	global.set("puts", newFn(puts, -1))
	global.set("def", newForm(def))
	global.set("fn", newForm(fn))
	global.set("if", newForm(_if))
	global.set("quote", newForm(quote))

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
		return e.get(v)
	case listType:
		return evalList(v, e)
	default:
		return v, nil
	}
}

func evalList(list Value, env *Env) (Value, error) {
	slice := vtos(list)
	args := make([]Value, len(slice)-1)
	id := slice[0]

	// Eval the first item in the list
	first, err := eval(id, env)
	if err != nil {
		return Value{}, err
	}

	// Check if the first item is a fn or a special form.
	// Returning an error if neither.
	switch first.typ {
	case fnType:
		fn := first.data.(Fn)

		// Loop over the rest of the list and eval each of
		// the function's arguments.
		for i, c := range slice[1:] {
			res, err := eval(c, env)
			if err != nil {
				return Value{}, err
			} else {
				args[i] = res
			}
		}

		if err := expectArgCount(id, args, fn.argc); err != nil {
			err := newError(id.origin, err.Error())
			return Value{}, err
		}
		return fn.fn(args...)
	case formType:
		form := first.data.(specialForm)
		return form(env, slice...)
	default:
		err := newError(first.origin, "not a function: %v", slice[0])
		return Value{}, err
	}
}

func quote(e *Env, vals ...Value) (Value, error) {
	quote := vals[0]
	vals = vals[1:]
	if err := expectArgCount(quote, vals, 1); err != nil {
		return Value{}, newError(quote.origin, err.Error())
	}
	return vals[0], nil
}

func fn(e *Env, vals ...Value) (Value, error) {
	fnForm := vals[0]
	vals = vals[1:] // Pop off fn keyword

	if err := expectArgCount(fnForm, vals, 2); err != nil {
		return Value{}, newError(fnForm.origin, err.Error())
	}
	if err := expectArg(fnForm, vals, 0, listType); err != nil {
		return Value{}, newError(fnForm.origin, err.Error())
	}

	params := vals[0]
	body := vals[1]

	fn := newFn(func(args ...Value) (Value, error) {
		res, err := eval(body, newFunctionEnv(e, params, args))
		if err != nil {
			return Value{}, err
		}
		return res, nil
	}, len(vtos(params)))

	return fn, nil
}

func def(e *Env, args ...Value) (Value, error) {
	def := args[0]
	args = args[1:]

	if err := expectArgCount(def, args, 2); err != nil {
		return Value{}, newError(def.origin, err.Error())
	}
	if err := expectArg(def, args, 0, idType); err != nil {
		return Value{}, newError(def.origin, err.Error())
	}

	id := args[0].data.(string)

	val, err := eval(args[1], e)
	if err != nil {
		return Value{}, err
	}
	e.set(id, val)
	return args[0], nil
}

func _if(env *Env, args ...Value) (Value, error) {
	var val Value
	var err error

	ifForm := args[0]
	args = args[1:]
	if err = expectArgCount(ifForm, args, 3); err != nil {
		return Value{}, newError(ifForm.origin, err.Error())
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
			err := newError(v.origin, "unexpected %s in +", v.typ)
			return Value{}, err
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
			err := newError(v.origin, "unexpected %s in +", v.typ)
			return Value{}, err
		}
	}
	return newFloat(sum), nil
}

func expectArgCount(val Value, args []Value, expect int) error {
	if expect != -1 && len(args) != expect {
		arguments := "argument"
		if expect != 1 {
			arguments += "s"
		}

		return fmt.Errorf("%s expected %d %s. Got %d.",
			val, expect, arguments, len(args))
	}
	return nil
}

func expectArg(val Value, args []Value, index int, expect Type) error {
	if args[index].typ != expect {
		return fmt.Errorf("argument %d of %s should be of type %s.",
			index+1, val, expect)
	}
	return nil
}
