// Copyright (c) 2015, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD-style license that can be found in the LICENSE file.
//
// This work is derived from the go tool source code
// Copyright 2011 The Go Authors.  All rights reserved.

package cmdapp

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// A Command is an implementation of a hosted command.
type Command struct {
	// Run runs the command.
	// The argument list is the set of unparsed arguments, that is the
	// arguments unparsed by the flag package.
	Run func(c *Command, args []string)

	// UsageLine is the ussage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is a short, single line description of the command.
	Short string

	// Long is a long description of the command.
	Long string

	// Set of flags specific to the command.
	Flag flag.FlagSet

	// Host is the name of the application that hosts the command.
	host string
}

// Name returns the commands's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// Usage prints the usage help of the command.
func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s %s\n\n", c.host, c.UsageLine)
	fmt.Fprintf(os.Stderr, "Type '%s help %s for more information.\n", c.host, c.Name())
	os.Exit(2)
}

func (c *Command) documentation(w io.Writer) {
	fmt.Fprintf(w, "%s\n\n", capitalize(c.Short))
	if c.Run != nil {
		fmt.Fprintf(w, "Usage:\n\n    %s %s\n\n", c.host, c.UsageLine)
	}
	fmt.Fprintf(w, "%s\n\n", strings.TrimSpace(c.Long))
}

func (c *Command) help(w io.Writer) {
	if c.Run != nil {
		fmt.Fprintf(w, "usage: %s %s\n\n", c.host, c.UsageLine)
	}
	fmt.Fprintf(w, "%s\n", strings.TrimSpace(c.Long))
}
