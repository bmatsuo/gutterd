package handler

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bmatsuo/gutterd/matcher/matchertest"
	"github.com/bmatsuo/torrent/metainfo"
)

func setupHandlerTest(t *testing.T) string {
	dir, err := ioutil.TempDir("", "gutterd-handler-test-")
	if err != nil {
		t.Fatalf("couldn't set up handler test: %v", err)
	}
	err = os.MkdirAll(filepath.Join(dir, "watch"), 0755)
	if err != nil {
		teardownHandlerTest(t, dir)
		t.Fatal("unable to create watch directory")
	}
	err = os.MkdirAll(filepath.Join(dir, "downloads"), 0755)
	if err != nil {
		teardownHandlerTest(t, dir)
		t.Fatal("unable to create downloads directory")
	}
	return dir
}

func teardownHandlerTest(t *testing.T, path string) {
	t.Logf("tearing down")
	err := os.RemoveAll(path)
	if err != nil {
		t.Errorf("couldn't tear down handler test: %v", err)
	}
}

type readRepeater struct {
	i int
	s string
}

func newReadRepeater(s string) *readRepeater {
	return &readRepeater{0, s}
}

func (r *readRepeater) Read(p []byte) (int, error) {
	var read int
	n := len(p)
	for n > 0 {
		m := copy(p[read:], r.s[r.i:])
		n -= m
		read += m
		r.i += m
		r.i %= len(r.s)
	}
	return read, nil
}

func TestHandlerWatch(t *testing.T) {
	dir := setupHandlerTest(t)
	defer teardownHandlerTest(t, dir)

	w, err := metainfo.NewWriterSingle(1<<19, "hello.txt")
	if err != nil {
		t.Fatalf("unable to create test writer: %v", err)
	}
	n, err := io.CopyN(w, newReadRepeater("hello torrents\n"), 1<<20)
	if err != nil {
		t.Fatalf("unable to write test data: %v", err)
	}
	if n != 1<<20 {
		t.Fatalf("bad write length: %v", n)
	}
	meta, err := w.Metainfo("", "http://example.com/announce")
	if err != nil {
		t.Fatalf("error constructing metainfo: %v", err)
	}
	watchdir := filepath.Join(dir, "watch")
	indir := filepath.Join(dir, "downloads")
	tpath := filepath.Join(indir, "hello.torrent")
	watchpath := filepath.Join(watchdir, "hello.torrent")
	err = metainfo.WriteFile(tpath, meta, 0644)
	if err != nil {
		t.Fatalf("unable to write torrent file: %v", err)
	}

	hmatch := NewWatch("hmatch", matchertest.Match, watchdir)
	err = hmatch.Handle(tpath, meta)
	if err != nil {
		t.Fatalf("handle %v", err)
	}
	_, err = os.Stat(watchpath)
	if err != nil {
		t.Fatalf("unable to locate torrent in watch directory: %v", err)
	}
	err = os.Rename(watchpath, tpath)
	if err != nil {
		t.Fatalf("unable to move torrent file: %v", err)
	}

	hnomatch := NewWatch("handler-no-match", matchertest.NoMatch, watchdir)
	err = hnomatch.Handle(tpath, meta)
	if err != NoMatch {
		t.Fatalf("handle %v", err)
	}
	_, err = os.Stat(watchpath)
	if !os.IsNotExist(err) {
		t.Fatalf("torrent was relocated")
	}
	_, err = os.Stat(tpath)
	if err != nil {
		t.Fatalf("unable to locate torrent")
	}

	herror := NewWatch("handler-error", matchertest.Error, watchdir)
	err = herror.Handle(tpath, meta)
	if err == nil || err == NoMatch {
		t.Fatalf("handle %v", err)
	}
	_, err = os.Stat(watchpath)
	if !os.IsNotExist(err) {
		t.Fatalf("torrent was relocated")
	}
	_, err = os.Stat(tpath)
	if err != nil {
		t.Fatalf("unable to locate torrent")
	}
}

func TestHandlerScript(t *testing.T) {
	dir := setupHandlerTest(t)
	defer teardownHandlerTest(t, dir)

	w, err := metainfo.NewWriterSingle(1<<19, "hello.txt")
	if err != nil {
		t.Fatalf("unable to create test writer: %v", err)
	}
	n, err := io.CopyN(w, newReadRepeater("hello torrents\n"), 1<<20)
	if err != nil {
		t.Fatalf("unable to write test data: %v", err)
	}
	if n != 1<<20 {
		t.Fatalf("bad write length: %v", n)
	}
	meta, err := w.Metainfo("", "http://example.com/announce")
	if err != nil {
		t.Fatalf("error constructing metainfo: %v", err)
	}
	indir := filepath.Join(dir, "downloads")
	tpath := filepath.Join(indir, "hello.torrent")
	err = metainfo.WriteFile(tpath, meta, 0644)
	if err != nil {
		t.Fatalf("unable to write torrent file: %v", err)
	}
	script := []string{"rm {{.Path}}"}

	hmatch, err := NewScript("hmatch", matchertest.Match, script...)
	if err != nil {
		t.Fatalf("script error: %v", err)
	}
	err = hmatch.Handle(tpath, meta)
	if err != nil {
		t.Fatalf("handle %v", err)
	}
	_, err = os.Stat(tpath)
	if !os.IsNotExist(err) {
		t.Fatalf("torrent file was not removed: %v", err)
	}
	err = metainfo.WriteFile(tpath, meta, 0644)
	if err != nil {
		t.Fatalf("unable to write torrent file: %v", err)
	}

	hnomatch, err := NewScript("handler-no-match", matchertest.NoMatch, script...)
	err = hnomatch.Handle(tpath, meta)
	if err != NoMatch {
		t.Fatalf("handle %v", err)
	}
	_, err = os.Stat(tpath)
	if err != nil {
		t.Fatalf("unable to locate torrent")
	}

	herror, err := NewScript("handler-error", matchertest.Error, script...)
	err = herror.Handle(tpath, meta)
	if err == nil || err == NoMatch {
		t.Fatalf("handle %v", err)
	}
	_, err = os.Stat(tpath)
	if err != nil {
		t.Fatalf("unable to locate torrent")
	}

}
