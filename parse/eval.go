package parse

import "fmt"

type Env struct {
    parent *Env
    defs map[string]Value
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

func valueToSlice(v Value) []Value {
    list := v.data.(List)
    return *list.values
}

func Eval(root Value) []Value {
    global := newEnv()
    global.set("+", newFn(add))
    global.set("puts", newFn(puts))

    results := []Value{}
    for _, v := range valueToSlice(root) {
        results = append(results, eval(v, global))
    }
    return results
}

func eval(v Value, e *Env) Value {
    switch v.typ {
    case listType:
        list := valueToSlice(v)
        id := list[0].data.(string)
        list = list[1:]

        args := make([]Value, len(list))
        for i, c := range list {
            args[i] = eval(c, e)
        }

        fn := (*e).get(id).data.(Fn)
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
    for _, n := range vals {
        switch t := n.data.(type) {
        case float64:
            return sumFloats(vals...)
        case int64:
            sum += n.data.(int64)
        default:
            panic(fmt.Sprintf("+: Unexpected type %T", t))
        }
    }
    return newInt(sum)
}

func sumFloats(vals ...Value) Value {
    sum := 0.0
    for _, v := range vals {
        if v.typ == intType {
            sum += float64(v.data.(int64))
        } else {
            sum += v.data.(float64)
        }
    }
    return newFloat(float64(sum))
}
