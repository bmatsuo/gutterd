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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/bmatsuo/gutterd/handler"
	"github.com/bmatsuo/gutterd/log"
	"github.com/bmatsuo/gutterd/watcher"
)

type Config struct {
	Path          string           `json:"-"`             // The path of the config file.
	HTTP          string           `json:"http"`          // HTTP service address.
	Logs          []log.Config     `json:"logs"`          // Log configurations.
	Watch         []watcher.Config `json:"watch"`         // Incoming watch directories.
	PollFrequency int64            `json:"pollFrequency"` // Poll frequency in seconds.
	Handlers      []handler.Config `json:"handlers"`      // Ordered set of handlers.
}

func (config Config) Validate() error {
	if config.Path == "" {
		return errors.New("config: no path")
	}
	for _, watcher := range config.Watch {
		if err := watcher.Validate(); err != nil {
			return fmt.Errorf("config: %v", err)
		}
	}
	if config.PollFrequency <= 0 {
		return fmt.Errorf("config: invalid pollFrequency: %d", config.PollFrequency)
	}
	for _, handler := range config.Handlers {
		if err := handler.Validate(); err != nil {
			return fmt.Errorf("config: %v", err)
		}
	}
	return nil
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
	config, err := loadConfigFromBytes(p, path, defaults)
	if err != nil {
		return nil, err
	}
	err = config.Validate()
	return config, err
}

func (c *Config) MakeHandlers() []*handler.Handler {
	handlers := make([]*handler.Handler, len(c.Handlers))
	for i := range c.Handlers {
		handlers[i] = c.Handlers[i].Handler()
	}
	return handlers
}
