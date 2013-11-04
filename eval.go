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
		name := val2str(p)
		value := args[i]
		env.set(name, value)
	}
	return env
}

func (e *Env) set(key string, v Value) {
	e.defs[key] = v
}

func (e Env) get(val Value) (Value, error) {
	id := val2str(val)
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
	global.set("add", newFn(add, 2, 2))
	global.set("subtract", newFn(subtract, 2, 2))
	global.set("multiply", newFn(multiply, 2, 2))
	global.set("divide", newFn(divide, 2, 2))
	global.set("compare", newFn(compare, 2, 2))
	global.set("puts", newFn(puts, 0, -1))

	global.set("def", newForm("def", def, 2, 2, []Type{idType}))
	global.set("fn", newForm("fn", fn, 2, 2, []Type{listType}))
	global.set("if", newForm("if", _if, 2, 3, []Type{}))
	global.set("quote", newForm("quote", quote, 1, 1, []Type{}))

	results := []Value{}
	for _, v := range val2slice(root) {
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
	slice := val2slice(list)
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

		fn := val2fn(first)
		if err := validateFnArgs(fn, args); err != nil {
			err := newError(id.origin, err.Error())
			return Value{}, err
		}
		val, err := fn.fn(args...)
		if err != nil {
			return val, err
		}
		return val, err
	case formType:
		form := val2form(first)
		if err := validateFormArgs(form, slice[1:]); err != nil {
			err := newError(id.origin, err.Error())
			return Value{}, err
		}
		return form.fn(env, slice...)
	default:
		err := newError(first.origin, "not a function: %v", slice[0])
		return Value{}, err
	}
}

func quote(e *Env, vals ...Value) (Value, error) {
	return vals[1], nil
}

func fn(e *Env, vals ...Value) (Value, error) {
	vals = vals[1:] // Pop off fn keyword

	params := vals[0]
	body := vals[1]

	min := len(val2slice(params))
	max := min
	fn := newFn(func(args ...Value) (Value, error) {
		res, err := eval(body, newFunctionEnv(e, params, args))
		if err != nil {
			return Value{}, err
		}
		return res, nil
	}, min, max)

	return fn, nil
}

func def(e *Env, args ...Value) (Value, error) {
	args = args[1:]

	id := val2str(args[0])

	val, err := eval(args[1], e)
	if err != nil {
		return Value{}, err
	}

	// if value is a fn, set the name on its signature
	// so it can be displayed in error messages
	if val.typ == fnType {
		fn := val2fn(val)
		fn.sig.name = id
	}

	e.set(id, val)
	return args[0], nil
}

func _if(env *Env, args ...Value) (Value, error) {
	var val Value
	var err error

	args = args[1:]
	val, err = eval(args[0], env)
	if err != nil {
		return Value{}, err
	}

	if truthy(val) {
		val, err = eval(args[1], env)
	} else if len(args) > 2 {
		val, err = eval(args[2], env)
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
	if v.typ == boolType && !val2bool(v) {
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

func validateFnArgs(fn *Fn, args []Value) error {
	if err := validateArgCount(fn.sig, args); err != nil {
		return err
	}
	return nil
}

func validateFormArgs(form *specialForm, args []Value) error {
	if err := validateArgCount(form.sig, args); err != nil {
		return err
	}
	if err := validateArgTypes(form.sig, args); err != nil {
		return err
	}
	return nil
}

func validateArgCount(sig signature, args []Value) error {
	argc := len(args)
	if argc < sig.minArgs {
		return argCountError(sig.name, sig.minArgs, argc)
	}
	if sig.maxArgs != -1 && argc > sig.maxArgs {
		return argCountError(sig.name, sig.maxArgs, argc)
	}
	return nil
}

func validateArgTypes(sig signature, args []Value) error {
	for i, typ := range sig.types {
		if typ != args[i].typ {
			err := fmt.Errorf("argument %d of %s should be %s, got %s",
				i+1, sig.name, typ, args[i].typ)
			return err
		}
	}
	return nil
}

func argCountError(name string, expected, actual int) error {
	arguments := "argument"
	if expected != 1 {
		arguments += "s"
	}

	return fmt.Errorf("%s expected %d %s, got %d",
		name, expected, arguments, actual)
}
