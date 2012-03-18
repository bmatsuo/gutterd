
[install go]: http://golang.org/install.html "Install Go"
[the godoc url]: http://localhost:6060/pkg/github.com/bmatsuo/gutterd/ "the Godoc URL"
[the gopkgdoc url]: http://gopkgdoc.appspot.com/pkg/github.com/bmatsuo/gutterd "the GoPkgDoc URL"

About gutterd
=============

Gutterd is a deamon process. It demuxes .torrent files found in directories like
your web browser's download directory and Dropbox directories. After analyzing
their metadata, .torrent files are moved to different BitTorrent clients'
watched directories.

Documentation
=============

Usage
-----

Run gutterd with the command

    gutterd [OPTIONS]

For help with command line options.

    gutterd -h

Configuration
-------------

Gutterd uses a JSON configuration stored in `~/.config/gutterd.json`
The configuration specifies directories to watch for incoming torrents,
as well as the handlers to match against those torrents. Here is an
example configuration.

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

Handlers
--------

When handler 'match' properties are unspecified, they will match any torrent.
Torrents are matched against handlers in order. So, in the example above, the
'other' handler acts as a catch-all and will match all torrents not matched by
any other handler.

Prerequisites
-------------

[Install Go][].

Use the weekly branch.

Installation
-------------

Use goinstall to install gutterd

    go get github.com/bmatsuo/gutterd

General Documentation
---------------------

Use godoc to vew the documentation for gutterd

    go doc github.com/bmatsuo/gutterd

Or alternatively, visit [the GoPkgDoc URL][]

Author
======

Bryan Matsuo &lt;bryan.matsuo@gmail.com&gt;

Copyright & License
===================

Copyright (c) 2012, Bryan Matsuo.
All rights reserved.
Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.
