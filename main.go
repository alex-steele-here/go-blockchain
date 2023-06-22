package main

import (
	"os"

	"github.com/alex-steele-here/go-blockchain.git/cli"
)

func main() {
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()
}
