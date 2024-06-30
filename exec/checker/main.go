package main

import (
	"os"

	"github.com/clkbug/konata-tools"
)

func run() error {
	if len(os.Args) < 2 {
		_, err := konata.ParseFile(os.Args[1])
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
