package handler

import (
	"fmt"
	"text/template"

	"github.com/bmatsuo/gutterd/matcher"
	"github.com/bmatsuo/torrent/metainfo"
)

var NoMatch = fmt.Errorf("no match")

type Interface interface {
	Name() string
	Handle(path string, meta *metainfo.Metainfo) error
}

// A Handler type's only function is to match torrents and add them to torrent
// clients.
type Handler struct {
	name              string   // Unique name for the Handler.
	exec              Exec     // A way to act on files.
	watch             string   // Destination for .torrent files (watched by a client).
	script            []string // Script to run on matched files (should destroy the path).
	matcher.Interface          // Implements matcher.Interface
	scriptTemplate    *template.Template
}

func New(name string, m matcher.Interface, exec Exec) *Handler {
	return &Handler{
		name:      name,
		exec:      exec,
		Interface: m,
	}
}

func NewWatch(name string, m matcher.Interface, watch string) *Handler {
	return New(name, m, WatchDir(watch))
}

func NewScript(name string, m matcher.Interface, script ...string) (*Handler, error) {
	st, err := NewScriptTemplate(script...)
	if err != nil {
		return nil, err
	}
	return New(name, m, st), nil
}

func (h *Handler) Name() string {
	return h.name
}

func (h *Handler) Handle(path string, meta *metainfo.Metainfo) error {
	err := h.MatchTorrent(meta)
	if err == matcher.NoMatch {
		return NoMatch
	}
	if err != nil {
		return err
	}
	return h.exec.Run(path, meta)
}

func (h *Handler) String() string {
	return h.name
}
