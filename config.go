// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    config.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     2012-03-04 19:23:27.544554 -0800 PST
 *  Description: 
 */

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

type MatcherConfig struct {
	Tracker  string `json:"tracker"`  // Matched tracker urls.
	Basename string `json:"basename"` // Matched (root) file basenames.
	Ext      string `json:"ext"`      // Matched (nested-)file extensions.
}

type HandlerConfig struct {
	Name  string        `json:"name"`  // A name for logging purposes.
	Watch string        `json:"watch"` // Matching .torrent file destination.
	Match MatcherConfig `json:"match"` // Describes .torrent files to handle.
}

type LogConfig struct {
	Path    string   `json:"path"`    // Log output path (&2/&1 for stderr/stdout).
	Accepts []string `json:"accepts"` // Names logs accepted ("gutterd", "http", ...).
}

type Config struct {
	Path          string          `json:"-"`             // The path of the config file.
	HTTP          string          `json:"http"`          // HTTP service address.
	Logs          []LogConfig     `json:"logs"`          // Log configurations.
	Watch         []string        `json:"watch"`         // Incoming watch directories.
	PollFrequency int64           `json:"pollFrequency"` // Poll frequency in seconds.
	Handlers      []HandlerConfig `json:"handlers"`      // Ordered set of handlers.
}

func loadConfigFromBytes(p []byte, path string, defaults *Config) (config *Config, err error) {
	if config = new(Config); defaults != nil { // Tightly coupled events.
		*config = *defaults
	}

	if err = json.Unmarshal(p, config); err != nil {
		return config, fmt.Errorf("json error: %v", err)
	}

	config.Path = path // Overwrite any JSON specified value.

	// Validate handlers.
	for i, handler := range config.Handlers {
		name := handler.Name
		if name == "" {
			name = strconv.FormatInt(int64(i), 10)
		}
		if handler.Watch == "" {
			return config, fmt.Errorf("handler %v: no watch directory", i)
		}
		if stat, err := os.Stat(handler.Watch); err != nil {
			return config, err // The return parameter "err" is shadowed.
		} else if !stat.IsDir() {
			return config, fmt.Errorf("'watch' entry is not a directory: %s", handler.Watch)
		}
	}
	return
}

func LoadConfig(path string, defaults *Config) (*Config, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return config, err
	}
	p, err := ioutil.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("read error: %v", err)
	}
	return loadConfigFromBytes(p, path, defaults)
}

func (c *Config) MakeHandlers() []*Handler {
	handlers := make([]*Handler, 0, len(c.Handlers))
	for _, config := range c.Handlers {
		mconfig := config.Match
		m := new(matcher)
		if mconfig.Tracker != "" {
			m.Tracker = regexp.MustCompile(mconfig.Tracker)
		}
		if mconfig.Basename != "" {
			m.Basename = regexp.MustCompile(mconfig.Basename)
		}
		if mconfig.Ext != "" {
			m.Ext = regexp.MustCompile(mconfig.Ext)
		}

		handlers = append(handlers, &Handler{config.Name, config.Watch, m})
	}
	return handlers
}
