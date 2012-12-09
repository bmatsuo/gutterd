package handler

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

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

func (mc MatcherConfig) Matcher() *Matcher {
	m := new(Matcher)
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
		return fmt.Errorf("Matcher tracker: %v", err)
	}
	if _, err := regexpCompile(mc.Basename); err != nil {
		return fmt.Errorf("Matcher basename: %v", err)
	}
	if _, err := regexpCompile(mc.Ext); err != nil {
		return fmt.Errorf("Matcher ext: %v", err)
	}
	return nil
}
