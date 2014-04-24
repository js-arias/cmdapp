// Copyright (c) 2013, J. Salvador Arias <jsalarias@csnat.unt.edu.ar>
// All rights reserved.
// Distributed under BSD-style license that can be found in the LICENSE file.
//
// This work is derived from the go tool source code
// Copyright 2011 The Go Authors.  All rights reserved.

/*
Package cmdapp implements a command line application that host a set of
commands as in the go tool or git.

In the program initialization the commands, the lists of commands and guides,
as well as the commands (and their flags) should be setup.

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
	"go/doc"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
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

	// List of commands in the the order in which they are printed by
	// help command.
	Commands []*Command
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

/*
RunArgs should be used if the application has its own set of flags that
must be checked before running any command.

The the main function should be:

	func main() {
		// Parse app owned flags
		args := flag.Parse()
		...	// use the flags, initialize, etc.

		// a is a cmdapp.App
		a.RunArgs(args)
	}
*/
func (a *App) RunArgs(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%s\n", a.usage())
		os.Exit(2)
	}

	// Setup's host name
	helpCmd.host = a.Name
	for _, c := range a.Commands {
		c.host = a.Name
		if c.Run == nil {
			c.IsCommon = false
		}
	}

	if args[0] == "help" {
		helpCmd.Flag.Usage = func() { helpCmd.Usage() }
		helpCmd.Flag.Parse(args[1:])
		help(a, helpCmd.Flag.Args())
		return
	}

	for _, c := range a.Commands {
		if c.Run == nil {
			continue
		}
		if c.Name == args[0] {
			c.Flag.Usage = func() { c.Usage() }
			c.Flag.Parse(args[1:])
			c.Run(c, c.Flag.Args())
			return
		}
	}

	fmt.Fprintf(os.Stderr, "%s: error: unknown command %s\n", a.Name, args[0])
	os.Exit(2)
}

// returns applications help
func (a *App) help() string {
	hlp := fmt.Sprintf("%s.\n", a.Short)
	hlp += fmt.Sprintf("\nSynopsis\n\n    %s %s\n", a.Name, a.Synopsis)
	hlp += fmt.Sprintf("\n%s\n", strings.TrimSpace(a.Long))
	return hlp
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

// returns application usage
func (a *App) usage() string {
	usg := fmt.Sprintf("%s.\n\n", a.Short)
	usg += fmt.Sprintf("usage: %s %s\n", a.Name, a.Synopsis)
	lr, cm := a.lenRun()
	if (lr > 15) && (cm > 0) {
		usg += fmt.Sprintf("\nThe most commonly used %s commands are:\n\n", a.Name)
		for _, c := range a.Commands {
			if !c.IsCommon {
				continue
			}
			usg += fmt.Sprintf("    %-16s %s\n", c.Name, c.Short)
		}
		usg += fmt.Sprintf("\nType '%s help --all' for a list of available commands. \n", a.Name)
		if len(a.Commands) > lr {
			usg += fmt.Sprintf("Type '%s help --guides' for a list of useful guides. \n", a.Name)
			usg += fmt.Sprintf("Type '%s help <command>' or '%s help <guide>' for more information \nabout a command or guide.\n", a.Name, a.Name)
		} else {
			usg += fmt.Sprintf("Type '%s help <command>' for more information about a command.\n", a.Name)
		}
	} else if lr > 0 {
		usg += fmt.Sprintf("\nThe commands are:\n\n")
		for _, c := range a.Commands {
			if c.Run == nil {
				continue
			}
			usg += fmt.Sprintf("    %-16s %s\n", c.Name, c.Short)
		}
		if (len(a.Commands) > lr) && (lr <= 15) {
			usg += fmt.Sprintf("\nThe common %s guides are:\n\n", a.Name)
			for _, c := range a.Commands {
				if c.Run != nil {
					continue
				}
				usg += fmt.Sprintf("    %-16s %s\n", c.Name, c.Short)
			}
		}
		if len(a.Commands) > lr {
			usg += fmt.Sprintf("\nType '%s help --guides' for a list of useful guides. \n", a.Name)
			usg += fmt.Sprintf("Type '%s help <command>' or '%s help <guide>' for more information \nabout a command or guide.\n", a.Name, a.Name)
		} else {
			usg += fmt.Sprintf("\nType '%s help <command>' for more information about a command.\n", a.Name)
		}
	} else if len(a.Commands) > 0 {
		usg += fmt.Sprintf("\nThe common %s guides are:\n\n", a.Name)
		for _, c := range a.Commands {
			usg += fmt.Sprintf("    %-16s %s\n", c.Name, c.Short)
		}
		usg += fmt.Sprintf("\nType '%s help <guide>' for more information about a guide.\n", a.Name)
	}
	return usg
}

// lenCmd returns the number of runnable and common commands
func (a *App) lenRun() (int, int) {
	lr, cm := 0, 0
	for _, c := range a.Commands {
		if c.Run != nil {
			lr++
		}
		if c.IsCommon {
			cm++
		}
	}
	return lr, cm
}

// runs the help command
func help(a *App, args []string) {
	if allList {
		if lr, _ := a.lenRun(); lr == 0 {
			fmt.Fprintf(os.Stderr, "%s has no commands\n", a.Name)
			os.Exit(1)
		}
		for _, c := range a.Commands {
			if c.Run == nil {
				continue
			}
			fmt.Fprintf(os.Stdout, "    %-16s %s\n", c.Name, c.Short)
		}
		fmt.Fprintf(os.Stdout, "\nType '%s help <command>' for more information about a command.\n", a.Name)
		return
	}
	if guideList {
		if lr, _ := a.lenRun(); len(a.Commands) == lr {
			fmt.Fprintf(os.Stderr, "%s has no guides\n", a.Name)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s guides are:\n\n", a.Name)
		for _, c := range a.Commands {
			if c.Run != nil {
				continue
			}
			fmt.Fprintf(os.Stdout, "    %-16s %s\n", c.Name, c.Short)
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
	if arg == "html" {
		f, err := os.Create(a.Name + ".html")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", helpCmd.ErrStr(err))
			os.Exit(2)
		}
		doc.ToHTML(f, a.help(), nil)
		f.Close()
		for _, c := range a.Commands {
			c.html()
		}
		return
	}
	for _, c := range a.Commands {
		if c.Name == arg {
			fmt.Fprintf(os.Stdout, "%s", c.help())
			return
		}
	}
	fmt.Fprintf(os.Stderr, "%s\n", helpCmd.ErrStr("unknown argument"))
	os.Exit(2)
}

func documentationHelp(w io.Writer, a *App) {
	fmt.Fprintf(w, "%s", a.help())
	for _, c := range a.Commands {
		fmt.Fprintf(w, "\n%s", c.help())
	}
}

var goHead = `// Authomatically generated doc.go file for use with godoc.

/*
`

var goFoot = `
*/
package main`
