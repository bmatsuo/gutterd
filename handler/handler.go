package handler

import "github.com/bmatsuo/gutterd/matcher"

// A Handler type's only function is to move matching torrents into
// media-specific client watch directories.
type Handler struct {
	Name             string   // Unique name for the Handler.
	Watch            string   // Destination for .torrent files (watched by a client).
	Script           []string // Script to run on matched files (should destroy the path).
	*matcher.Matcher          // Acts as a Matcher.
}

func (h *Handler) String() string { return h.Name }
