package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strconv"
)


func callFunctionByName(name string, args ...interface{}) (interface{}, error) {
	// Get the function by name using reflection
	function := reflect.ValueOf(getFunctionByName(name))

	// Check if the function is valid and callable
	if !function.IsValid() || function.Kind() != reflect.Func {
		return nil, fmt.Errorf("function %s not found or not callable", name)
	}

	// Prepare the argument values as a slice of reflect.Value
	argValues := make([]reflect.Value, len(args))
	for i, arg := range args {
		argValues[i] = reflect.ValueOf(arg)
	}

	// Call the function with the provided arguments
	result := function.Call(argValues)

	// Extract and return the result as an interface{}
	if len(result) > 0 {
		return result[0].Interface(), nil
	}
	return nil, nil
}


func getFunctionByName(name string) interface{} {
	// Map the function names to their corresponding functions
	functions := map[string]interface{}{
		"hello": hello,
		"add":   add,
		// Add more function mappings as needed
	}

	// Return the function with the specified name, or nil if not found
	if function, ok := functions[name]; ok {
		return function
	}
	return nil
}


func hello(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}


func add(a, b int) int {
	return a + b
}


/** 

Get global function names

*/

func getFunctionNames(filePath string) ([]string, error) {
	fset := token.NewFileSet()

	// Parse the file
	file, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	var functionNames []string

	// Enumerate the global variables
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			// This is a function declaration
			functionNames = append(functionNames, decl.Name.Name)
		}
	}

	return functionNames, nil
}


func main() {
	

	result, err := callFunctionByName("hello", "Alice")
	if err != nil {
		// Handle the error
	}
	fmt.Println(result) // Output: Hello, Alice!

	result, err = callFunctionByName("add", 2, 3)
	if err != nil {
		// Handle the error
	}
	fmt.Println(result) // Output: 5

	// functions, err := getFunctionNames("example.go")
	// if err != nil {
	// 	panic(err)
	// }

	if len(os.Args) < 2 {
		errors.New("apply: not a function")
	}

	// Get the command line arguments.
	args := os.Args[1:]
	fmt.Println(args)

	a, _ := strconv.Atoi(args[1])
	b, _ := strconv.Atoi(args[2])
	// Call the function "myFunction"
	res, er := callFunctionByName("add", a, b)
	if er != nil {
		panic(er)
	}

	// Print the result
	fmt.Println(res)
}


