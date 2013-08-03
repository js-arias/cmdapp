// Copyright (c) 2013, J. Salvador Arias <jsalarias@csnat.unt.edu.ar>
// All rights reserved.
// Distributed under BSD-style license that can be found in the LICENSE file.
//
// This work is derived from the go tool source code
// Copyright 2011 The Go Authors.  All rights reserved.

package cmdapp

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// A command is a hosted subcommand.
type Command struct {
	// Run runs the command.
	// The argument list is the set of unparsed arguments, that is the
	// arguments unparsed by the flag package.
	Run func(c *Command, args []string)

	// Name is the command's name.
	Name string

	// Host is the name of the application that hosts the command.
	Host string

	// Synopsis is the command usage line.
	Synopsis string

	// Short is a short, single line description of the command.
	Short string

	// Long is a long description of the command.
	Long string

	// Set of flags specific to the command.
	Flag flag.FlagSet
}

// ErrStr returns an error description from the command
func (c *Command) ErrStr(err interface{}) string {
	return fmt.Sprintf("%s %s: error: %v", c.Host, c.Name, err)
}

// prints command usage
func (c *Command) usage() {
	fmt.Fprintf(os.Stderr, "%s-%s - %s\n", c.Host, c.Name, c.Short)
	fmt.Fprintf(os.Stderr, "Usage: %s %s %s\n", c.Host, c.Name, c.Synopsis)
	fmt.Fprintf(os.Stderr, "Type '%s help %s' for more information\n", c.Host, c.Name)
	os.Exit(2)
}

// returns command help
func (c *Command) help() string {
	hlp := fmt.Sprintf("%s-%s - %s\n", c.Host, c.Name, c.Short)
	hlp += fmt.Sprintf("\nSYNOPSIS\n\n    %s %s %s\n", c.Host, c.Name, c.Synopsis)
	hlp += fmt.Sprintf("\n%s\n", strings.TrimSpace(c.Long))
	return hlp
}

// help command, a dummy non-runnable command
var helpCmd = Command{
	Name:     "help",
	Synopsis: "[-g|--guide] [<command>|<guide>]",
	Short:    "Display help information",
	Long: `
DESCRIPTION

With no option and no COMMAND or GUIDE given, the list of commands are printed
to the standard output.

If the option --guide is given, a list of useful guides is also printed on the
standard output.

If a command, or a guide, is given, the information for that command or guide
is printed in the standard output.

OPTIONS

    -g
    --guides
      Prints a list of useful guides on the standard output. This option
      overrides any given command or guide name.
	`,
}

var guideList = false

func init() {
	helpCmd.Flag.BoolVar(&guideList, "guides", false, "")
	helpCmd.Flag.BoolVar(&guideList, "g", false, "")
}