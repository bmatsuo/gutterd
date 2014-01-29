// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    options.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     2012-03-04 17:28:31.729424 -0800 PST
 */

import (
	"flag"
	"path/filepath"

	"github.com/bmatsuo/gutterd/watcher"
)

// A struct that holds parsed command line flags.
type Options struct {
	HTTP          string
	ConfigPath    string
	PollFrequency int64
	watchStr      string
	Watch         []watcher.Config
	LogPath       string
	LogAccepts    string
}

// attach command line flags to opt. call flag.Parse() after.
func setupFlags(opt *Options) {
	flag.Int64Var((*int64)(&opt.PollFrequency), "poll", 0, "Specify a polling frequency (in seconds).")
	flag.StringVar(&opt.HTTP, "http", "", "Address to serve web requests from (e.g. ':6060').")
	flag.StringVar(&opt.watchStr, "watch", "", "Specify a set of directories to watch.")
	flag.StringVar(&opt.ConfigPath, "config", "", "A config file to use instead of ~/.config/gutterd.json.")
}

// check flags for acceptable values.
func verifyFlags(opt *Options) error {
	if opt.watchStr != "" {
		for _, dir := range filepath.SplitList(opt.watchStr) {
			opt.Watch = append(opt.Watch, watcher.Config(dir))
		}
		for _, w := range opt.Watch {
			if err := w.Validate(); err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

// parse flags and validate them.
func parseFlags() *Options {
	opt := new(Options)
	setupFlags(opt)
	flag.Parse()
	verifyFlags(opt)
	return opt
}
