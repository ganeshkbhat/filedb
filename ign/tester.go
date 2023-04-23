package main

// import (
// 	"errors"
// 	"fmt"
// 	"reflect"
// 	"os"
// 	"go/ast"
// 	"go/parser"
// 	"go/token"
// 	"strconv"
// )

// //
// func add(a, b int) int {
// 	return a + b
// }

// //
// func getFunctionNames(filePath string) []string {
// 	fset := token.NewFileSet()

// 	// Parse the file
// 	file, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var functionNames []string
// 	// Enumerate the global variables
// 	for _, decl := range file.Decls {
// 		switch decl := decl.(type) {
// 		case *ast.FuncDecl:
// 			// This is a function declaration
// 			functionNames = append(functionNames, decl.Name.Name)
// 		}
// 	}
// 	return functionNames
// }

// //
// func callGlobalFunction(funcName string, args ...interface{}) (interface{}, error) {
// 	// Get a reference to the function by name
// 	funcValue := reflect.ValueOf(nil)
// 	funcType := reflect.TypeOf(nil)
// 	for _, global := range getFunctionNames("tester.go") {
// 		if global == funcName {
// 			funcValue, funcType = reflect.ValueOf(global), reflect.TypeOf(global)
// 			fmt.Println("global, funcValue, funcType: ", global, funcValue, funcType, reflect.TypeOf(funcValue))
// 			break
// 		}
// 	}

// 	fmt.Println("funcValue, funcType: ", funcValue, funcType)
// 	if funcValue == reflect.ValueOf(nil) {
// 		return nil, errors.New("function not found")
// 	}

// 	// Check the function arguments
// 	if funcType.Kind() != reflect.Func {
// 		return nil, errors.New("not a function")
// 	}

// 	// if funcType.NumIn() != len(args) {
// 	// 	return nil, errors.New("wrong number of arguments")
// 	// }

// 	// for i, arg := range args {
// 	// 	if !reflect.TypeOf(arg).AssignableTo(funcType.In(i)) {
// 	// 		return nil, fmt.Errorf("argument %d type mismatch", i)
// 	// 	}
// 	// }

// 	// in := make([]reflect.Value, len(args))
// 	// // // Convert to string
// 	// for i := range args {
// 	// 	v, err := strconv.Atoi(args[i].(string))
// 	// 	if err != nil {
// 	// 		fmt.Errorf("apply: incorrect value of argument needed for function")
// 	// 	}
// 	// 	in[i] = reflect.ValueOf(v)
// 	// }

// 	// // Call the function
// 	// result := funcValue.Call(in)
// 	// if len(result) == 0 {
// 	// 	return nil, nil
// 	// }
// 	// return result[0].Interface(), nil

// 	return funcValue.Interface(), nil
// }

// // callFunction by getting global names of variables
// func callFunctions(filename string, funcName string, args ...interface{}) (interface{}, error) {
// 	// Get a reference to the function by name
// 	funcValue := reflect.ValueOf(nil)
// 	funcType := reflect.TypeOf(nil)
// 	for _, global := range getFunctionNames(filename) {
// 		if global == funcName {
// 			funcValue, funcType = reflect.ValueOf(global), reflect.TypeOf(global)
// 			break
// 		}
// 	}

// 	if funcValue == reflect.ValueOf(nil) {
// 		return nil, errors.New("function not found")
// 	}

// 	// Check the function arguments
// 	if funcType.NumIn() != len(args) {
// 		return nil, errors.New("wrong number of arguments")
// 	}
// 	inArgs := make([]reflect.Value, len(args))
// 	for i, arg := range args {
// 		inArgs[i] = reflect.ValueOf(arg)
// 		if !inArgs[i].IsValid() || !inArgs[i].Type().AssignableTo(funcType.In(i)) {
// 			return nil, fmt.Errorf("argument %d type mismatch", i)
// 		}
// 	}
// 	// Call the function
// 	result := funcValue.Call(inArgs)
// 	if len(result) == 0 {
// 		return nil, nil
// 	}
// 	return result[0].Interface(), nil
// }

// func GetGlobalFunction(funcName string) (interface{}, error) {
//     // Get the value of the global symbol with the given name
//     funcValue := reflect.ValueOf(funcName)

//     // Check if the symbol is a function
//     if funcValue.Kind() != reflect.Func {
//         return nil, fmt.Errorf("%s is not a function", funcName)
//     }

//     // Get the type of the function
//     funcType := funcValue.Type()

//     // Check if the function is defined in the global scope
//     if funcType.PkgPath() != "" {
//         return nil, fmt.Errorf("%s is not a global function", funcName)
//     }

//     // Return the function value
//     return funcValue.Interface(), nil
// }

// func main() {

// 	if len(os.Args) < 2 {
// 		errors.New("apply: function not provided")
// 	}

// 	// // Get the command line arguments.
// 	args := os.Args[1:]
// 	fmt.Println("args: ", args)

// 	in := make([]reflect.Value, len(args))
// 	// // Convert to string
// 	for i := range args {
// 		v, err := strconv.Atoi(args[i])
// 		if err != nil {
// 			fmt.Errorf("apply: incorrect value of argument needed for function")
// 		}
// 		in[i] = reflect.ValueOf(v)
// 	}

// 	result, err := GetGlobalFunction("add")
// 	if err != nil {
// 		fmt.Println("Error calling function:", err)
// 	}

// 	fmt.Println("Result of function call: result, err:", result, err)

// 	fnc, ok := result.(func())
// 	if !ok {
// 		fmt.Println("Error calling function:", ok)
// 	}

// 	r, o := fnc(1 ,2)
// 	if !o {
// 		fmt.Println("Error calling function:", o)
// 	}
// 	fmt.Println("Result of function call: fnc", r, fnc)

// 	// // GET FUNCTION STRING FUNCTION
// 	// result := getFunctionNames("tester.go")
// 	// fmt.Println(result)

// 	// // GET FUNCTIONS
// 	// args := []int{1, 2}
// 	// fmt.Println(callGlobalFunction("add", args))
// }

import (
	"fmt"
	// "os"
	"reflect"
	// "errors"
	// "strconv"
)

func myGlobalFunction() {
	fmt.Println("Hello, world!")
}

func GetGlobalFunction(funcName string) (interface{}, error) {
	// Get the value of the global symbol with the given name
	funcValue := reflect.ValueOf(myGlobalFunction)

	// Check if the symbol is a function
	if funcValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("%s is not a function", funcName)
	}

	// Get the type of the function
	funcType := funcValue.Type()

	// Check if the function is defined in the global scope
	if funcType.PkgPath() != "" {
		return nil, fmt.Errorf("%s is not a global function", funcName)
	}

	// Return the function value
	return funcValue.Interface(), nil
}

func add(a, b int) int {
	return a + b
}

func main() {

	// if len(os.Args) < 2 {
	// 	errors.New("apply: not a function")
	// }
	// // Get the command line arguments.
	// args := os.Args[1:]

	// Call GetGlobalFunction to retrieve the myGlobalFunction function
	myFunc, err := GetGlobalFunction("myGlobalFunction")
	if err != nil {
		fmt.Println(err)
	}

	// Cast the function value to a function with the appropriate signature
	fn, ok := myFunc.(func())
	if !ok {
		fmt.Println("Function type assertion failed")
	}

	// a, _ := strconv.Atoi(args[1])
	// b, e := strconv.Atoi(args[2])
	// if e != nil {
	// 	fmt.Println(e)
	// }
	// fmt.Println(a, b)

	fmt.Println(reflect.TypeOf(fn), reflect.ValueOf(fn))

	// // call the function
	fn()

	// // Call the function
	// res := fn(a, b)
	// fmt.Println(res)
}
