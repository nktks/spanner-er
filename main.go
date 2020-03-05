package main

import (
	"os"
)

func main() {
	cli := &cli{}
	os.Exit(cli.run(os.Args[1:]))
}
