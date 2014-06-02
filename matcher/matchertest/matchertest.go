package matchertest

import (
	"fmt"

	"github.com/bmatsuo/gutterd/matcher"
	"github.com/bmatsuo/torrent/metainfo"
)

var Match = matcher.Must(matcher.Register(func() matcher.Interface {
	return &Mock{"test-match", nil}
}))
var NoMatch = matcher.Must(matcher.Register(func() matcher.Interface {
	return &Mock{"test-no-match", matcher.NoMatch}
}))
var Error = matcher.Must(matcher.Register(func() matcher.Interface {
	return &Mock{"test-error", fmt.Errorf("this is a test error message")}
}))

type Mock struct {
	name string
	err  error
}

func (m *Mock) Name() string {
	return m.name
}

func (m *Mock) MatchTorrent(t *metainfo.Metainfo) error {
	return m.err
}
