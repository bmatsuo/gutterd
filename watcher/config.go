package watcher

import (
	"fmt"
	"os"
)

type Config string

func (c Config) Validate() error {
	stat, err := os.Stat(string(c))
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("watch is not a directory: %s", c)
	}
	return nil
}
