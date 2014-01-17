//
// Copyright (c) 2014 Canonical Ltd.
//
// Author: Jani Monoses <jani@ubuntu.com>
//
package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

const magicHeader = "!<arch>\n"

type arWriter struct {
	f *os.File
	w io.Writer
}

func NewArWriter(path string) (*arWriter, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &arWriter{f, f}, nil
}

func (w *arWriter) AddMagic() {
	w.w.Write([]byte(magicHeader))
}

func (w *arWriter) AddHeader(name string, size int) error {
	if len(name) > 15 {
		return errors.New("ar member name length too long")
	}
	if size > 9999999999 {
		return errors.New("ar member size too large")
	}
	header := fmt.Sprintf("%-16s%-12d0     0     100644  %-10d`\n",
		name, time.Now().Unix(), size)
	if len(header) != 60 {
		return errors.New("Header len not 60")
	}
	w.w.Write([]byte(header))
	return nil
}

func (w *arWriter) AddData(name string, data []byte) error {
	size := len(data)
	w.AddHeader(name, size)
	_, err := w.w.Write(data)
	if err != nil {
		return err
	}
	if size&1 == 1 {
		w.w.Write([]byte("\n"))
	}
	return nil
}

func (w *arWriter) AddFile(name, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	size := fi.Size()
	w.AddHeader(name, int(size))
	_, err = io.Copy(w.w, f)
	if err != nil {
		return err
	}

	if size&1 == 1 {
		w.w.Write([]byte("\n"))
	}
	return nil
}

func (w *arWriter) Close() {
	w.f.Close()
}
