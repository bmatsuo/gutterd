package handler

import (
	"errors"
	"fmt"
	"os"

	"github.com/bmatsuo/gutterd/matcher"
)

type Config struct {
	Name  string         `json:"name"`  // A name for logging purposes.
	Watch string         `json:"watch"` // Matching .torrent file destination.
	Match matcher.Config `json:"match"` // Describes .torrent files to handle.
}

func (c Config) Handler() *Handler { return &Handler{c.Name, c.Watch, c.Match.Matcher()} }

func (hc Config) Validate() error {
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
