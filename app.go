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
	"io"
	"os"
	"strings"
)

// Short is a short description of the application.
var Short string

// Commands is the list of available commands and help topics. The order in
// the list is used for help output.
var Commands []*Command

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
	if args[0] == "help" {
		help(args[1:])
		return
	}

	for _, c := range Commands {
		if (c.Name() == args[0]) && (c.Run != nil) {
			c.Flag.Usage = func() { c.Usage() }
			c.Flag.Parse(args[1:])
			err := c.Run(c, c.Flag.Args())
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s: %v\n", Name, c.Name(), err)
				os.Exit(1)
			}
			return
		}
	}

	fmt.Fprintf(os.Stderr, "%s: unknown subcommand %s\nRun '%s help' for usage.\n", Name, args[0], Name)
	os.Exit(1)
}

// usage printd application's help and exists.
func usage() {
	printUsage(os.Stderr)
	os.Exit(1)
}

// printUsage outputs the application usage help.
func printUsage(w io.Writer) {
	fmt.Fprintf(w, "%s\n\n", Short)
	fmt.Fprintf(w, "Usage:\n\n    %s [help] <command> [<args>...]\n\n", Name)
	topics := false
	fmt.Fprintf(w, "The commands are:\n")
	for _, c := range Commands {
		if c.Run == nil {
			topics = true
			continue
		}
		fmt.Fprintf(w, "    %-16s %s\n", c.Name(), c.Short)
	}
	fmt.Fprintf(w, "\nUse '%s help <command>' for more information about a command.\n\n", Name)
	if !topics {
		return
	}
	fmt.Fprintf(w, "Additional help topics:\n\n")
	for _, c := range Commands {
		if c.Run != nil {
			continue
		}
		fmt.Fprintf(w, "    %-16s %s\n", c.Name(), c.Short)
	}
	fmt.Fprintf(w, "\nUse '%s help <topic>' for more information about that topic.\n\n", Name)
}

// help implements the 'help' command.
func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "%s: help: too many arguments.\nUsage: '%s help [<command>]'\n", Name, Name)
		os.Exit(1)
	}

	arg := args[0]

	// 'help documentation' generates doc.go
	if arg == "documentation" {
		f, err := os.Create("doc.go")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: help: %v\n", Name, err)
			os.Exit(1)
		}
		defer f.Close()
		fmt.Fprintf(f, "%s\n", strings.TrimSpace(goHead))
		printUsage(f)
		for _, c := range Commands {
			c.documentation(f)
		}
		fmt.Fprintf(f, "\n%s", strings.TrimSpace(goFoot))
		return
	}

	for _, c := range Commands {
		if c.Name() == arg {
			c.help()
			return
		}
	}

	fmt.Fprintf(os.Stderr, "%s: help: unknown help topic %s.\nRun '% help'\n", Name, arg, Name)
	os.Exit(1)
}

var goHead = `// Authomatically generated doc.go file for use with godoc.

/*`

var goFoot = `*/
package main`
