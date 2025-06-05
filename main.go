package main

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

// run executes all tasks and returns an OS exit code
func run() (code int) {
	if len(os.Args) < 2 {
		code = 1
		fmt.Fprintf(os.Stderr, "%v\n", errors.New("requires an argument"))
		return
	}

	root := os.Args[1]

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		ProcessFile(root, &wg, errCh)
	}()
	
	wg.Wait()
	close(errCh)
	
	for err := range errCh {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		code = 2
	}

	return
}

func main() {
	code := run()
	os.Exit(code)
}
