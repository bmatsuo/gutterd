package handler

import (
	"errors"
	"fmt"
	"os"

	"github.com/bmatsuo/gutterd/matcher"
)

type Config struct {
	Name   string         `json:"name"`   // A name for logging purposes.
	Watch  string         `json:"watch"`  // Matching .torrent file destination.
	Script string         `json:"script"` // Executed on matched files (should delete the file).
	Match  matcher.Config `json:"match"`  // Describes .torrent files to handle.
}

func (c Config) Handler() *Handler {
	return &Handler{
		Name:    c.Name,
		Watch:   c.Watch,
		Script:  c.Script,
		Matcher: c.Match.Matcher(),
	}
}

func (hc Config) Validate() error {
	if hc.Name == "" {
		return errors.New("nameless handler")
	}
	if hc.Watch == "" && hc.Script == "" {
		return fmt.Errorf("either watch or script must be provided for handler %q", hc.Name)
	}
	if hc.Watch != "" && hc.Script != "" {
		return fmt.Errorf("script and watch may not both be provided for %q", hc.Name)
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
