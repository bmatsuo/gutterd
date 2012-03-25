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
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var whitespace = regexp.MustCompile(`.\s+`)

func regexpCompile(s string) (r *regexp.Regexp, err error) {
	normalized := whitespace.ReplaceAllStringFunc(
		strings.TrimFunc(s, unicode.IsSpace),
		func(s string) string {
			if err != nil {
				return s
			}
			if s[0] == '\\' {
				sr := strings.NewReader(s[1:])
				space, _, e := sr.ReadRune()
				if e != nil {
					err = e
				}
				return string([]rune{'\\', space})
			}
			return ""
		})
	if err != nil {
		return
	}
	return regexp.Compile(normalized)
}

func regexpMustCompile(s string) *regexp.Regexp {
	r, err := regexpCompile(s)
	if err != nil {
		panic(err)
	}
	return r
}

type MatcherConfig struct {
	Tracker  string `json:"tracker"`  // Matched tracker urls.
	Basename string `json:"basename"` // Matched (root) file basenames.
	Ext      string `json:"ext"`      // Matched (nested-)file extensions.
}

func (mc MatcherConfig) Matcher() *matcher {
	m := new(matcher)
	if mc.Tracker != "" {
		m.Tracker = regexpMustCompile(mc.Tracker)
	}
	if mc.Basename != "" {
		m.Basename = regexpMustCompile(mc.Basename)
	}
	if mc.Ext != "" {
		m.Ext = regexpMustCompile(mc.Ext)
	}
	return m
}

func (mc MatcherConfig) Validate() error {
	if _, err := regexpCompile(mc.Tracker); err != nil {
		return fmt.Errorf("matcher tracker: %v", err)
	}
	if _, err := regexpCompile(mc.Basename); err != nil {
		return fmt.Errorf("matcher basename: %v", err)
	}
	if _, err := regexpCompile(mc.Ext); err != nil {
		return fmt.Errorf("matcher ext: %v", err)
	}
	return nil
}

type HandlerConfig struct {
	Name  string        `json:"name"`  // A name for logging purposes.
	Watch string        `json:"watch"` // Matching .torrent file destination.
	Match MatcherConfig `json:"match"` // Describes .torrent files to handle.
}

func (c HandlerConfig) Handler() *Handler { return &Handler{c.Name, c.Watch, c.Match.Matcher()} }

func (hc HandlerConfig) Validate() error {
	if hc.Name == "" {
		return errors.New("nameless handler")
	}
	if hc.Watch == "" {
		return fmt.Errorf("handler %q: no watch directory.", hc.Name)
	}
	stat, err := os.Stat(hc.Watch)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("handler %q: watch is not a directory: %s", hc.Name, hc.Watch)
	}
	err = hc.Match.Validate()
	if err != nil {
		return fmt.Errorf("handler %q: %v", hc.Name, err)
	}
	return nil
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

func (config Config) Validate() error {
	if config.Path == "" {
		return errors.New("config: no path")
	}
	for _, watch := range config.Watch {
		stat, err := os.Stat(watch)
		if err != nil {
			return err
		}
		if !stat.IsDir() {
			return fmt.Errorf("config: watch is not a directory: %s", config.Watch)
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

func (c *Config) MakeHandlers() []*Handler {
	handlers := make([]*Handler, len(c.Handlers))
	for i := range c.Handlers {
		handlers[i] = c.Handlers[i].Handler()
	}
	return handlers
}
