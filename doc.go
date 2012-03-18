// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*  Filename:    doc.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     2012-03-04 17:28:31.729616 -0800 PST
 *  Description: Godoc documentation for gutterd
 */


/*
Gutterd is a deamon process. It demuxes .torrent files found in directories like
your web browser's download directory and Dropbox directories. After analyzing
their metadata, .torrent files are moved to different BitTorrent clients'
watched directories.

Usage:

    gutterd [options]

Options:

Command line options can override system defaults and values specified in the
configuration file. Default values will be ignored.

	-config=""
  			A config file to use instead of ~/.config/gutterd.json.

	-http=""
  			Address to serve web requests from (e.g. ':6060').

	-log=""
			A path to log output. This overrides all logs specified in the config file.
			If -log-accepts is not specified, all logs will be output to the file.

	-log-accepts=""
			Log names for filtering the log specified with -log. If -log is not
			specified, accepted logs are printed to stderr. Similar to -log, logs in the
			config file wil be overridden when the flag is provided.

	-poll=0
			Specify a polling frequency (in seconds).

	-watch=""
			Specify a set of directories to watch.

Configuration:

Gutterd uses a JSON configuration. Most improtantly, the configuration
specifies directories to watch for incoming .torrent files, as well as handlers
to match against and demux those files. Here is an example configuration.

    {
        "http": ":6060",
		"logs": [
			{
				"path": "&2",
				"accepts": [ "gutterd", "http" ]
			}
		],
        "watch": [ "/Users/b/Downloads" ],
        "pollFrequency": 60,
        "handlers": [
            {
                "name": "music",
                "watch": "/Users/b/Music",
                "match": {
                    "tracker": "tracker\\.music\\.net",
                    "ext": "\\.(mp3|m4a|mp4)"
                }
            },
            {
                "name": "tv",
                "watch": "/Users/b/Movies",
                "match": { "tracker": "tracker\\.tv\\.net" }
            },
            {
                "name": "movies",
                "watch": "/Users/b/Movies",
                "match": { "tracker": "tracker\\.movies\\.net" }
            },
            {
                "name": "other",
                "watch": "/Users/b"
            }
        ]
    }

Handlers:

When handler 'match' properties are unspecified, they will match any torrent.
Torrents are matched against handlers in order. So, in the example above, the
'other' handler acts as a catch-all and will match all torrents not matched by
any other handler.
*/
package documentation
