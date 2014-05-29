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
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"code.google.com/p/go.exp/fsnotify"
	"github.com/bmatsuo/torrent/bencoding"
	"github.com/bmatsuo/torrent/metainfo"
	"github.com/golang/glog"

	"github.com/bmatsuo/gutterd/handler"
	"github.com/bmatsuo/gutterd/statsd"
	"github.com/bmatsuo/gutterd/watcher"
)

var (
	config *Config  // Deamon configuration.
	opt    *Options // Command line options.
)

func HomeDirectory() (home string, err error) {
	if home = os.Getenv("HOME"); home == "" {
		err = fmt.Errorf("HOME is not set")
	}
	return
}

// Handle a .torrent file.
func handle(handlers []*handler.Handler, path string) {
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

func main() {
	opt = parseFlags()

	// Read the deamon configuration. flag overrides default (~/.config/gutterd.json)
	var err error
	defconfig := &Config{}
	if opt.ConfigPath == "" {
		home, err := HomeDirectory()
		if err != nil {
			glog.Fatalf("unable to locate home directory: %v", err)
		}
		opt.ConfigPath = filepath.Join(home, ".config", "gutterd.json")
	}
	if config, err = LoadConfig(opt.ConfigPath, defconfig); err != nil {
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
	if opt.Watch != nil {
		config.Watch = opt.Watch
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
		glog.Fatalf("error creating file system watcher: %v", err)
	}
	go func() {
		<-kill
		glog.Infof("shutting down watchers")
		err := fs.Close()
		if err != nil {
			glog.Warningf("error closing filesystem watcher")
		}
	}()

	if err = fs.Watch(config.Watch...); err != nil {
		glog.Fatalf("error initializing file system watcher; %v", err)
		return
	}

	for event := range fs.Event {
		statsd.Incr("torrents.matches", 1, 1)
		handle(handlers, event.Name)
	}
}
