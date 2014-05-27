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
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/howeyc/fsnotify"

	"github.com/bmatsuo/gutterd/handler"
	"github.com/bmatsuo/gutterd/metadata"
	"github.com/bmatsuo/gutterd/statsd"
	"github.com/bmatsuo/gutterd/watcher"
)

var (
	config   *Config            // Deamon configuration.
	handlers []*handler.Handler // The ordered set of torrent handlers.
	opt      *Options           // Command line options.
	fs       *watcher.Watcher   // Filesystem event watcher
)

func HomeDirectory() (home string, err error) {
	if home = os.Getenv("HOME"); home == "" {
		err = fmt.Errorf("HOME is not set")
	}
	return
}

// Read the config file and setup global variables.
func init() {
	opt = parseFlags()

	// Read the deamon configuration. flag overrides default (~/.config/gutterd.json)
	var err error
	defconfig := &Config{}
	configPath := opt.ConfigPath
	if configPath == "" {
		home, err := HomeDirectory()
		if err != nil {
			glog.Fatalf("unable to locate home directory: %v", err)
		}
		configPath = filepath.Join(home, ".config", "gutterd.json")
	}
	if config, err = LoadConfig(configPath, defconfig); err != nil {
		glog.Fatalf("unable to load configuration: %v", err)
	}

	if config.Statsd != "" {
		err := statsd.Init(config.Statsd, "gutterd")
		if err != nil {
			glog.Warningf("statsd init error (no stats will be recorded); %v", err)
		}
		statsd.Incr("proc.start", 1, 1)
	}

	handlers = config.MakeHandlers()

	// command line flag overrides
	if opt.Watch != nil {
		config.Watch = opt.Watch
	}
	if opt.HTTP != "" {
		config.HTTP = opt.HTTP
	}

	statsd.Incr("proc.boot", 1, 1)
}

// Handle a .torrent file.
func handleFile(path string) {
	torrent, err := metadata.ReadMetadataFile(path)
	if err != nil {
		statsd.Incr("torrent.error", 1, 1)
		glog.Errorf("error reading torrent (%q); %v", path, err)
		return
	}
	// Find the first handler matching the supplied torrent.
	for _, handler := range handlers {
		if handler.Match(torrent) {
			name := "torrent.match." + handler.Name
			statsd.Incr(name, 1, 1)
			glog.Infof("match file:%q handler:%q watch:%q",
				torrent.Info.Name,
				handler.Name,
				handler.Watch,
			)
			mvpath := filepath.Join(handler.Watch, filepath.Base(path))
			if err := os.Rename(path, mvpath); err != nil {
				glog.Error("watch import failed (%q); %v", torrent.Info.Name, err)
			}
			return
		}
	}
	statsd.Incr("torrent.no-match", 1, 1)
	glog.Warningf("no handler matched torrent: %q", torrent.Info.Name)
}

func signalHandler() {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt)
	for _ = range sig {
		fs.Close()
	}
}

func fsInit() (err error) {
	fs, err = watcher.NewInstr(
		func(event *fsnotify.FileEvent) bool {
			statsd.Incr("watcher.fs.events", 1, 1) //  filter sees all events
			return event.IsCreate() && strings.HasSuffix(event.Name, ".torrent")
		},
		func(err error) {
			statsd.Incr("watcher.fs.errors", 1, 1)
			glog.Warningf("watcher error: %v", err)
		})
	if err != nil {
		return
	}

	if err = fs.Watch(config.Watch...); err != nil {
		return
	}

	return
}

func main() {
	if err := fsInit(); err != nil {
		glog.Error("error initializing file system watcher; %v", err)
		os.Exit(1)
	}
	for event := range fs.Event {
		statsd.Incr("torrents.matches", 1, 1)
		handleFile(event.Name)
	}
}
