package matcher

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"github.com/bmatsuo/torrent/metainfo"
)

var JSONKey = "match"

type Error struct {
	name string
	err  error
}

// Err return the cause of err.
func (err Error) Err() error {
	return err.err
}

func (err Error) Error() string {
	return fmt.Sprintf("matcher %q %s", err.name, err.err.Error())
}

var ErrNoMatch = fmt.Errorf("no match")

type JSONRegexp struct {
	*regexp.Regexp
}

func (r *JSONRegexp) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *JSONRegexp) UnarshalJSON(p []byte) error {
	var s string
	err := json.Unmarshal(p, &s)
	if err != nil {
		return err
	}
	r.Regexp, err = regexp.Compile(s)
	return err
}

// Interface is implemented by matchers available to guttered when a
// constructor is registered with Register().
type Interface interface {
	Name() string
	MatchTorrent(t *metainfo.Metainfo) error
}

type Func func(*metainfo.Metainfo) error

func (fn Func) MatchTorrent(t *metainfo.Metainfo) error {
	return fn(t)
}

var _M = struct {
	sync.Mutex
	m map[string]func() Interface
}{
	m: make(map[string]func() Interface),
}

type M struct {
	name string
	m    Interface
}

func (m *M) Name() string {
	return m.name
}

func (m *M) MatchTorrent(t *metainfo.Metainfo) error {
	if m.m == nil {
		return fmt.Errorf("no match")
	}
	return m.m.MatchTorrent(t)
}

func (m *M) MarshalJSON() ([]byte, error) {
	if m.m == nil {
		return nil, fmt.Errorf("unitialized")
	}
	p, err := json.Marshal(m.m)
	if err != nil {
		return nil, Error{m.name, err}
	}
	_map := make(map[string]*json.RawMessage)
	err = json.Unmarshal(p, _map)
	if err != nil {
		return nil, Error{m.name, fmt.Errorf("bad serialization")}
	}
	jsonName, _ := json.Marshal(m.name)
	_map[JSONKey] = (*json.RawMessage)(&jsonName)
	return json.Marshal(_map)
}

func (m *M) UnmarshalJSON(p []byte) error {
	_map := make(map[string]*json.RawMessage)
	err := json.Unmarshal(p, &_map)
	if err != nil {
		return err
	}
	_name := _map[JSONKey]
	if _name == nil {
		return fmt.Errorf("missing %s", JSONKey)
	}
	err = json.Unmarshal(*_name, &m.name)
	if err != nil {
		return err
	}
	if m.name == "" {
		return fmt.Errorf("empty %s", JSONKey)
	}
	_M.Lock()
	fn := _M.m[m.name]
	_M.Unlock()
	if fn == nil {
		return Error{m.name, fmt.Errorf("unrecognized")}
	}
	m.m = fn()
	err = json.Unmarshal(p, m.m)
	if err != nil {
		return Error{m.name, err}
	}
	return nil
}

func Must(m *M, err error) *M {
	if err != nil {
		panic(err)
	}
	return m
}

func Register(fn func() Interface) (*M, error) {
	if fn == nil {
		return nil, fmt.Errorf("nil interface")
	}
	name := fn().Name()
	_M.Lock()
	defer _M.Unlock()
	_, ok := _M.m[name]
	if ok {
		return nil, fmt.Errorf("already registered")
	}
	_M.m[name] = fn
	return &M{name, nil}, nil
}
