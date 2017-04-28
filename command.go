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
	"unicode"
	"unicode/utf8"
)

// A Command is a hosted command.
type Command interface {
	// Run runs the command.
	// It returns an error if the command finish on error.
	// The argument list is the set of arguments unparsed by flag package.
	Run(args []string) error

	// Name returns the command's name.
	Name() string

	// Args is the command's argument list.
	Args() string

	// Short is a short description of the command.
	Short() string

	// Long is the long message shown in 'help <this-command>' output.
	Long() string

	// Register sets command-specific flags.
	Register(*flag.FlagSet)

	// Runnable reports whether the command can be run;
	// otherwise it is a documentation pseudo-command.
	Runnable() bool
}

// Usage prints the usage message and exits the program.
func Usage(c Command) {
	fmt.Fprintf(os.Stderr, "usage: %s %s %s\n\n", Name, c.Name(), c.Args())
	fmt.Fprintf(os.Stderr, "Type '%s help %s' for more information.\n", Name, c.Name())
	os.Exit(1)
}

// documentation prints command documentation.
func documentation(w io.Writer, c Command) {
	fmt.Fprintf(w, "%s\n\n", capitalize(c.Short()))
	if c.Runnable() {
		fmt.Fprintf(w, "Usage:\n\n    %s %s %s\n\n", Name, c.Name(), c.Args())
	}
	fmt.Fprintf(w, "%s\n\n", strings.TrimSpace(c.Long()))
}

// capitalize set the first rune of a string as upper case.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}
