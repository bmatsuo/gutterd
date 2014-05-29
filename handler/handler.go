package handler

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/bmatsuo/gutterd/matcher"
	"github.com/bmatsuo/torrent/metainfo"
)

// A Handler type's only function is to move matching torrents into
// media-specific client watch directories.
type Handler struct {
	Name             string   // Unique name for the Handler.
	Watch            string   // Destination for .torrent files (watched by a client).
	Script           []string // Script to run on matched files (should destroy the path).
	*matcher.Matcher          // Acts as a Matcher.
	scriptTemplate   *template.Template
}

var ErrNoMatch = fmt.Errorf("no match")

func (h *Handler) Handle(path string, meta *metainfo.Metainfo) error {
	if !h.Match(meta) {
		return ErrNoMatch
	}
	if h.Watch != "" {
		return h.handleWatch(path, meta)
	}
	if len(h.Script) > 0 {
		return h.handleScript(path, meta)
	}
	return fmt.Errorf("%q matched but has no action for %q", h.Name, meta.Info.Name)
}

// handleWatch moves path to h.Watch
func (h *Handler) handleWatch(path string, meta *metainfo.Metainfo) error {
	mvpath := filepath.Join(h.Watch, filepath.Base(path))
	return os.Rename(path, mvpath)
}

// handleScript creates a temporary script from the template, executes it, and removes it.
func (h *Handler) handleScript(path string, meta *metainfo.Metainfo) error {
	contxt := map[string]interface{}{
		"Path": path,
	}
	f, err := ioutil.TempFile("", "gutterd-script-")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	err = h.scriptTemplate.Execute(f, contxt)
	if err != nil {
		f.Close()
		return err
	}
	err = f.Close()
	if err != nil {
		return fmt.Errorf("%q couldn't create script file: %v", err)
	}
	cmd := exec.Command("/bin/bash", f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%q script failed: %v", err)
	}
	return nil
}

func (h *Handler) String() string { return h.Name }
