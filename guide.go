// Copyright (c) 2013, J. Salvador Arias <jsalarias@csnat.unt.edu.ar>
// All rights reserved.
// Distributed under BSD-style license that can be found in the LICENSE file.
//
// This work is derived from the go tool source code
// Copyright 2011 The Go Authors.  All rights reserved.

package cmdapp

import (
	"fmt"
	"strings"
)

// Guide is useful guide or concept topic that usually covers a lot of
// information about a command.
type Guide struct {
	// Name is the guide's name.
	Name string

	// Short is a short, single line description of the guide.
	Short string

	// Long is a long description of the guide.
	Long string
}

// returns guide description
func (g *Guide) help() string {
	hlp := fmt.Sprintf("%s - %s\n", capitalize(g.Name), g.Short)
	hlp += fmt.Sprintf("\n%s\n", strings.TrimSpace(g.Long))
	return hlp
}
