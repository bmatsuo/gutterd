// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

/*  Filename:    metadata.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     2012-03-04 20:29:46.043613 -0800 PST
 *  Description:
 */

import (
	"fmt"
	"io/ioutil"

	"github.com/bmatsuo/gorrent/bencode"
)

// One file in a multi-file Metadata object.
type FileInfo struct {
	Path   []string // File path components.
	Length int64    // Length in bytes.
	MD5Sum string   // Optional.
}

// The main contents of a Metadata type
type TorrentInfo struct {
	Name        string      // Name of file (single-file mode) or directory (multi-file mode)
	Files       []*FileInfo // Nil if and only if single-file mode
	MD5Sum      string      // Optional -- Non-empty if and only if single-file mode.
	Pieces      string      // SHA-1 hash values of all pieces
	PieceLength int64       // Length in bytes.
	Private     bool        // Optional
}

// Returns true if info is in Single file mode.
func (info *TorrentInfo) SingleFileMode() bool { return info.Files == nil }

// The contents of a .torrent file.
type Metadata struct {
	Info         *TorrentInfo // Required
	Announce     string       // Required
	CreationDate int64        // Optional
	Encoding     string       // Optional
	CreatedBy    string       // Optional
	Comment      string       // Optional
}

func tryCastKey(m map[string]interface{}, key string, action func(interface{}), required bool) {
	tryCast(key, m[key], action, required)
}

func tryCast(name string, v interface{}, action func(interface{}), required bool) {
	defer func() {
		if e := recover(); e != nil {
			if required {
				panic(fmt.Errorf("%s: %v", name, e))
			}
		}
	}()
	action(v)
}

func ReadMetadataFile(torrent string) (meta *Metadata, err error) {
	p, err := ioutil.ReadFile(torrent)
	if err != nil {
		return nil, err
	}
	data, err := bencode.NewDecoder(p).DecodeAll()
	switch {
	case err != nil:
		return nil, err
	case len(data) > 1:
		return nil, fmt.Errorf("unexpected bencoded data")
	case len(data) == 0:
		return nil, fmt.Errorf("no bencoded data")
	}
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("type error: %v", e)
		}
	}()
	var _meta map[string]interface{}
	tryCast("torrent", data[0], func(v interface{}) { _meta = v.(map[string]interface{}) }, true)
	meta = new(Metadata)
	tryCastKey(_meta, "announce", func(v interface{}) { meta.Announce = v.(string) }, true)
	tryCastKey(_meta, "encoding", func(v interface{}) { meta.Encoding = v.(string) }, false)
	tryCastKey(_meta, "comment", func(v interface{}) { meta.Comment = v.(string) }, false)
	tryCastKey(_meta, "created by", func(v interface{}) { meta.CreatedBy = v.(string) }, false)
	tryCastKey(_meta, "creation date", func(v interface{}) { meta.CreationDate = v.(int64) }, false)
	var _info map[string]interface{}
	tryCastKey(_meta, "info", func(v interface{}) { _info = v.(map[string]interface{}) }, true)
	info := new(TorrentInfo)
	meta.Info = info
	tryCastKey(_info, "name", func(v interface{}) { info.Name = v.(string) }, true)
	tryCastKey(_info, "pieces", func(v interface{}) { info.Pieces = v.(string) }, true)
	tryCastKey(_info, "md5sum", func(v interface{}) { info.MD5Sum = v.(string) }, false)
	tryCastKey(_info, "private", func(v interface{}) { info.Private = v.(int64) == 1 }, false)
	tryCastKey(_info, "piece length", func(v interface{}) { info.PieceLength = v.(int64) }, true)
	var _fileIs []interface{}
	tryCastKey(_info, "files", func(v interface{}) { _fileIs = v.([]interface{}) }, false)
	for i, _fileI := range _fileIs {
		var _file map[string]interface{}
		tryCast(fmt.Sprintf("file %d", i), _fileI,
			func(v interface{}) { _file = v.(map[string]interface{}) }, true)
		file := new(FileInfo)
		tryCastKey(_file, "md5sum", func(v interface{}) { file.MD5Sum = v.(string) }, false)
		tryCastKey(_file, "length", func(v interface{}) { file.Length = v.(int64) }, true)
		var path []interface{}
		tryCastKey(_file, "path", func(v interface{}) { path = v.([]interface{}) }, true)
		for j, elem := range path {
			tryCast(fmt.Sprintf("file %d: path element %d", i, j), elem,
				func(v interface{}) { file.Path = append(file.Path, v.(string)) }, true)
		}
		info.Files = append(info.Files, file)
	}
	return meta, nil
}
