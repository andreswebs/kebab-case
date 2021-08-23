package main

import (
	"errors"
	"fmt"
	"os"
)

// run executes all tasks and returns an OS exit code and error
func run() (code int, err error) {

	if len(os.Args) < 2 {
		return 1, errors.New("requires an argument")
	}

	start := os.Args[1]

	err = ProcessFile(start)
	if err != nil {
		return 1, err
	}

	return

}

func main() {
	code, err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	os.Exit(code)
}
