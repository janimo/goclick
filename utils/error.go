//
// Copyright (c) 2014 Canonical Ltd.
//
// Author: Jani Monoses <jani@ubuntu.com>
//
package utils

func ExitMsg(msg string) {
	panic(msg)
}
func ExitError(err error) {
	panic(err)
}
