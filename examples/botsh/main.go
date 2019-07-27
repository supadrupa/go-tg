package main

import (
	"fmt"
	"os"

	"github.com/mr-linch/go-tg/examples/botsh/cmd"
)

func main() {
	if err := cmd.NewApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
