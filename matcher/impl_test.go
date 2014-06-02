package matcher

import (
	"testing"

	"github.com/bmatsuo/torrent/metainfo"
)

func TestMatchExt(t *testing.T) {
	for _, test := range []MTestJSON{
		{
			`{"match":"ext","pattern":"iso"}`,
			&metainfo.Metainfo{Info: metainfo.Info{Name: "test.iso"}},
			MTestMatch,
		},
		{
			`{"match":"ext","pattern":"iso"}`,
			&metainfo.Metainfo{Info: metainfo.Info{Files: []metainfo.FileInfo{
				{Path: []string{"path", "to", "test.iso"}},
			}}},
			MTestMatch,
		},
		{
			`{"match":"ext","pattern":"iso"}`,
			&metainfo.Metainfo{Info: metainfo.Info{Name: "iso.no"}},
			MTestNoMatch,
		},
		{
			`{"match":"ext","pattern":"iso"}`,
			&metainfo.Metainfo{Info: metainfo.Info{Files: []metainfo.FileInfo{
				{Path: []string{"path", "to", "not-an.iso", "iso.no"}},
			}}},
			MTestNoMatch,
		},
		{
			`{"match":"ext","pattern":"iso"}`,
			nil,
			MTestError,
		},
	} {
		testM(t, test)
	}
}

func TestMatchBase(t *testing.T) {
	for _, test := range []MTestJSON{
		{
			`{"match":"base","pattern":"iso"}`,
			&metainfo.Metainfo{Info: metainfo.Info{Name: "test.iso"}},
			MTestMatch,
		},
		{
			`{"match":"base","pattern":"iso"}`,
			&metainfo.Metainfo{Info: metainfo.Info{Files: []metainfo.FileInfo{
				{Path: []string{"path", "to", "test.iso"}},
			}}},
			MTestMatch,
		},
		{
			`{"match":"base","pattern":"iso"}`,
			&metainfo.Metainfo{Info: metainfo.Info{Name: "test.no"}},
			MTestNoMatch,
		},
		{
			`{"match":"base","pattern":"iso"}`,
			&metainfo.Metainfo{Info: metainfo.Info{Files: []metainfo.FileInfo{
				{Path: []string{"path", "to", "not-an-iso", "test.no"}},
			}}},
			MTestNoMatch,
		},
		{
			`{"match":"base","pattern":"iso"}`,
			nil,
			MTestError,
		},
	} {
		testM(t, test)
	}
}

func TestMatchAnnounce(t *testing.T) {
	for _, test := range []MTestJSON{
		{`{"match":"announce","pattern":"example[.]com"}`, &metainfo.Metainfo{Announce: "http://example.com/announce"}, MTestMatch},
		{`{"match":"announce","pattern":"example[.]com"}`, &metainfo.Metainfo{Announce: "https://ex1.example.com/announce"}, MTestMatch},
		{`{"match":"announce","pattern":"example[.]com"}`, &metainfo.Metainfo{Announce: "https://ex1.example.io/announce"}, MTestNoMatch},
		{`{"match":"announce","pattern":"example[.]com"}`, nil, MTestError},
	} {
		testM(t, test)
	}
}

func TestAll(t *testing.T) {
	for _, test := range []MTestJSON{
		{`{"match":"all"}`, nil, MTestMatch},
		{`{"match":"all","of":[]}`, nil, MTestMatch},
		{`{"match":"all","of":[{"match":"test-match"}]}`, nil, MTestMatch},
		{`{"match":"all","of":[{"match":"test-no-match"}]}`, nil, MTestNoMatch},
		{`{"match":"all","of":[{"match":"test-error"}]}`, nil, MTestError},
		{`{"match":"all","of":[{"match":"test-match"},{"match":"test-no-match"}]}`, nil, MTestNoMatch},
		{`{"match":"all","of":[{"match":"test-no-match"},{"match":"test-match"}]}`, nil, MTestNoMatch},
		{`{"match":"all","of":[{"match":"test-match"},{"match":"test-error"}]}`, nil, MTestError},
		{`{"match":"all","of":[{"match":"test-error"},{"match":"test-match"}]}`, nil, MTestError},
		{`{"match":"all","of":[{"match":"test-no-match"},{"match":"test-error"}]}`, nil, MTestNoMatch},
		{`{"match":"all","of":[{"match":"test-error"},{"match":"test-no-match"}]}`, nil, MTestError},
	} {
		testM(t, test)
	}
}

func TestAny(t *testing.T) {
	for _, test := range []MTestJSON{
		{`{"match":"any"}`, nil, MTestNoMatch},
		{`{"match":"any","of":[]}`, nil, MTestNoMatch},
		{`{"match":"any","of":[{"match":"test-match"}]}`, nil, MTestMatch},
		{`{"match":"any","of":[{"match":"test-no-match"}]}`, nil, MTestNoMatch},
		{`{"match":"any","of":[{"match":"test-error"}]}`, nil, MTestError},
		{`{"match":"any","of":[{"match":"test-match"},{"match":"test-no-match"}]}`, nil, MTestMatch},
		{`{"match":"any","of":[{"match":"test-no-match"},{"match":"test-match"}]}`, nil, MTestMatch},
		{`{"match":"any","of":[{"match":"test-match"},{"match":"test-error"}]}`, nil, MTestMatch},
		{`{"match":"any","of":[{"match":"test-error"},{"match":"test-match"}]}`, nil, MTestMatch},
		{`{"match":"any","of":[{"match":"test-no-match"},{"match":"test-error"}]}`, nil, MTestError},
		{`{"match":"any","of":[{"match":"test-error"},{"match":"test-no-match"}]}`, nil, MTestError},
	} {
		testM(t, test)
	}
}
