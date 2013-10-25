package fatlisp

import "fmt"

func vtos(v Value) []Value {
	list := v.data.(List)
	return *list.values
}

func vtoi(v Value) Int {
	return v.data.(Int)
}

func vtof(v Value) Float {
	return v.data.(Float)
}

func vtob(v Value) bool {
	return v.data.(bool)
}

func vtofn(v Value) *Fn {
	return v.data.(*Fn)
}

func vtoform(v Value) *specialForm {
	return v.data.(*specialForm)
}

func val2str(v Value) string {
	return v.data.(string)
}

func val2num(v Value) Number {
	switch v.typ {
	case intType:
		return v.data.(Int)
	case floatType:
		return v.data.(Float)
	default:
		panic("Can't convert value to num")
	}
}

func num2val(n Number) Value {
	if n.isFloat() {
		return newFloat(n.toFloat())
	}
	return newInt(n.toInt())
}

func checkTypes(vals []Value, types ...Type) error {
Outer:
	for _, v := range vals {
		for _, t := range types {
			if v.typ == t {
				continue Outer
			}
		}
		return typeError(v)
	}

	return nil
}

func isNumeric(v Value) bool {
	return v.typ == intType || v.typ == floatType
}

func typeError(v Value) error {
	return newError(v.origin, "unexpected type %s", v.typ)
}

func newError(origin item, msg string, args ...interface{}) error {
	msg = fmt.Sprintf(msg, args...)
	return fmt.Errorf("%s %s", origin.pos, msg)
}
