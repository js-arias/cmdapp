// Copyright (c) 2015, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD-style license that can be found in the LICENSE file.
//
// This work is derived from the go tool source code
// Copyright 2011 The Go Authors.  All rights reserved.

// Package cmdapp implements a command line application that host a set of
// commands as in the go tool and git.
//
// During program initialization the commands (as well as their flags), and
// the list of commands should be set up.
//
// In most simple case, the Run function will execute the required command:
//	import "github.com/js-arias/cmdapp"
//
//	// initialize commands...
//
//	func main() {
//		cmdapp.Run()
//	}
package cmdapp

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Short is a short description of the application.
var Short string

// commands is the list of available commands and help topics.
var (
	mutex    sync.Mutex
	commands = make(map[string]Command)
)

// Add adds a new command to the application.
// Command names should be unique,
// otherwise it will trigger a panic.
func Add(c Command) {
	name := strings.ToLower(c.Name())
	mutex.Lock()
	defer mutex.Unlock()
	if _, dup := commands[name]; dup {
		msg := fmt.Sprintf("cmdapp: Repeated command name: %s %s", name, c.Short())
		panic(msg)
	}
	commands[name] = c
}

// Name stores the application name, the default is based on the arguments of
// the program.
var Name = os.Args[0]

// Run runs the application.
func Run() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	mutex.Lock()
	c, ok := commands[args[0]]
	mutex.Unlock()
	if !ok || !c.Runnable() {
		fmt.Fprintf(os.Stderr, "%s: unknown subcommand %s\nRun '%s help' for usage.\n", Name, args[0], Name)
		os.Exit(1)
	}

	fs := flag.NewFlagSet(c.Name(), flag.ExitOnError)
	fs.Usage = func() { Usage(c) }
	c.Register(fs)
	fs.Parse(args[1:])
	err := c.Run(fs.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s: %v\n", Name, c.Name(), err)
		os.Exit(1)
	}
}

// usage printd application's help and exists.
func usage() {
	printUsage(os.Stderr)
	os.Exit(1)
}
