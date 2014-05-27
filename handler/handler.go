package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/bmatsuo/gutterd/matcher"
	"github.com/bmatsuo/gutterd/metadata"
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

func (h *Handler) Handle(path string, meta *metadata.Metadata) error {
	if !h.Match(meta) {
		return ErrNoMatch
	}
	if h.Watch != "" {
		mvpath := filepath.Join(h.Watch, filepath.Base(path))
		if err := os.Rename(path, mvpath); err != nil {
			return err
		}
	}
	if len(h.Script) > 0 {
		h.handleScript(path, meta)
	}
	return nil
}

func (h *Handler) handleScript(path string, meta *metadata.Metadata) error {
	return nil
}

func (h *Handler) String() string { return h.Name }
