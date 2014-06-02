package handler

import (
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/bmatsuo/gutterd/matcher"
)

type Config struct {
	Name    string     `json:"name"`             // A name for logging purposes.
	Watch   string     `json:"watch,omitempty"`  // Matching .torrent file destination.
	Script  []string   `json:"script,omitempty"` // Executed on matched files (should delete the file).
	Matcher *matcher.M `json:"matcher"`          // Describes .torrent files to handle.
}

func (c Config) Handler() (*Handler, error) {
	switch {
	case c.Watch != "" && len(c.Script) > 0:
		return nil, fmt.Errorf("both watch and script present")
	case c.Watch != "":
		return NewWatch(c.Name, c.Matcher, c.Watch), nil
	case len(c.Script) > 0:
		return NewScript(c.Name, c.Matcher, c.Script...)
	default:
		return nil, fmt.Errorf("neither watch no script present")
	}
}

func (hc Config) Validate() error {
	if hc.Name == "" {
		return errors.New("nameless handler")
	}
	if hc.Watch == "" && len(hc.Script) == 0 {
		return fmt.Errorf("neither watch nor script provided for %q", hc.Name)
	}
	if hc.Watch != "" && len(hc.Script) == 0 {
		return fmt.Errorf("both script and watch provided for %q", hc.Name)
	}
	for i := range hc.Script {
		_, err := template.New("").Parse(hc.Script[i])
		if err != nil {
			return fmt.Errorf("script %d for %q is invalid: %v", err)
		}
	}
	if hc.Watch != "" {
		stat, err := os.Stat(hc.Watch)
		if err != nil {
			return err
		}
		if !stat.IsDir() {
			return fmt.Errorf("watch for %q is not a directory: %s", hc.Name, hc.Watch)
		}
	}
	return nil
}
