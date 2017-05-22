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
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// help is the help command.
type help struct{}

func init() {
	Add(help{})
}

const helpCmdLong = `
Command help displays help information for a command or a help topic.

With no arguments prints to the standard output the list of available commands
and help topics.
`

func (h help) Name() string              { return "help" }
func (h help) Args() string              { return "[<command>]" }
func (h help) Short() string             { return "displays help information about " + Name }
func (h help) Long() string              { return helpCmdLong }
func (h help) Register(fs *flag.FlagSet) {}
func (h help) Runnable() bool            { return true }

func (h help) Run(args []string) error {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return nil
	}
	if len(args) != 1 {
		return errors.New("help: too many arguments.")
	}

	arg := args[0]

	// 'help documentation' generates doc.go
	if arg == "documentation" {
		f, err := os.Create("doc.go")
		if err != nil {
			return errors.Wrap(err, "help:")
		}
		defer f.Close()
		fmt.Fprintf(f, "%s\n", strings.TrimSpace(goHead))
		printUsage(f)
		mutex.Lock()
		defer mutex.Unlock()
		var cmds []string
		for _, c := range commands {
			cmds = append(cmds, c.Name())
		}
		sort.Strings(cmds)

		for _, c := range cmds {
			documentation(f, commands[c])
		}
		fmt.Fprintf(f, "\n%s", strings.TrimSpace(goFoot))
		return nil
	}

	mutex.Lock()
	c, ok := commands[arg]
	mutex.Unlock()
	if !ok {
		return errors.Errorf("help: unknown help topic: %s", arg)
	}
	documentation(os.Stdout, c)
	return nil
}

// printUsage outputs the application usage help.
func printUsage(w io.Writer) {
	fmt.Fprintf(w, "%s\n\n", Short)
	fmt.Fprintf(w, "Usage:\n\n    %s [help] <command> [<args>...]\n\n", Name)
	topics := false
	fmt.Fprintf(w, "The commands are:\n")

	mutex.Lock()
	defer mutex.Unlock()
	var cmds []string
	for _, c := range commands {
		cmds = append(cmds, c.Name())
		if !c.Runnable() {
			topics = true
		}
	}
	sort.Strings(cmds)

	for _, nm := range cmds {
		c := commands[nm]
		if !c.Runnable() {
			continue
		}
		fmt.Fprintf(w, "    %-16s %s\n", c.Name(), c.Short())
	}
	fmt.Fprintf(w, "\nUse '%s help <command>' for more information about a command.\n\n", Name)
	if !topics {
		return
	}
	fmt.Fprintf(w, "Additional help topics:\n\n")
	for _, nm := range cmds {
		c := commands[nm]
		if c.Runnable() {
			continue
		}
		fmt.Fprintf(w, "    %-16s %s\n", c.Name(), c.Short())
	}
	fmt.Fprintf(w, "\nUse '%s help <topic>' for more information about that topic.\n\n", Name)
}

var goHead = `// Authomatically generated doc.go file for use with godoc.

/*`

var goFoot = `*/
package main`
