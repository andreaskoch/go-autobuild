// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"os"
)

const (
	defaultVerbose = false
)

type BuilderSettings struct {
	Path    string
	Verbose bool
}

var Settings BuilderSettings = BuilderSettings{}

func init() {

	// use the current directory as the default path
	defaultPath, err := os.Getwd()
	if err != nil {
		defaultPath = "."
	}

	flag.StringVar(&Settings.Path, "path", defaultPath, "The directory of the project that shall be build")

	flag.BoolVar(&Settings.Verbose, "verbose", defaultVerbose, "A flag for enabling verbose output")
}
