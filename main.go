package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}

func run() error {
	return nil
}
