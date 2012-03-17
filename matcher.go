// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    matcher.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     2012-03-04 19:59:29.661375 -0800 PST
 *  Description: 
 */

import (
	"path/filepath"
	"regexp"
)

// Matched against torrents (Metadata) by Handler types.
type matcher struct {
	Tracker  *regexp.Regexp
	Basename *regexp.Regexp
	Ext      *regexp.Regexp
}

// Match a torrent against the patterns of m. If all non-nil patterns match the
// corresponding fields in torrent, then the method returns true.
func (m *matcher) Match(torrent *Metadata) bool {
	if m.Tracker != nil {
		if !m.Tracker.MatchString(torrent.Announce) {
			return false
		}
	}
	if m.Ext != nil {
		var exts []string
		if torrent.Info.SingleFileMode() {
			exts = append(exts, filepath.Ext(torrent.Info.Name))
		} else {
			for _, file := range torrent.Info.Files {
				path := file.Path
				exts = append(exts, filepath.Ext(path[len(path)-1]))
			}
		}
		matches := false
		for _, ext := range exts {
			if m.Ext.MatchString(ext) {
				matches = true
			}
		}
		if !matches {
			return false
		}
	}
	return true
}

// A Handler type's only function is to move matching torrents into
// media-specific client watch directories.
type Handler struct {
	Name     string // Unique name for the Handler.
	Watch    string // Destination for .torrent files (watched by a client).
	*matcher        // Acts as a matcher.
}

func (h *Handler) String() string { return h.Name }
