//
// Copyright (c) 2014 Canonical Ltd.
//
// Author: Jani Monoses <jani@ubuntu.com>
//
package commands

import (
	"fmt"
)

type clickCommand struct {
	name        string
	description string
	run         func(args []string)
}

var allCommands = []clickCommand{
	buildCommand,
	buildsourceCommand,
	contentsCommand,
	infoCommand,
}

func HelpText() string {
	t := ""
	for _, c := range allCommands {
		t += fmt.Sprintf("  %-22s%-60s\n", c.name, c.description)
	}
	return t
}

func RunCommand(args []string) bool {
	command := args[0]
	for _, c := range allCommands {
		if command == c.name {
			c.run(args[1:])
			return true
		}
	}
	return false
}

var chrootCommand = clickCommand{"chroot", "Use and manage a Click chroot.", nil}
var hookCommand = clickCommand{"hook", "Install or remove a Click system hook.", nil}
var installCommand = clickCommand{"install", "Install a Click package (low-level; consider pkcon instead).", nil}
var listCommand = clickCommand{"list", "List installed Click packages.", nil}
var pkgdirCommand = clickCommand{"pkgdir", "Print the directory where a Click package is unpacked.", nil}
var registerCommand = clickCommand{"register", "Register an installed Click package for a user.", nil}
var unregisterCommand = clickCommand{"unregister", "Unregister an installed Click package for a user.", nil}
var verifyCommand = clickCommand{"verify", "Verify a Click package.", nil}
