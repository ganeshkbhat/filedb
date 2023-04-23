package commands

// package commands

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"go/ast"
	"go/parser"
	"go/token"
	// _ "github.com/mattn/go-sqlite3"
)


func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}


// func getFunctionInterface(i interface{}) string {
// 	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Interface{}
// }

/*
	func convertExample() {
		// Example input values
		inputString := "hello world"
		inputInt := "123"
		inputMap := `{"name": "Alice", "age": 30}`

		// Convert the inputs to valid types
		outputString := convertToValidType(inputString).(string)
		outputInt := convertToValidType(inputInt).(int)
		outputMap := convertToValidType(inputMap).(map[string]interface{})

		// Print the results
		fmt.Printf("outputString: %v, type: %T\n", outputString, outputString)
		fmt.Printf("outputInt: %v, type: %T\n", outputInt, outputInt)
		fmt.Printf("outputMap: %v, type: %T\n", outputMap, outputMap)
	}
*/
func ConvertToValidType(input string) interface{} {
	// Try to convert the input to an integer
	if intValue, err := strconv.Atoi(input); err == nil {
		return intValue
	}

	// Try to parse the input as a JSON string to create a map
	var mapValue map[string]interface{}
	if err := json.Unmarshal([]byte(input), &mapValue); err == nil {
		return mapValue
	}

	// If all else fails, just return the input as a string
	return input
}


// // get function using static map to call using function name
func GetFunctionStaticMap(name string) interface{} {
	// Map the function names to their corresponding functions
	functions := map[string]interface{}{
		"add": add,
		// Add more function mappings as needed
	}
	// Return the function with the specified name, or nil if not found
	if function, ok := functions[name]; ok {
		return function
	}
	return nil
}


// // get function using static switch to call using function name
func getFuncSwitch(name string) interface{} {
	switch name {
	case "add":
		return add
	case "apply":
		return apply
	default:
		return nil
	}
}


// // apply - call function from file
func apply(filename string, funcName interface{}, args []string) error {

	// fnValue := reflect.ValueOf(nil)
	// funcType := reflect.TypeOf(nil)

	// fnValue := reflect.ValueOf(funcName)
	// if fnValue.Kind() != reflect.Func {
	// 	return fmt.Errorf("apply: not a function")
	// }

	// for _, global := range getFileFunctionNames(filename) {
	// 	fmt.Println("global", global)
	// 	if global == funcName {
	// 		fnValue, funcType = reflect.ValueOf(global), reflect.TypeOf(global)
	// 		fmt.Println("funcName, fnValue, funcType: ", funcName, fnValue, funcType)
	// 		break
	// 	}
	// }

	// // reflect.ValueOf(fn) is a string and giving errors
	fnValue := reflect.ValueOf(funcName)
	// f := reflect.ValueOf(getFunc(fn.(string)))

	// fmt.Println(funcName, funcType)

	// working with asking for function manually
	// fnValue := reflect.ValueOf(getFuncSwitch(funcName.(string)))

	if fnValue.Kind() != reflect.Func {
		return fmt.Errorf("apply: not a function")
		// return errors.New("apply: not a function")
	}

	fnType := fnValue.Type()
	if len(args) != fnType.NumIn() {
		return fmt.Errorf("apply: incorrect number of arguments")
	}

	in := make([]reflect.Value, len(args))

	// // Convert to string
	for i := range args {
		v, err := strconv.Atoi(args[i])
		if err != nil {
			fmt.Errorf("apply: incorrect value of argument needed for function")
		}
		in[i] = reflect.ValueOf(v)
	}

	// // Check string value
	// for i, arg := range args {
	//     in[i] = reflect.ValueOf(arg)
	// }

	out := fnValue.Call(in)
	fmt.Println(len(args), args)

	// fmt.Println(r)
	// // Call the function with the given arguments
	// out := fnValue.Call(in)
	// // Print the result
	fmt.Printf("%s returned %s with a value of %v\n", fnType, out[0], out[0])

	return nil
}


// // getFileFunctionNames from file.go
func getFileFunctionNames(filePath string) []string {
	fset := token.NewFileSet()

	// Parse the file
	file, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		panic(err)
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
	return functionNames
}


// // callFunction by getting global names of variables
func callFunction(filename string, funcName string, args ...interface{}) (interface{}, error) {
	// Get a reference to the function by name
	funcValue := reflect.ValueOf(nil)
	funcType := reflect.TypeOf(nil)
	for _, global := range getFileFunctionNames(filename) {
		if global == funcName {
			funcValue, funcType = reflect.ValueOf(global), reflect.TypeOf(global)
			break
		}
	}
	if funcValue == reflect.ValueOf(nil) {
		return nil, errors.New("function not found")
	}
	// Check the function arguments
	if funcType.NumIn() != len(args) {
		return nil, errors.New("wrong number of arguments")
	}
	inArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		inArgs[i] = reflect.ValueOf(arg)
		if !inArgs[i].IsValid() || !inArgs[i].Type().AssignableTo(funcType.In(i)) {
			return nil, fmt.Errorf("argument %d type mismatch", i)
		}
	}
	// Call the function
	result := funcValue.Call(inArgs)
	if len(result) == 0 {
		return nil, nil
	}
	return result[0].Interface(), nil
}


// // add two numbers
func add(a, b int) int {
	return a + b
}


// func main() {
// 	if len(os.Args) < 2 {
// 		errors.New("apply: not a function")
// 	}

// 	// Get the command line arguments.
// 	args := os.Args[1:]

// 	fmt.Println(args)
// 	// Call the function specified in the first argument with the rest of the arguments.
// 	if len(args) > 0 {
// 		fmt.Printf("%v\n", reflect.TypeOf(args[0]))
// 		err := apply("sql.go", args[0], args[1:])
// 		// val, err := callFunction("sql.go", args[0], args[1:])
// 		if err != nil {
// 			fmt.Printf("Test: %s\n", err)
// 		}
// 		// fmt.Println(val)
// 	}
// }

