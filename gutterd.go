// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    gutterd.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     2012-03-04 17:28:31.728667 -0800 PST
 *  Description: Main source file in gutterd
 */

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"code.google.com/p/go.exp/fsnotify"
	"github.com/bmatsuo/gutterd/handler"
	"github.com/bmatsuo/gutterd/statsd"
	"github.com/bmatsuo/gutterd/watcher"
	"github.com/bmatsuo/torrent/bencoding"
	"github.com/bmatsuo/torrent/metainfo"
	"github.com/golang/glog"
)

func HomeDirectory() (home string, err error) {
	if home = os.Getenv("HOME"); home == "" {
		err = fmt.Errorf("HOME is not set")
	}
	return
}

func main() {
	// A struct that holds parsed command line flags.
	var opts struct {
		ConfigPath    string
		PollFrequency int64
		watchStr      string
		Watch         []watcher.Config
		LogPath       string
		LogAccepts    string
	}

	// attach command line flags to opt. call flag.Parse() after.
	flag.Int64Var((*int64)(&opts.PollFrequency), "poll", 0, "Specify a polling frequency (in seconds).")
	flag.StringVar(&opts.watchStr, "watch", "", "Specify a set of directories to watch.")
	flag.StringVar(&opts.ConfigPath, "config", "", "A config file to use instead of ~/.config/gutterd.json.")
	flag.Parse()

	// check flags for acceptable values.
	if opts.watchStr != "" {
		for _, dir := range filepath.SplitList(opts.watchStr) {
			opts.Watch = append(opts.Watch, watcher.Config(dir))
		}
		for _, w := range opts.Watch {
			if err := w.Validate(); err != nil {
				glog.Fatal("watch: %v", err)
			}
			continue
		}
	}

	// Read the deamon configuration.
	var err error
	defconfig := &Config{}
	if opts.ConfigPath == "" {
		home, err := HomeDirectory()
		if err != nil {
			glog.Fatalf("unable to locate configuration: no home directory -- %v", err)
		}
		opts.ConfigPath = filepath.Join(home, ".config", "gutterd.json")
	}
	config, err := LoadConfig(opts.ConfigPath, defconfig)
	if err != nil {
		glog.Fatalf("unable to load configuration: %v", err)
	}

	if config.Statsd != "" {
		err := statsd.Init(config.Statsd, "gutterd")
		if err != nil {
			glog.Warningf("statsd init error (no stats will be recorded); %v", err)
		}
		statsd.Incr("proc.start", 1, 1)
	}

	handlers := config.MakeHandlers()

	// command line flag overrides
	if opts.Watch != nil {
		config.Watch = opts.Watch
	}

	statsd.Incr("proc.boot", 1, 1)

	sig := make(chan os.Signal, 2)
	kill := make(chan struct{})
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		signal.Stop(sig)
		close(kill)
	}()

	fs, err := SetupWatcher(kill)
	if err != nil {
		glog.Fatalf("error creating up file system watcher; %v", err)
	}

	if err = fs.Watch(config.Watch...); err != nil {
		glog.Fatalf("error initializing file system watcher; %v", err)
		return
	}

	for event := range fs.Event {
		statsd.Incr("torrents.matches", 1, 1)
		HandlePath(handlers, event.Name)
	}
}

func SetupWatcher(stop <-chan struct{}) (*watcher.Watcher, error) {
	fs, err := watcher.NewInstr(
		func(event *fsnotify.FileEvent) bool {
			statsd.Incr("watcher.fs.events", 1, 1) //  filter sees all events
			return event.IsCreate() && strings.HasSuffix(event.Name, ".torrent")
		},
		func(err error) {
			statsd.Incr("watcher.fs.errors", 1, 1)
			glog.Warningf("watcher error: %v", err)
		})
	if err != nil {
		return nil, err
	}
	go func() {
		<-stop
		glog.Infof("shutting down watchers")
		err := fs.Close()
		if err != nil {
			glog.Warningf("error closing filesystem watcher")
		}
	}()
	return fs, err
}

// Handle a .torrent file.
func HandlePath(handlers []*handler.Handler, path string) {
	p, err := ioutil.ReadFile(path)
	if err != nil {
		statsd.Incr("torrent.error", 1, 1)
		glog.Errorf("error reading torrent (%q); %v", path, err)
		return
	}
	var torrent *metainfo.Metainfo
	err = bencoding.Unmarshal(&torrent, p)
	if err != nil {
		statsd.Incr("torrent.error", 1, 1)
		glog.Errorf("error reading torrent (%q); %v", path, err)
		return
	}

	// Find the first handler matching the supplied torrent.
	for _, h := range handlers {
		err := h.Handle(path, torrent)
		if err == handler.ErrNoMatch {
			continue
		}
		if err != nil {
			glog.Warning(err)
			continue
		}
		name := "torrent.match." + h.Name
		statsd.Incr(name, 1, 1)
		glog.Infof("%q matched file %q", torrent.Info.Name, h.Name)
		return
	}

	statsd.Incr("torrent.no-match", 1, 1)
	glog.Warningf("no handler matched torrent: %q", torrent.Info.Name)
}
