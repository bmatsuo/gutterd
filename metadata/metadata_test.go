package metadata

/*  Filename:    metadata_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     2012-03-04 20:29:46.043866 -0800 PST
 *  Description: For testing metadata.go
 */

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadMetadataFile(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to determine working directory: %v", err)
	}
	t.Logf("pwd: %v", cwd)
	testfiles, err := filepath.Glob(filepath.Join(cwd, "test", "*.torrent"))
	if err != nil {
		t.Fatalf("failed to find test torrent files: %v", err)
	}
	for _, filename := range testfiles {
		meta, err := ReadMetadataFile(filename)
		if err != nil {
			t.Errorf("failed to read file: %v", err)
			continue
		}
		if meta.Announce == "" {
			t.Errorf("no announce url for file %q", filename)
		}
	}
}
