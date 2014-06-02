package main

/*  Filename:    config_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     2012-03-04 19:23:27.544848 -0800 PST
 *  Description: For testing config.go
 */
/*
import (
	"reflect"
	"testing"

	"github.com/bmatsuo/gutterd/handler"
	"github.com/bmatsuo/gutterd/matcher"
	"github.com/bmatsuo/gutterd/watcher"
)

var configValidTests = []struct {
	cstr string
	path string
	cdef *Config
	c    *Config
}{
	{ // Test the marshalling for each configuration name.
		`{
			"http": ":8080",
			"pollFrequency": 20,
			"watch": [ "/home/foo/Downloads" ],
			"logs": [ { "path": "&1", "accepts": [ "gutterd" ] } ],
			"handlers": [
			{
				"name": "foo",
				"watch": "./",
				"matcher": {
					"match": "all",
					"of": [
					{"tracker": "tracker\\.baz\\.net"},
					{"basename": "qux"},
					{"ext": ".quux"}
					]
				}
			}
			]
		}`,
		"./gutterd.json",
		&Config{},
		&Config{
			Path:          "./gutterd.json",
			PollFrequency: 20,
			Watch:         []watcher.Config{"/home/foo/Downloads"},
			Handlers: []handler.Config{{
				Name:  "foo",
				Watch: "./",
				Match: matcher.Config{
					Tracker:  `tracker\.baz\.net`,
					Basename: "qux",
					Ext:      ".quux"}}},
		}},
	{ // Test the overloading of default values.
		`{
			"logs": [ { "path": "&1", "accepts": [ "gutterd" ] } ],
			"handlers": [ {
				"name": "foo",
				"watch": "./",
				"match": {
					"tracker": "tracker\\.baz\\.net",
					"basename": "qux",
					"ext": ".quux" } } ]
		}`,
		"./gutterd.json",
		&Config{
			PollFrequency: 30,
			Watch:         []watcher.Config{"./"},
		},
		&Config{
			Path:          "./gutterd.json",
			PollFrequency: 30,
			Watch:         []watcher.Config{"./"},
			Handlers: []handler.Config{{
				Name:  "foo",
				Watch: "./",
				Match: matcher.Config{
					Tracker:  `tracker\.baz\.net`,
					Basename: "qux",
					Ext:      ".quux"}}},
		}},
}

func TestConfig(t *testing.T) {
	for i, test := range configValidTests {
		c, err := loadConfigFromBytes([]byte(test.cstr), test.path, test.cdef)
		if err != nil {
			t.Errorf("Test %d: Load error: %v", i, err)
			continue
		}
		if !reflect.DeepEqual(c, test.c) {
			t.Errorf("Test %d: Unequal configs:\nLoaded: %#v\nExpected: %#v", i, c, test.c)
			continue
		}
	}
}
*/
