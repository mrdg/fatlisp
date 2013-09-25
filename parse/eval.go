package parse

import "fmt"

type Env struct {
	parent *Env
	defs   map[string]Value
}

func newEnv() *Env {
	return &Env{defs: make(map[string]Value, 1)}
}

func (e *Env) set(key string, v Value) {
	(*e).defs[key] = v
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
