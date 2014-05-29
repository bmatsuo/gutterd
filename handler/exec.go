package handler

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bmatsuo/torrent/metainfo"
)

type Exec interface {
	Run(path string, meta *metainfo.Metainfo) error
}

type ExecFunc func(path string, meta *metainfo.Metainfo) error

func (fn ExecFunc) Run(path string, meta *metainfo.Metainfo) error {
	return fn(path, meta)
}

type errExec struct{ err error }

func (ee errExec) Run(path string, meta *metainfo.Metainfo) error {
	return ee.err
}

type WatchDir string

func (dir WatchDir) Run(path string, _ *metainfo.Metainfo) error {
	return os.Rename(path, filepath.Join(string(dir), filepath.Base(path)))
}

type ScriptTemplate struct {
	t *template.Template
}

func NewScriptTemplate(script ...string) (*ScriptTemplate, error) {
	t, err := template.New("").Parse(strings.Join(script, "\n"))
	if err != nil {
		return nil, err
	}
	st := &ScriptTemplate{t}
	return st, nil
}

func (st *ScriptTemplate) Run(path string, meta *metainfo.Metainfo) error {
	contxt := map[string]interface{}{
		"Path": path,
	}
	f, err := ioutil.TempFile("", "gutterd-script-")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	err = st.t.Execute(f, contxt)
	if err != nil {
		f.Close()
		return err
	}
	err = f.Close()
	if err != nil {
		return fmt.Errorf("%q couldn't create script file: %v", err)
	}
	cmd := exec.Command("/bin/bash", f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%q script failed: %v", err)
	}
	return nil
}
