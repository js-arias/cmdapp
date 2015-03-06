// Copyright (c) 2015, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD-style license that can be found in the LICENSE file.
//
// This work is derived from the go tool source code
// Copyright 2011 The Go Authors.  All rights reserved.

// Package cmdapp implements a command line application that host a set of
// commands as in the go tool or git.
//
// In the program initialization the commands (as well as their flags), and
// the list of commands should be set up.
//
// If the main program do not have its own set of flags, just call the App
// Run's method to execute the commands, otherwise use RunArgs:
//
//	a = cmdapp.App{
//		// initialization
//	}
//
//	func main() {
//		a.Run()
//	}
package cmdapp

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// App is a command line application.
type App struct {
	// UsageLine is the usage message.
	// The first word in the line is taken to be the application name.
	UsageLine string

	// Short is a short description of the application.
	Short string

	// Commands lists the available commands and help topics. The order
	// in this list is the order in which they are printed by 'help'.
	Commands []*Command
}

// Name returns the application's name: the first word in the usage line.
func (a *App) Name() string {
	name := a.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// Run runs the application.
func (a *App) Run() {
	flag.Usage = a.usage
	flag.Parse()

	args := flag.Args()
	a.RunWithArgs(args)
}

// RunWithArgs should be used if the application has its own set of flags that
// must be checked before running any command:
//
//	func main() {
//		// parse app owned flags
//		args := flag.Parse()
//		...	// use the flags, initialize, etc.
//
//		// a is a cmdapp.App, args include all non parsed flags
//		a.RunArgs(args)
//	}
func (a *App) RunWithArgs(args []string) {
	if len(args) < 1 {
		a.usage()
	}

	// Setup's host name
	for _, c := range a.Commands {
		c.host = a.Name()
	}

	if args[0] == "help" {
		a.help(args[1:])
		return
	}

	for _, c := range a.Commands {
		if (c.Name() == args[0]) && (c.Run != nil) {
			c.Flag.Usage = func() { c.Usage() }
			c.Flag.Parse(args[1:])
			c.Run(c, c.Flag.Args())
			return
		}
	}

	fmt.Fprintf(os.Stderr, "%s: Unknown subcommand %s.\nRun '%s help' for usage.\n", a.Name(), args[0], a.Name())
	os.Exit(2)
}

func (a *App) usage() {
	a.printUsage(os.Stderr)
	os.Exit(2)
}

func (a *App) printUsage(w io.Writer) {
	fmt.Fprintf(w, "%s\n\n", a.Short)
	fmt.Fprintf(w, "Usage:\n\n    %s\n\n", a.UsageLine)
	if len(a.Commands) == 0 {
		return
	}
	fmt.Fprintf(w, "The commands are:\n")
	top := false
	for _, c := range a.Commands {
		if c.Run == nil {
			top = true
			continue
		}
		fmt.Fprintf(w, "    %-16s %s\n", c.Name(), c.Short)
	}
	fmt.Fprintf(w, "\nUse '%s help [command]' for more information about a command.\n\n")
	if !top {
		return
	}
	fmt.Fprintf(w, "Additional help topics:\n")
	for _, c := range a.Commands {
		if c.Run != nil {
			continue
		}
		fmt.Fprintf(w, "    %-16s %s\n", c.Name(), c.Short)
	}
	fmt.Fprintf(w, "\nUse '%s help [topic]' for more information about that topic.\n\n")
}

// help implements the 'help' command.
func (a *App) help(args []string) {
	if len(args) == 0 {
		a.printUsage(os.Stdout)
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "help: Too many arguments.\n\nusage: %s help <command>\n", a.Name())
		os.Exit(2)
	}

	arg := args[0]

	// 'help documentation' generates doc.go.
	if arg == "documentation" {
		f, err := os.Create("doc.go")
		if err != nil {
			fmt.Fprintf(os.Stderr, "help: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		fmt.Fprintf(f, "%s", goHead)
		a.printUsage(f)
		for _, c := range a.Commands {
			c.documentation(f)
		}
		fmt.Fprintf(f, "%s", goFoot)
		return
	}

	for _, c := range a.Commands {
		if c.Name() == arg {
			c.help(os.Stdout)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "help: Unknown help topic %s.\nRun '%s help'.\n", arg, a.Name())
	os.Exit(2)
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

var goHead = `// Authomatically generated doc.go file for use with godoc.

/*
`

var goFoot = `
*/
package main`
