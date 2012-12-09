package watcher

import (
	"fmt"
	"github.com/howeyc/fsnotify"
)

type Watcher struct {
	Event chan *fsnotify.FileEvent
	*fsnotify.Watcher
}

func New(filter Filter) (w *Watcher, err error) {
	w = &Watcher{
		Event: make(chan *fsnotify.FileEvent, 1),
	}
	w.Watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	go func() {
		for event := range w.Watcher.Event {
			if filter(event) {
				w.Event <- event
			}
		}
		close(w.Event)
	}()
	go func() {
		for err := range w.Watcher.Error {
			// just suck these out for now
			fmt.Println(err)
		}
	}()
	return
}

func (w *Watcher) Watch(dirs ...Config) error {
	for _, d := range dirs {
		if err := w.Watcher.Watch(string(d)); err != nil {
			return err
		}
	}
	return nil
}

type Filter func(*fsnotify.FileEvent) bool
