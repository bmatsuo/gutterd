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
	"path/filepath"

	"github.com/bmatsuo/gutterd/watcher"
)

// TODO Customize exported (capitalized) variables, types, and functions.

var (
	cmdHelpUsage = "gutterd [options]"
	cmdHelpFoot  string
)

// A struct that holds levyd's parsed command line flags.
type Options struct {
	HTTP          string
	ConfigPath    string
	PollFrequency int64
	watchStr      string
	Watch         []watcher.Config
	LogPath       string
	LogAccepts    string
}

//  Create a flag.FlagSet to parse the levyd's flags.
func setupFlags(opt *Options) *flag.FlagSet {
	fs := flag.NewFlagSet("levyd", flag.ExitOnError)
	fs.Int64Var((*int64)(&opt.PollFrequency), "poll", 0, "Specify a polling frequency (in seconds).")
	fs.StringVar(&opt.HTTP, "http", "", "Address to serve web requests from (e.g. ':6060').")
	fs.StringVar(&opt.watchStr, "watch", "", "Specify a set of directories to watch.")
	fs.StringVar(&opt.ConfigPath, "config", "", "A config file to use instead of ~/.config/gutterd.json.")
	fs.StringVar(&opt.LogPath, "log", "", "Log output path.")
	fs.StringVar(&opt.LogAccepts, "log-accepts", "", "Comma separated list of logs (e.g. 'gutterd,http').")
	return setupUsage(fs)
}

// Check the levyd's flags and arguments for acceptable values.
// When an error is encountered, panic, exit with a non-zero status, or override
// the error.
func verifyFlags(opt *Options, fs *flag.FlagSet) {
	if opt.watchStr != "" {
		for _, dir := range filepath.SplitList(opt.watchStr) {
			opt.Watch = append(opt.Watch, watcher.Config(dir))
		}
		for _, w := range opt.Watch {
			if err := w.Validate(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			continue
		}
	}
}

/**************************/
/* Do not edit below here */
/**************************/

//  Print a help message to standard error. See cmdHelpUsage and cmdHelpFoot.
func printHelp() { setupFlags(&Options{}).Usage() }

//  Hook up cmdHelpUsage and cmdHelpFoot with flag defaults to function flag.Usage.
func setupUsage(fs *flag.FlagSet) *flag.FlagSet {
	printNonEmpty := func(s string) {
		if s != "" {
			fmt.Fprintf(os.Stderr, "%s\n", s)
		}
	}
	fs.Usage = func() {
		printNonEmpty(cmdHelpUsage)
		fs.PrintDefaults()
		printNonEmpty(cmdHelpFoot)
	}
	return fs
}

//  Parse the flags, validate them, and post-process (e.g. Initialize more complex structs).
func parseFlags() Options {
	var opt Options
	fs := setupFlags(&opt)
	fs.Parse(os.Args[1:])
	verifyFlags(&opt, fs)
	// Process the verified Options...
	return opt
}
