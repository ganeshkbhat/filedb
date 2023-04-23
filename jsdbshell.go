package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
    "syscall"
)

func sqlite3(filepath string, command string, args ...interface{}) {

}

func document(filepath string, command string, args ...interface{}) {
	fmt.Printf(filepath)

	// if filetype === csv

	// if filetype === json

	// if filetype === yaml

	// if filetype === text

	// fmt.Printf(filetype)
}

func responses(sql string, str string) string {
	return str + "response result "
}

func JsdbShell() {
	// Create a buffered writer for stdin
	scanner := bufio.NewScanner(os.Stdin)

	var text string
	// // Create a buffered writer for stdout
    // text := bufio.NewWriter(os.Stdout)

	// Set up a signal handler to catch SIGINT
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT)
	
	// break the loop if text == "exit();"
	for text != "exit();" {
		fmt.Print("jsdbsql> ")
		scanner.Scan()
		sql := scanner.Text()
		text = responses(sql, ">>\n")
		// Check for SIGINT
        select {
			case <-sig:
				fmt.Println("Caught SIGINT, exiting...")
				return
			default:
				// No signal, continue writing
				if text != "exit();" {
					fmt.Println(text, "  \n")
				}
        }
	}
}

func main() {
	JsdbShell()
}
