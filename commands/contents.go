//
// Copyright (c) 2014 Canonical Ltd.
//
// Author: Jani Monoses <jani@ubuntu.com>
//
package commands

import (
	"fmt"
	"os/exec"

	"github.com/janimo/goclick/utils"
)

var contentsCommand = clickCommand{"contents", "Show the file-list contents of a Click package file.", runContents}

func runContents(args []string) {
	if len(args) == 0 {
		utils.ExitMsg("error: need file")
	}

	clickfile := args[0]
	if !utils.Exists(clickfile) {
		utils.ExitMsg(clickfile + " does not exist")
	}
	cmd := exec.Command("dpkg-deb", "-c", clickfile)
	out, err := cmd.Output()
	if err != nil {
		utils.ExitError(err)
	}
	fmt.Printf("%s", out)
}
