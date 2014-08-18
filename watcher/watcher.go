package watcher

import (
	"gopkg.in/fsnotify.v0"
)

type Watcher struct {
	Event chan *fsnotify.FileEvent
	*fsnotify.Watcher
}

func New(filter Filter) (*Watcher, error) {
	return NewInstr(filter, nil)
}

// instrumentable
func NewInstr(filter Filter, errHandler func(error)) (*Watcher, error) {
	w := &Watcher{
		Event: make(chan *fsnotify.FileEvent, 1),
	}
	var err error
	w.Watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	go func() {
		for event := range w.Watcher.Event {
			// TODO filter could take a long time...
			if filter(event) {
				w.Event <- event
			}
		}
		close(w.Event)
	}()
	go func() {
		for err := range w.Watcher.Error {
			if errHandler != nil {
				go errHandler(err)
			}
		}
	}()
	return w, nil
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
