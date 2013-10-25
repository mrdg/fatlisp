package fatlisp

import "fmt"

func equal(vals ...Value) (Value, error) {
	x := vals[0]
	y := vals[1]

	if x.typ != y.typ {
		return newInt(-1), nil
	}

	var eq bool
	switch x.typ {
	case intType:
		eq = vtoi(x) == vtoi(y)
	case floatType:
		eq = vtof(x) == vtof(y)
	case stringType:
	case idType:
		eq = val2str(x) == val2str(y)
	case listType:
		return Value{}, fmt.Errorf("TODO: implement equals for lists")
	case nilType:
		eq = true
	case boolType:
		eq = vtob(x) == vtob(y)
	case fnType:
		xp := vtofn(x)
		yp := vtofn(y)
		eq = &xp == &yp
	}

	return newBool(eq), nil
}

func compare(vals ...Value) (Value, error) {
	x := vals[0]
	y := vals[1]
	return x.compare(y)
}

func (x Value) compare(y Value) (Value, error) {
	if x.typ != y.typ {
		if !isNumeric(x) || !isNumeric(y) {
			return newInt(-1), nil
		}
	}
	switch x.typ {
	case intType:
		return vtoi(x).compare(val2num(y)), nil
	case floatType:
		return vtof(x).compare(val2num(y)), nil
	case stringType:
		return String(val2str(x)).compare(String(val2str(y))), nil
	default:
		return Value{}, newError(x.origin, "can't compare type %s", x.typ)
	}
}

func (x Int) compare(y Number) Value {
	a := Float(x)
	b := y.toFloat()
	if a > b {
		return newInt(1)
	} else if a == b {
		return newInt(0)
	} else {
		return newInt(-1)
	}
}

func (x Float) compare(y Number) Value {
	a := x
	b := y.toFloat()
	if a > b {
		return newInt(1)
	} else if a == b {
		return newInt(0)
	} else {
		return newInt(-1)
	}
}

func (x String) compare(y String) Value {
	if x > y {
		return newInt(1)
	} else if x == y {
		return newInt(0)
	} else {
		return newInt(-1)
	}
}
