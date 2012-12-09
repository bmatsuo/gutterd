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
	"errors"
	"io"
	l "log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/howeyc/fsnotify"

	"github.com/bmatsuo/gutterd/handler"
	"github.com/bmatsuo/gutterd/log"
	"github.com/bmatsuo/gutterd/metadata"
	"github.com/bmatsuo/gutterd/watcher"
)

var (
	config   *Config            // Deamon configuration.
	handlers []*handler.Handler // The ordered set of torrent handlers.
	opt      Options            // Command line options.
	fs       *watcher.Watcher   // Filesystem event watcher
)

func HomeDirectory() (home string, err error) {
	if home = os.Getenv("HOME"); home == "" {
		err = errors.New("Environment variable HOME not set.")
	}
	return
}

func logNamesFromString(s string) []string {
	accepts := strings.TrimFunc(s, unicode.IsSpace)
	namesraw := strings.Split(accepts, ",")
	names := make([]string, 0, len(namesraw))
	for i := range namesraw {
		name := strings.TrimFunc(namesraw[i], unicode.IsSpace)
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	return names
}

// Read the config file and setup global variables.
func init() {
	log.DefaultLoggerMux = new(log.LoggerMux)
	log.DefaultLogger = log.DefaultLoggerMux.NewSource("gutterd")
	initLogger := l.New(os.Stderr, "", 0)
	log.DefaultLoggerMux.NewSink(initLogger, "gutterd")

	opt = parseFlags()

	// Read the deamon configuration.
	var err error
	defconfig := &Config{
		Logs: []log.Config{
			{"&2", []string{"gutterd", "http"}},
		},
	}
	if opt.ConfigPath != "" {
		if config, err = LoadConfig(opt.ConfigPath, defconfig); err != nil {
			log.Fatal("Couldn't load configuration: ", err)
		}
	} else if home, err := HomeDirectory(); err != nil {
		log.Fatal(err)
	} else if config, err = LoadConfig(home+"/.config/gutterd.json", defconfig); err != nil {
		log.Fatal("Couldn't load configuration: ", err)
	}

	handlers = config.MakeHandlers()

	if opt.Watch != nil {
		config.Watch = opt.Watch
	}

	if opt.HTTP != "" {
		config.HTTP = opt.HTTP
	}
	if config.HTTP != "" {
		_initHTTP()
	}

	// Setup logging destinations.
	if opt.LogPath != "" {
		accepts := logNamesFromString(opt.LogAccepts)
		if len(accepts) == 0 {
			accepts = defconfig.Logs[0].Accepts
		}
		config.Logs = []log.Config{{opt.LogPath, accepts}}
	} else if accepts := logNamesFromString(opt.LogAccepts); len(accepts) > 0 {
		config.Logs = []log.Config{{defconfig.Logs[0].Path, accepts}}
	}
	for _, logConfig := range config.Logs {
		var logfile io.Writer
		switch logConfig.Path {
		case "":
			fallthrough
		case "&2":
			logfile = os.Stderr
		case "&1":
			logfile = os.Stdout
		default:
			logfile, err = os.OpenFile(logConfig.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				log.Fatalf("Couldn't open log file: %s", logConfig.Path)
			}
		}
		log.DefaultLoggerMux.NewSink(l.New(logfile, "", l.LstdFlags), logConfig.Accepts...)
	}

	log.DefaultLoggerMux.RemoveSink(initLogger)
}

// Handle a .torrent file.
func handleFile(path string) {
	torrent, err := metadata.ReadMetadataFile(path)
	if err != nil {
		log.Error(err)
		return
	}
	// Find the first handler matching the supplied torrent.
	log.Info("matching torrent to handlers ", handlers)
	for _, handler := range handlers {
		if handler.Match(torrent) {
			log.Printf("MATCH\t%s\t%s\t%s", torrent.Info.Name, handler.Name, handler.Watch)
			mvpath := filepath.Join(handler.Watch, filepath.Base(path))
			if err := os.Rename(path, mvpath); err != nil {
				log.Error(err)
			}
			return
		}
	}
	log.Print("NO MATCH\t", torrent.Info.Name)
}

func signalHandler() {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Kill, os.Interrupt)
	for _ = range sig {
		fs.Close()
	}
}

func fsInit() (err error) {
	fs, err = watcher.New(func(event *fsnotify.FileEvent) bool {
		return event.IsCreate() && strings.HasSuffix(event.Name, ".torrent")
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
	if config.HTTP != "" {
		go ListenAndServe()
	}
	if err := fsInit(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	for event := range fs.Event {
		handleFile(event.Name)
	}
}
