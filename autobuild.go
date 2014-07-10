// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/andreaskoch/go-fswatch"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	VERSION = "0.1.0"
)

var usage = func() {
	message("Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	// print application info
	message("%s (Version: %s)\n\n", os.Args[0], VERSION)

	// parse the flags
	flag.Parse()

	// get the path
	directory := Settings.Path

	// check if the supplied path exists
	if !pathExists(directory) {
		message("Path %q does not exist.", Settings.Path)
		os.Exit(1)
	}

	// clean the path
	directory = filepath.Clean(directory)

	// get the absolute path
	directory, err := filepath.Abs(directory)
	if err != nil {
		message("Cannot determine the absolute path from %q.", Settings.Path)
		os.Exit(1)
	}

	// make sure the path is a directory
	if !isDirectory(Settings.Path) {
		directory = filepath.Dir(directory)
	}

	// start the build
	stopFilesystemWatcher := make(chan bool, 1)
	build(Settings.Path, stopFilesystemWatcher)

	// stop checker
	message(`Write "stop" and press <Enter> to stop.`)

	stopApplication := make(chan bool, 1)
	go func() {
		input := bufio.NewReader(os.Stdin)

		for {

			userInput, err := input.ReadString('\n')
			if err != nil {
				fmt.Println("%s\n", err)
			}

			if command := strings.ToLower(strings.TrimSpace(userInput)); command == "stop" {

				// empty line
				message("")

				stopFilesystemWatcher <- true
				stopApplication <- true
			}
		}
	}()

	select {
	case <-stopApplication:
		debug("Stopped building %q.", Settings.Path)
	}

	os.Exit(0)
}

func build(directory string, stop chan bool) {

	recurse := true
	skipNonGoFiles := func(path string) bool {
		return !strings.HasSuffix(strings.ToLower(path), ".go")
	}

	go func() {
		folderWatcher := fswatch.NewFolderWatcher(directory, recurse, skipNonGoFiles, 2)
		folderWatcher.Start()

		for folderWatcher.IsRunning() {

			select {
			case <-folderWatcher.Modified():
				log.Printf("Buiding %q.\n", directory)

				go func() {
					execute(directory, "go install")
				}()

			case <-stop:
				debug("Stopping build for %q.", directory)
				folderWatcher.Stop()

			case <-folderWatcher.Stopped():
				break
			}
		}

		debug("Stopped build for %q.", directory)
	}()
}

func execute(directory, commandText string) {

	// get the command
	command := getCmd(directory, commandText)

	// execute the command
	if err := command.Start(); err != nil {
		fmt.Println(err)
	}

	// wait for the command to finish
	command.Wait()
}

func getCmd(directory, commandText string) *exec.Cmd {
	if commandText == "" {
		return nil
	}

	components := strings.Split(commandText, " ")

	// get the command name
	commandName := components[0]

	// get the command arguments
	arguments := make([]string, 0)
	if len(components) > 1 {
		arguments = components[1:]
	}

	// create the command
	command := exec.Command(commandName, arguments...)

	// set the working directory
	command.Dir = directory

	// redirect command io
	redirectCommandIO(command)

	return command
}

func redirectCommandIO(cmd *exec.Cmd) (*os.File, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	//direct. Masked passwords work OK!
	cmd.Stdin = os.Stdin
	return nil, err
}

func debug(text string, args ...interface{}) {
	if !Settings.Verbose {
		return
	}

	message(text, args)
}

func message(text string, args ...interface{}) {

	// append newline character
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}

	fmt.Printf(text, args...)
}
