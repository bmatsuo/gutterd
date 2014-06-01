package matcher

import (
	"encoding/json"
	"path/filepath"

	"github.com/bmatsuo/torrent/metainfo"
)

var matchExt = Must(Register(func() Interface { return new(MatchExt) }))
var matchAnnounce = Must(Register(func() Interface { return new(MatchAnnounce) }))
var matchBase = Must(Register(func() Interface { return new(MatchBase) }))
var matchAll = Must(Register(func() Interface { return new(MatchAll) }))
var matchAny = Must(Register(func() Interface { return new(MatchAny) }))

type MatchExt struct {
	Pattern *JSONRegexp `json:"pattern"`
}

func (m *MatchExt) Name() string {
	return "ext"
}

func (m *MatchExt) MatchTorrent(t *metainfo.Metainfo) error {
	if m.Pattern == nil {
		return ErrNoMatch
	}
	if len(t.Info.Files) == 0 {
		if m.Pattern.MatchString(filepath.Ext(t.Info.Name)) {
			return nil
		}
		return ErrNoMatch
	}
	for _, f := range t.Info.Files {
		p := f.Path
		ext := filepath.Ext(p[len(p)-1])
		if m.Pattern.MatchString(ext) {
			return nil
		}
	}
	return ErrNoMatch
}

type MatchBase struct {
	Pattern *JSONRegexp `json:"pattern"`
}

func (m *MatchBase) Name() string {
	return "base"
}

func (m *MatchBase) MatchTorrent(t *metainfo.Metainfo) error {
	if m.Pattern == nil {
		return ErrNoMatch
	}
	if len(t.Info.Files) == 0 {
		if m.Pattern.MatchString(filepath.Ext(t.Info.Name)) {
			return nil
		}
		return ErrNoMatch
	}
	for _, f := range t.Info.Files {
		p := f.Path
		base := p[len(p)-1]
		if m.Pattern.MatchString(base) {
			return nil
		}
	}
	return ErrNoMatch
}

type MatchAnnounce struct {
	Pattern *JSONRegexp `json:"pattern"`
}

func (m *MatchAnnounce) Name() string {
	return "tracker"
}

func (m *MatchAnnounce) MatchTorrent(t *metainfo.Metainfo) error {
	if m.Pattern == nil {
		return ErrNoMatch
	}
	if m.Pattern.MatchString(filepath.Ext(t.Announce)) {
		return nil
	}
	return ErrNoMatch
}

type MatchAll []*M

type msRawJSON MatchAll

func (m MatchAll) Name() string {
	return "all"
}

func (m *MatchAll) UnmarshalJSON(p []byte) error {
	_map := make(map[string]*json.RawMessage)
	err := json.Unmarshal(p, _map)
	if err != nil {
		return err
	}
	_p, ok := _map["of"]
	if !ok {
		return nil
	}
	return json.Unmarshal(*_p, (*msRawJSON)(m))
}

func (m *MatchAll) MarshalJSON() ([]byte, error) {
	_map := make(map[string]*json.RawMessage)
	p, err := json.Marshal(msRawJSON(*m))
	if err != nil {
		return nil, err
	}
	_map["of"] = (*json.RawMessage)(&p)
	return json.Marshal(_map)
}

func (m MatchAll) MatchTorrent(t *metainfo.Metainfo) error {
	for i := range m {
		err := m[i].MatchTorrent(t)
		if err != nil {
			return err
		}
	}
	return nil
}

type MatchAny []*M

func (m MatchAny) Name() string {
	return "any"
}

func (m *MatchAny) UnmarshalJSON(p []byte) error {
	_map := make(map[string]*json.RawMessage)
	err := json.Unmarshal(p, _map)
	if err != nil {
		return err
	}
	_p, ok := _map["of"]
	if !ok {
		return nil
	}
	return json.Unmarshal(*_p, (*msRawJSON)(m))
}

func (m *MatchAny) MarshalJSON() ([]byte, error) {
	_map := make(map[string]*json.RawMessage)
	p, err := json.Marshal(msRawJSON(*m))
	if err != nil {
		return nil, err
	}
	_map["of"] = (*json.RawMessage)(&p)
	return json.Marshal(_map)
}

func (m MatchAny) MatchTorrent(t *metainfo.Metainfo) error {
	for i := range m {
		err := m[i].MatchTorrent(t)
		if err == ErrNoMatch {
			continue
		}
		if err != nil {
			return err
		}
		return nil
	}
	return ErrNoMatch
}
