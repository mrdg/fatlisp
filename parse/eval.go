package parse
// 
// import "fmt"
// 
// type Env struct {
//     parent *Env
//     defs map[string]lispValue
// }
// 
// type Node struct {
//     value lispValue
//     children []Node
// }
// 
// func eval(node Node) lispValue {
//     // TODO: handle special forms
// 
//     if len(node.children) > 0 {
//         args := make([]lispValue, len(node.children))
//         for i, c := range node.children {
//             args[i] = eval(c)
//         }
//         fn, _ := node.value.(lispFunc)
//         return fn.call(args...)
//     } else {
//         return node.value
//     }
// }
// 
// func add(args... lispValue) lispValue {
//     returnFloat := false
//     fsum := 0.0
//     isum := 0
// 
//     for _, n := range args {
//         switch t := n.(type) {
//         case lispFloat:
//             returnFloat = true
//             fsum += n.(lispFloat).value
//         case lispInt:
//             isum += n.(lispInt).value
//         default:
//             panic(fmt.Sprintf("+: Unexpected type %T", t))
//         }
//     }
// 
//     if returnFloat {
//         fsum += float64(isum)
//         return lispFloat{fsum}
//     } else {
//         return lispInt{isum}
//     }
// }
