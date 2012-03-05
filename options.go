// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    options.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     2012-03-04 17:28:31.729424 -0800 PST
 *  Description: Option parsing for levyd
 */

import (
    "flag"
    "fmt"
    "os"
)

// TODO Customize exported (capitalized) variables, types, and functions.

var (
    CmdHelpUsage string // Custom usage string.
    CmdHelpFoot  string // Printed after help.
)

// A struct that holds levyd's parsed command line flags.
type Options struct {
    Verbose bool
}

//  Create a flag.FlagSet to parse the levyd's flags.
func SetupFlags(opt *Options) *flag.FlagSet {
    fs := flag.NewFlagSet("levyd", flag.ExitOnError)
    fs.BoolVar(&opt.Verbose, "v", false, "Verbose program output.")
    return setupUsage(fs)
}

// Check the levyd's flags and arguments for acceptable values.
// When an error is encountered, panic, exit with a non-zero status, or override
// the error.
func VerifyFlags(opt *Options, fs *flag.FlagSet) {
}

/**************************/
/* Do not edit below here */
/**************************/

//  Print a help message to standard error. See CmdHelpUsage and CmdHelpFoot.
func PrintHelp() { SetupFlags(&Options{}).Usage() }

//  Hook up CmdHelpUsage and CmdHelpFoot with flag defaults to function flag.Usage.
func setupUsage(fs *flag.FlagSet) *flag.FlagSet {
    printNonEmpty := func (s string) {
        if s != "" {
            fmt.Fprintf(os.Stderr, "%s\n", s)
        }
    }
    fs.Usage = func() {
        printNonEmpty(CmdHelpUsage)
        fs.PrintDefaults()
        printNonEmpty(CmdHelpFoot)
    }
    return fs
}

//  Parse the flags, validate them, and post-process (e.g. Initialize more complex structs).
func parseFlags() Options {
    var opt Options
    fs := SetupFlags(&opt)
    fs.Parse(os.Args[1:])
    VerifyFlags(&opt, fs)
    // Process the verified Options...
    return opt
}
