// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
package matcher provides torrent matchers
*/
package matcher

/*  Filename:    matcher.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     2012-03-04 19:59:29.661375 -0800 PST
 *  Description:
 */

import (
	"path/filepath"
	"regexp"

	"github.com/bmatsuo/torrent/metainfo"
)

// Matcher matches torrents (Metadata).
type Matcher struct {
	Tracker  *regexp.Regexp
	Basename *regexp.Regexp
	Ext      *regexp.Regexp
}

// Match a torrent against the patterns of m. If all non-nil patterns match the
// corresponding fields in torrent, then the method returns true.
func (m *Matcher) Match(torrent *metainfo.Metainfo) bool {
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
	if m.Basename != nil {
		basename := filepath.Base(torrent.Info.Name)
		if !m.Basename.MatchString(basename) {
			return false
		}
	}
	return true
}
