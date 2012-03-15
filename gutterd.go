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
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	config   *Config    // Deamon configuration.
	watching []string   // Base-level directories to watch for torrents.
	handlers []*Handler // The ordered set of torrent handlers.
	opt      Options    // Command line options.
)

// Return from a pollFunc type to stop poll().
var ErrPollStop = fmt.Errorf("STOP POLLING")

// A function that can be used polling.
type pollFunc func() (time.Duration, error)

// Repeatedly call fn until ErrPollStop is returned.
func poll(fn pollFunc) {
	for {
		d, err := fn()
		if err == ErrPollStop {
			break
		} else if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		time.Sleep(d)
	}
}

func HomeDirectory() (home string, err error) {
	if home = os.Getenv("HOME"); home == "" {
		err = errors.New("Environment variable HOME not set.")
	}
	return
}

// Read the config file and setup global variables.
func init() {
	var err error
	defconfig := &Config{
		PollFrequency: 60,
		LogPath:       "&2",
	}

	loggerMux = new(LoggerMux)
	DefaultLogger = loggerMux.NewSource("gutterd")
	initLogger := log.New(os.Stderr, "", 0)
	loggerMux.NewSink(initLogger, "gutterd")

	opt = parseFlags()

	// Read the deamon configuration.
	if opt.ConfigPath != "" {
		if config, err = LoadConfig(opt.ConfigPath, defconfig); err != nil {
			Fatalf("%s: %v", "Couldn't load configuration", err)
		}
	} else if home, err := HomeDirectory(); err != nil {
		Fatal(err)
	} else if config, err = LoadConfig(home+"/.config/gutterd.json", defconfig); err != nil {
		Fatalf("%s: %v", "Couldn't load configuration", err)
	}

	handlers = config.MakeHandlers()

	watching = config.Watch

	if opt.PollFrequency > 0 {
		config.PollFrequency = opt.PollFrequency
	}

	if opt.Watch != nil {
		config.Watch = opt.Watch
	}

	// Setup the logging destination.
	if opt.LogPath != "" {
		config.LogPath = opt.LogPath
	}
	var logfile io.Writer
	switch config.LogPath {
	case "":
		fallthrough
	case "&2":
		logfile = os.Stderr
	case "&1":
		logfile = os.Stdout
	default:
		logfile, err = os.OpenFile(config.LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			Fatalf("Couldn't open log file: %s", config.LogPath)
		}
	}
	loggerMux.NewSink(log.New(logfile, "", log.LstdFlags), "gutterd")
	loggerMux.RemoveSink(initLogger)
}

// Handle a .torrent file.
func handleFile(path string) {
	torrent, err := ReadMetadataFile(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	Printf("%-8s%s", "TORRENT", torrent.Info.Name)
	// Find the first handler matching the supplied torrent.
	for _, handler := range handlers {
		if handler.Match(torrent) {
			Printf("%-8s%-14s", "MATCH", handler.Name)
			Printf("%-8s%s\n\n", "MOVING", handler.Watch)
			mvpath := filepath.Join(handler.Watch, filepath.Base(path))
			if err := os.Rename(path, mvpath); err != nil {
				Printf("%-8s%v", "ERROR", err)
			}
			return
		}
	}
	Printf("%-8s%-14s%s\n\n", "NO MATCH", "", torrent.Info.Name)
}

func main() {
	if len(watching) == 0 {
		Fatalf("%s\t%s", "ERROR", "Not watching any directories")
	}

	// Poll watch directories, handling all torrents found.
	poll(pollFunc(func() (time.Duration, error) {
		for _, watch := range watching {
			torrents, err := filepath.Glob(filepath.Join(watch, "*.torrent"))
			if err != nil {
				Printf("Error polling %s:\n%v", watch, err)
				continue
			}
			for _, _torrent := range torrents {
				handleFile(_torrent)
				continue
			}
		}
		return (config.PollFrequency * 1e9), nil
	}))
}
