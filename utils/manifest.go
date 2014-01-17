//
// Copyright (c) 2014 Canonical Ltd.
//
// Author: Jani Monoses <jani@ubuntu.com>
//
package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
)

type Manifest map[string]interface{}

func ReadManifest(manifestPath string) Manifest {
	var v interface{}

	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		ExitError(err)
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		ExitError(err)
	}

	m := v.(map[string]interface{})

	if m["architecture"] == nil {
		m["architecture"] = "all"
	}

	return m
}

func EpochlessVersion(v string) string {
	re, err := regexp.Compile("^[0-9]+:")
	if err != nil {
		ExitError(err)
	}
	return re.ReplaceAllString(v, "")
}

func WriteManifest(manifestPath string, m Manifest) {
	b, err := json.MarshalIndent(m, "", "  ")

	mf, err := os.Create(manifestPath)
	if err != nil {
		ExitError(err)
	}

	mf.Write(b)
	mf.Write([]byte("\n"))

	mf.Chmod(0644)

	mf.Close()
}
