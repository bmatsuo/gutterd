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
	"time"
)

type MatcherConfiguration struct {
	Tracker  string // Regexp pattern matching tracker urls.
	Basename string // Regexp pattern matching the (root) file's basename.
	Ext      string // Regexp pattern matching file extensions.
}

type HandlerConfiguration struct {
	Name  string // The log name for the handler.
	Watch string // The directory handled torrents are placed in.
	Match MatcherConfiguration
}

// The configuration is stored as JSON. Attributes are camelCase.
type Configuration struct {
	LogPath       string                 // File path (if any) to direct the log to.
	Watch         []string               // Directories watched for new torrents to handle.
	PollFrequency time.Duration          // Delay (seconds) between polling watch directories.
	Handlers      []HandlerConfiguration // A prioritized list of torrent handlers.
}

func LoadConfig(path string, defaults *Configuration) (*Configuration, error) {
	config := new(Configuration)
	if defaults != nil {
		*config = *defaults
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return config, err
	}
	if configBytes, err := ioutil.ReadFile(path); err != nil {
		return config, fmt.Errorf("read error: %v", err)
	} else if err = json.Unmarshal(configBytes, config); err != nil {
		return config, fmt.Errorf("json error: %v", err)
	}
	for i, handler := range config.Handlers {
		name := handler.Name
		if name == "" {
			name = strconv.FormatInt(int64(i), 10)
		}
		if handler.Watch == "" {
			return config, fmt.Errorf("handler %v: no watch directory", i)
		}
		if stat, err := os.Stat(handler.Watch); err != nil {
			return config, fmt.Errorf("can't stat watch directory %s: %v", handler.Watch, err)
		} else if !stat.IsDir() {
			return config, fmt.Errorf("watch entry is not a directory: %s", handler.Watch)
		}
	}
	return config, nil
}

func (c *Configuration) MakeHandlers() []*Handler {
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
