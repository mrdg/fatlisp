package fatlisp

import "fmt"

func equal(vals ...Value) (Value, error) {
	x := vals[0]
	y := vals[1]

	if x.typ != y.typ {
		return int2val(-1), nil
	}

	var eq bool
	switch x.typ {
	case intType:
		eq = val2int(x) == val2int(y)
	case floatType:
		eq = val2float(x) == val2float(y)
	case stringType:
	case idType:
		eq = val2str(x) == val2str(y)
	case listType:
		return Value{}, fmt.Errorf("TODO: implement equals for lists")
	case nilType:
		eq = true
	case boolType:
		eq = val2bool(x) == val2bool(y)
	case fnType:
		xp := val2fn(x)
		yp := val2fn(y)
		eq = &xp == &yp
	}

	return bool2val(eq), nil
}

func compare(vals ...Value) (Value, error) {
	x := vals[0]
	y := vals[1]
	return x.compare(y)
}

func (x Value) compare(y Value) (Value, error) {
	if x.typ != y.typ {
		if !isNumeric(x) || !isNumeric(y) {
			return int2val(-1), nil
		}
	}
	switch x.typ {
	case intType:
		return val2int(x).compare(val2num(y)), nil
	case floatType:
		return val2float(x).compare(val2num(y)), nil
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
		return int2val(1)
	} else if a == b {
		return int2val(0)
	} else {
		return int2val(-1)
	}
}

func (x Float) compare(y Number) Value {
	a := x
	b := y.toFloat()
	if a > b {
		return int2val(1)
	} else if a == b {
		return int2val(0)
	} else {
		return int2val(-1)
	}
}

func (x String) compare(y String) Value {
	if x > y {
		return int2val(1)
	} else if x == y {
		return int2val(0)
	} else {
		return int2val(-1)
	}
}
