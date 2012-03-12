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
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	config   *Config     // Deamon configuration.
	watching []string    // Base-level directories to watch for torrents.
	handlers []*Handler  // The ordered set of torrent handlers.
	logger   *log.Logger // The global logger.
	opt      Options     // Command line options.
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

// Read the config file and setup global variables.
func init() {
	var err error
	defconfig := &Config{
		PollFrequency: 60,
		LogPath:       "&2",
	}
	// Read the deamon configuration.
	config, err = LoadConfig("/Users/b/.config/gutterd.json", defconfig)
	if err != nil {
		fmt.Printf("%-8s%s: %v", "ERROR", "Couldn't load configuration", err)
		os.Exit(1)
	}
	opt = parseFlags()
	if opt.LogPath != "" {
		config.LogPath = opt.LogPath
	}
	if opt.PollFrequency > 0 {
		config.PollFrequency = opt.PollFrequency
	}

	// Setup the logging destination.
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
			fmt.Fprintf(os.Stderr, "Couldn't open log file %s", config.LogPath)
			os.Exit(0)
		}
	}
	logger = log.New(logfile, "gutterd ", log.LstdFlags)

	handlers = config.MakeHandlers()

	watching = config.Watch
}

// Handle a .torrent file.
func handleFile(path string) {
	torrent, err := ReadMetadataFile(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	logger.Printf("%-8s%s", "TORRENT", torrent.Info.Name)
	// Find the first handler matching the supplied torrent.
	for _, handler := range handlers {
		if handler.Match(torrent) {
			logger.Printf("%-8s%-14s", "MATCH", handler.Name)
			logger.Printf("%-8s%s\n\n", "MOVING", handler.Watch)
			mvpath := filepath.Join(handler.Watch, filepath.Base(path))
			if err := os.Rename(path, mvpath); err != nil {
				logger.Printf("%-8s%v", "ERROR", err)
			}
			return
		}
	}
	logger.Printf("%-8s%-14s%s\n\n", "NO MATCH", "", torrent.Info.Name)
}

func main() {
	if len(watching) == 0 {
		logger.Printf("%-8s%s", "ERROR", "Not watching any directories")
		os.Exit(1)
	}

	// Poll watch directories, handling all torrents found.
	poll(pollFunc(func() (time.Duration, error) {
		for _, watch := range watching {
			torrents, err := filepath.Glob(filepath.Join(watch, "*.torrent"))
			if err != nil {
				logger.Printf("Error polling %s:\n%v", watch, err)
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
