cmdapp
======

Package cmdapp implements a command line application that host a set of
commands as in the go tool or git.

This work is derived from the [go tool](http://golang.org/cmd/go/) source
code.


Quick usage
-----------

    go get github.com/js-arias/cmdapp

In the program initialization the commands (and their flags), as well as
the lists of commands should be set up.

If nothing else is required (there are no special cases to check before
the formal run of the application), the main funcion can be reduced to:

	import "github.com/js-arias/cmdapp"
	
	// initialization ...
	
	}
	func main() {
		Run()
	}

And the specified command will be run.

Authorship and license
----------------------

Copyright (c) 2013, J. Salvador Arias <jsalarias@csnat.unt.edu.ar>
All rights reserved.
Distributed under BSD-style license that can be found in the LICENSE file.

This work is derived from the [go tool](http://golang.org/cmd/go/) source
code. Copyright 2011 The Go Authors.  All rights reserved.
