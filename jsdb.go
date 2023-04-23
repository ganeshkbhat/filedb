package main

// package jsdb

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

func main() {
	if len(os.Args) < 2 {
		errors.New("apply: not a function")
	}

	// Get the command line arguments.
	args := os.Args[1:]
	fmt.Println(args)

	// Call the function specified in the first argument with the rest of the arguments.
	if len(args) > 0 {
		fmt.Printf("%v\n", reflect.TypeOf(args[0]))
		// fmt.Println(val)
	}
}
