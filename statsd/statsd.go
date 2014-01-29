package statsd

import (
	"github.com/cactus/go-statsd-client/statsd"
)

var client *statsd.Client // Statsd

func Init(addr string, ns string) error {
	var err error
	client, err = statsd.New(addr, ns)
	if err != nil {
		client = nil
	}
	return err
}

func stat(fn func() error) error {
	if client != nil {
		return fn()
	}
	// ignore if no client is configured
	return nil
}

func Incr(name string, value int64, rate float32) error {
	return stat(func() error { return client.Inc(name, value, rate) })
}
