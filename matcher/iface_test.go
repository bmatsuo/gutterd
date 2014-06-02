package matcher

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bmatsuo/torrent/metainfo"
)

var matchTestMatch = Must(Register(func() Interface {
	return &MatchTestMock{"test-match", nil}
}))
var matchTestNoMatch = Must(Register(func() Interface {
	return &MatchTestMock{"test-no-match", ErrNoMatch}
}))
var matchTestError = Must(Register(func() Interface {
	return &MatchTestMock{"test-error", fmt.Errorf("this is a test error message")}
}))

type MatchTestMock struct {
	name string
	err  error
}

func (m *MatchTestMock) Name() string {
	return m.name
}

func (m *MatchTestMock) MatchTorrent(t *metainfo.Metainfo) error {
	return m.err
}

type MTestJSON struct {
	JSON   string
	Meta   *metainfo.Metainfo
	Result MTestResult
}
type MTestResult uint

const (
	MTestMatch MTestResult = iota
	MTestNoMatch
	MTestError
)

func testM(t *testing.T, test MTestJSON) {
	var m M
	err := json.Unmarshal([]byte(test.JSON), &m)
	if err != nil {
		t.Errorf("unmarshal %s %v", test.JSON, err)
		return
	}
	switch err = m.MatchTorrent(test.Meta); err {
	case nil:
		var qual string
		switch test.Result {
		case MTestMatch:
			return
		case MTestNoMatch:
			qual = " (expected no-match)"
		case MTestError:
			qual = " (expected error)"
		}
		t.Errorf("%s unexpected match%q", test.JSON, qual)
	case ErrNoMatch:
		var qual string
		switch test.Result {
		case MTestMatch:
			qual = " (expected match)"
		case MTestNoMatch:
			return
		case MTestError:
			qual = " (expected error)"
		}
		t.Errorf("%s unexpected no-match%s", test.JSON, qual)
	default:
		var qual string
		switch test.Result {
		case MTestMatch:
			qual = " (expected match)"
		case MTestNoMatch:
			qual = " (expected no-match)"
		case MTestError:
			return
		}
		t.Errorf("%s unexpected error%s: %v", test.JSON, qual, err)
	}
}

func TestM(t *testing.T) {
	for _, test := range []MTestJSON{
		{`{"match":"test-match"}`, nil, MTestMatch},
		{`{"match":"test-no-match"}`, nil, MTestNoMatch},
		{`{"match":"test-error"}`, nil, MTestError},
	} {
		testM(t, test)
	}
}

func TestRegister(t *testing.T) {
	// register an interface with a unique name
	test := MatchTestMock{"test-register-dup-name", nil}
	m, err := Register(func() Interface { return &test })
	if err != nil {
		t.Fatalf("unable to register unique name", err)
	}
	if m == nil {
		t.Fatalf("nil matcher returned")
	}

	// attempt to register an interface with the same name
	m, err = Register(func() Interface { return &test })
	if err == nil {
		t.Fatalf("registering duplicate name succeeded")
	}
	if m != nil {
		t.Fatalf("non-nil matcher returned")
	}
}
