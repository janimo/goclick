//
// Copyright (c) 2014 Canonical Ltd.
//
// Author: Jani Monoses <jani@ubuntu.com>
//
package utils

//Preinst for Click packages.

//In general there is a rule that Click packages may not have maintainer
//scripts.  However, there is one exception: a static preinst used to cause
//dpkg to fail if people attempt to install Click packages directly using dpkg
//rather than via "click install".  This avoids accidents, since Click
//packages use a different root of their filesystem tarball.

const StaticPreinst = `#! /bin/sh
echo "Click packages may not be installed directly using dpkg."
echo "Use 'click install' instead."
exit 1
`
