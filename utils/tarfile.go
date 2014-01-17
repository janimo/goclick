//
// Copyright (c) 2014 Canonical Ltd.
//
// Author: Jani Monoses <jani@ubuntu.com>
//
package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"strings"
)

// Adds a file to the tar archive tw
func tarFile(root, name string, tw *tar.Writer, fakeroot bool) {
	fi, err := os.Stat(name)
	if err != nil {
		ExitError(err)
	}

	//Write header. TODO symlinks
	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		ExitError(err)
	}

	hdr.Name = strings.Replace(name, root, "./", 1)
	//Force root owner and group in the tarball
	if fakeroot {
		hdr.Uid, hdr.Gid = 0, 0
		hdr.Uname, hdr.Gname = "root", "root"
	}
	err = tw.WriteHeader(hdr)
	if err != nil {
		ExitError(err)
	}

	//Create it if it's a directory
	if fi.IsDir() {
		tw.Flush()
		return
	}
	//Write the contents of the input file
	infile, err := os.Open(name)
	if err != nil {
		ExitError(err)
	}
	defer infile.Close()

	_, err = io.Copy(tw, infile)

	if err != nil {
		ExitError(err)
	}
	tw.Flush()
}

//create a tarball
func CreateTarball(path, arcname string, ignores []string) {
	createTarballHelper(path, arcname, false, ignores)
}

// create tarball with members' ownership set to root
func CreateFakerootTarball(path, arcname string, ignores []string) {
	createTarballHelper(path, arcname, true, ignores)
}

// creates a taball, helper for the above two functions. Skip files matching any pattern in ignores
func createTarballHelper(path, arcname string, fakeroot bool, ignores []string) {
	var buf bytes.Buffer

	path = AppendSlash(path)
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	Walk(path, func(n string, fi os.FileInfo) {
		tarFile(path, n, tw, fakeroot)
	}, ignores)

	gzfile, err := os.Create(arcname)
	if err != nil {
		ExitError(err)
	}

	defer gzfile.Close()
	gzw := gzip.NewWriter(gzfile)
	defer gzw.Close()
	_, err = io.Copy(gzw, &buf)

	if err != nil {
		ExitError(err)
	}
}
