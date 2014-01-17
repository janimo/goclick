//
// Copyright (c) 2014 Canonical Ltd.
//
// Author: Jani Monoses <jani@ubuntu.com>
//
package main

import (
	"fmt"
	"os"

	"github.com/janimo/goclick/commands"
)

func printUsage() {
	fmt.Println("Usage: goclick COMMAND [options]\n")
	fmt.Println("Commands are as follows ('click COMMAND --help' for more):\n")
	fmt.Printf("%s", commands.HelpText())
	os.Exit(1)
}

func main() {
	if len(os.Args) == 1 {
		printUsage()
	}
	if !commands.RunCommand(os.Args[1:]) {
		printUsage()
	}
}
