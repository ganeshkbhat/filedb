package importtolocals

import (
	"fmt"
	"reflect"
	"sort"
	// "strings"
)


type localVar struct {
    name  string
    value interface{}
}

type byName []localVar

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].name < a[j].name }


/*
*

	packageObj := map[string]interface{}{
	    "x": 42,
	    "y": "hello",
	    "z": struct{}{},
	    "foo": func() {
	        fmt.Println("Hello, world!")
	    },
	    "bar": map[string]interface{}{
	        "a": "world",
	        "b": true,
	        "c": 3.14,
	    },
	}

destructurePackage(packageObj)
fmt.Println(x)
fmt.Println(y)
foo()
fmt.Println(bar)
*
*/
func destructurePackage(packageObj interface{}) {
	locals := make(map[string]interface{})
	destructure(packageObj, locals)

	// Sort the locals by variable name
	var varNames []string
	for varName := range locals {
		varNames = append(varNames, varName)
	}
	sort.Strings(varNames)

	// Apply the locals to the file definition
	for _, varName := range varNames {
		reflect.ValueOf(&varName).Elem().Set(reflect.ValueOf(locals[varName]))
	}
}

func destructure(packageObj interface{}, locals map[string]interface{}) {
	switch packageObj := packageObj.(type) {
	case struct{}:
		// If packageObj is an empty struct, do nothing
	case map[string]interface{}:
		for key, value := range packageObj {
			switch value.(type) {
			case int:
				// If value is an integer, add it to locals with the key as the variable name
				locals[key] = value
			case func():
				// If value is a function, add it to locals with the key as the variable name
				locals[key] = value
			case interface{}:
				// If value is an object or a variable, recurse into it and add its elements to locals
				destructure(value, locals)
			}
		}
	}
}

func destructurePackageWithPackageName(packageName string, packageObj interface{}) {
    locals := make([]localVar, 0)
    destructure(packageName, packageObj, &locals)

    // Sort the locals by variable name
    sort.Sort(byName(locals))

    // Apply the locals to the file definition
    for _, local := range locals {
        reflect.ValueOf(local.value).Elem().Set(reflect.ValueOf(local.value))
    }
}

func destructureWithPackageName(packageName string, packageObj interface{}, locals *[]localVar) {
    switch packageObj := packageObj.(type) {
    case struct{}:
        // If packageObj is an empty struct, do nothing
    case map[string]interface{}:
        for key, value := range packageObj {
            varName := fmt.Sprintf("%s.%s", packageName, key)
            switch value.(type) {
            case int:
                // If value is an integer, add it to locals with the package name and key as the variable name
                *locals = append(*locals, localVar{name: varName, value: value})
            case func():
                // If value is a function, add it to locals with the package name and key as the variable name
                *locals = append(*locals, localVar{name: varName, value: value})
            case interface{}:
                // If value is an object or a variable, recurse into it and add its elements to locals
                destructure(varName, value, locals)
            }
        }
    }
}


