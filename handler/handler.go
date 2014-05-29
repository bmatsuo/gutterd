package handler

import (
	"fmt"
	"text/template"

	"github.com/bmatsuo/gutterd/matcher"
	"github.com/bmatsuo/torrent/metainfo"
)

var ErrNoMatch = fmt.Errorf("no match")

type Interface interface {
	Name() string
	Handle(path string, meta *metainfo.Metainfo) error
}

// A Handler type's only function is to move matching torrents into
// media-specific client watch directories.
type Handler struct {
	name             string   // Unique name for the Handler.
	exec             Exec     // A way to act on files.
	watch            string   // Destination for .torrent files (watched by a client).
	script           []string // Script to run on matched files (should destroy the path).
	*matcher.Matcher          // Implements matcher.Interface
	scriptTemplate   *template.Template
}

func New(name string, m *matcher.Matcher, exec Exec) *Handler {
	return &Handler{
		name:    name,
		exec:    exec,
		Matcher: m,
	}
}

func NewWatch(name string, m *matcher.Matcher, watch string) *Handler {
	return New(name, m, WatchDir(watch))
}

func NewScript(name string, m *matcher.Matcher, script ...string) (*Handler, error) {
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
	if !h.Match(meta) {
		return ErrNoMatch
	}
	return h.exec.Run(path, meta)
}

func (h *Handler) String() string {
	return h.name
}
