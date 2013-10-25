package fatlisp

type Number interface {
	add(y Number) Value
	subtract(y Number) Value
	multiply(y Number) Value
	divide(y Number) Value
	toFloat() Float
	toInt() Int
	isFloat() bool
}

func (i Int) toInt() Int {
	return i
}

func (i Int) isFloat() bool {
	return false
}

func (i Int) toFloat() Float {
	return Float(i)
}

func (f Float) isFloat() bool {
	return true
}

func (f Float) toFloat() Float {
	return f
}

func (f Float) toInt() Int {
	return Int(f)
}

func (x Int) add(y Number) Value {
	if y.isFloat() {
		return num2val(x.toFloat() + y.toFloat())
	}
	return num2val(x + y.toInt())
}

func (x Float) add(y Number) Value {
	return num2val(x + y.toFloat())
}

func add(vals ...Value) (Value, error) {
	err := checkTypes(vals, intType, floatType)
	if err != nil {
		return Value{}, err
	}
	x := val2num(vals[0])
	y := val2num(vals[1])
	return x.add(y), nil
}

func (x Int) subtract(y Number) Value {
	if y.isFloat() {
		return num2val(x.toFloat() - y.toFloat())
	}
	return num2val(x - y.toInt())
}

func (x Float) subtract(y Number) Value {
	return num2val(x - y.toFloat())
}

func subtract(vals ...Value) (Value, error) {
	err := checkTypes(vals, intType, floatType)
	if err != nil {
		return Value{}, err
	}
	x := val2num(vals[0])
	y := val2num(vals[1])
	return x.subtract(y), nil
}

func (x Int) multiply(y Number) Value {
	if y.isFloat() {
		return num2val(x.toFloat() * y.toFloat())
	}
	return num2val(x * y.toInt())
}

func (x Float) multiply(y Number) Value {
	return num2val(x * y.toFloat())
}

func multiply(vals ...Value) (Value, error) {
	err := checkTypes(vals, intType, floatType)
	if err != nil {
		return Value{}, err
	}
	x := val2num(vals[0])
	y := val2num(vals[1])
	return x.multiply(y), nil
}

func (x Int) divide(y Number) Value {
	if y.isFloat() {
		return num2val(x.toFloat() / y.toFloat())
	}
	return num2val(x / y.toInt())
}

func (x Float) divide(y Number) Value {
	return num2val(x / y.toFloat())
}

func divide(vals ...Value) (Value, error) {
	err := checkTypes(vals, intType, floatType)
	if err != nil {
		return Value{}, err
	}
	x := val2num(vals[0])
	y := val2num(vals[1])
	return x.divide(y), nil
}
