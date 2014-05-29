package handler

import (
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/bmatsuo/gutterd/matcher"
)

type Config struct {
	Name   string         `json:"name"`             // A name for logging purposes.
	Watch  string         `json:"watch,omitempty"`  // Matching .torrent file destination.
	Script []string       `json:"script,omitempty"` // Executed on matched files (should delete the file).
	Match  matcher.Config `json:"match"`            // Describes .torrent files to handle.
}

func (c Config) Handler() (*Handler, error) {
	switch {
	case c.Watch != "" && len(c.Script) > 0:
		return nil, fmt.Errorf("both watch and script present")
	case c.Watch != "":
		return NewWatch(c.Name, c.Match.Matcher(), c.Watch), nil
	case len(c.Script) > 0:
		return NewScript(c.Name, c.Match.Matcher(), c.Script...)
	default:
		return nil, fmt.Errorf("nother watch no script present")
	}
}

func (hc Config) Validate() error {
	if hc.Name == "" {
		return errors.New("nameless handler")
	}
	if hc.Watch == "" && len(hc.Script) == 0 {
		return fmt.Errorf("either watch or script must be provided for handler %q", hc.Name)
	}
	if hc.Watch != "" && len(hc.Script) == 0 {
		return fmt.Errorf("script and watch may not both be provided for %q", hc.Name)
	}
	for i := range hc.Script {
		_, err := template.New("").Parse(hc.Script[i])
		if err != nil {
			return fmt.Errorf("script %d for %q is invalid: %v", err)
		}
	}
	stat, err := os.Stat(hc.Watch)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("watch for handler %q is not a directory: %s", hc.Name, hc.Watch)
	}
	err = hc.Match.Validate()
	if err != nil {
		return fmt.Errorf("handler %q %v", hc.Name, err)
	}
	return nil
}
