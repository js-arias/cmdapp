// Copyright (c) 2013, J. Salvador Arias <jsalarias@csnat.unt.edu.ar>
// All rights reserved.
// Distributed under BSD-style license that can be found in the LICENSE file.
//
// This work is derived from the go tool source code
// Copyright 2011 The Go Authors.  All rights reserved.

/*
Package cmdapp implements a command line application that host a set of
commands as in the go tool or git.

In the program initialization the commands, the lists of commands (subjects)
and guides, as well as the commands (and their flags) should be setup.

If nothing else is required (there are no special cases to check before the
formal run of the application), the main funcion can be reduced to:

	a = cmdapp.App{
		// initialization
	}

	func main() {
		a.Run()
	}

And the specified command will be run.
*/
package cmdapp

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

//App is a command line application.
type App struct {
	// Name is the application's name.
	Name string

	// Synopsis is the application usage line.
	Synopsis string

	// Short is a short, single line description of the application.
	Short string

	// Long is a long description of the application.
	Long string

	// Guides is the list of help guides.
	Guides []*Guide

	// Subjects is the list of subjects (and therefore, of commands).
	Subject []*Subject
}

// Run runs the application.
func (a *App) Run() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", a.usage())
		os.Exit(2)
	}
	flag.Parse()

	args := flag.Args()
	a.RunArgs(args)
}

// RunArgs should be used if the application has its own set of flags that
// must be checked before running any command.
func (a *App) RunArgs(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%s\n", a.usage())
		os.Exit(2)
	}

	// Setup's host name
	helpCmd.Host = a.Name
	for _, s := range a.Subject {
		for _, c := range s.Commands {
			c.Host = a.Name
		}
	}

	if args[0] == "help" {
		helpCmd.Flag.Parse(args[1:])
		help(a, helpCmd.Flag.Args())
		return
	}

	for _, s := range a.Subject {
		for _, c := range s.Commands {
			if c.Name == args[0] {
				c.Flag.Usage = func() { c.usage() }
				c.Flag.Parse(args[1:])
				c.Run(c, c.Flag.Args())
				return
			}
		}
	}

	fmt.Fprintf(os.Stderr, "%s: error: unknown command %s\n", a.Name, args[0])
	os.Exit(2)
}

// returns applications help
func (a *App) help() string {
	hlp := fmt.Sprintf("%s - %s\n", a.Name, a.Short)
	hlp += fmt.Sprintf("\nSYNOPSIS\n\n    %s %s\n", a.Name, a.Synopsis)
	hlp += fmt.Sprintf("\n%s\n", strings.TrimSpace(a.Long))
	return hlp
}

// returns application usage
func (a *App) usage() string {
	usg := fmt.Sprintf("%s - %s\n", a.Name, a.Short)
	usg += fmt.Sprintf("Usage: %s %s\n", a.Name, a.Synopsis)
	for _, s := range a.Subject {
		usg += fmt.Sprintf("\n%s\n\n", s.Name)
		for _, c := range s.Commands {
			usg += fmt.Sprintf("    %-11s %s\n", c.Name, c.Short)
		}
	}
	usg += fmt.Sprintf("\nType '%s help <command>' for more information about a command.\n", a.Name)
	usg += fmt.Sprintf("Type '%s help %s' for more information about %s.\n", a.Name, a.Name, a.Name)
	usg += fmt.Sprintf("Type '%s help --guides' for a list of useful guides\n", a.Name)
	return usg
}

// runs the help command
func help(a *App, args []string) {
	if guideList {
		if len(a.Guides) == 0 {
			fmt.Fprintf(os.Stderr, "%s has no guides\n", a.Name)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s guides are:\n\n", a.Name)
		for _, g := range a.Guides {
			fmt.Fprintf(os.Stdout, "    %-11s %s\n", g.Name, g.Short)
		}
		fmt.Fprintf(os.Stdout, "\nType '%s help <guide>' for more information about a guide.\n", a.Name)
		return
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stdout, "%s", a.usage())
		return
	}

	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "%s\n", helpCmd.ErrStr("too many arguments"))
		os.Exit(2)
	}

	arg := args[0]
	if arg == a.Name {
		fmt.Fprintf(os.Stdout, "%s", a.help())
		return
	}
	if arg == "help" {
		fmt.Fprintf(os.Stdout, "%s", helpCmd.help())
		return
	}
	if arg == "documentation" {
		documentationHelp(os.Stdout, a)
		return
	}
	if arg == "doc.go" {
		f, err := os.Create("doc.go")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", helpCmd.ErrStr(err))
			os.Exit(2)
		}
		defer f.Close()
		fmt.Fprintf(f, "%s", goHead)
		documentationHelp(f, a)
		fmt.Fprintf(f, "%s", goFoot)
		return
	}
	for _, s := range a.Subject {
		for _, c := range s.Commands {
			if c.Name == arg {
				fmt.Fprintf(os.Stdout, "%s", c.help())
				return
			}
		}
	}
	for _, g := range a.Guides {
		if g.Name == arg {
			fmt.Fprintf(os.Stdout, "%s", g.help())
			return
		}
	}
	fmt.Fprintf(os.Stderr, "%s\n", helpCmd.ErrStr("unknown argument"))
	os.Exit(2)
}

func documentationHelp(w io.Writer, a *App) {
	fmt.Fprintf(w, "%s", a.help())
	for _, s := range a.Subject {
		for _, c := range s.Commands {
			fmt.Fprintf(w, "\n%s", c.help())
		}
	}
	for _, g := range a.Guides {
		fmt.Fprintf(w, "\n%s", g.help())
	}
}

var goHead = `// Authomatically generated doc.go file for use with godoc.

/*
`

var goFoot = `
*/
package main
`
