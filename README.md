
[install go]: http://golang.org/install.html "Install Go"
[the godoc url]: http://localhost:6060/pkg/github.com/bmatsuo/gutterd/ "the Godoc URL"

About gutterd
=============

Gutterd is a deamon process that inspects downloaded .torrent files and
organizes them into specific directories based on their content. 

This is meant to be used in conjuction with multiple sessions of `rtorrent`.

Documentation
=============

Gutterd works by having the user specify a set of torrent handlers.
Handlers handlers inspect the torrent and can be matched against them
based on several criteria. If a handler matches a given torrent, the
handler moves the torrent into a watch directory specific to the handler.

Configuration
-------------

Gutterd uses a json configuration stored in `~/.config/gutterd.json`
The configuration specifies directories to watch for incoming torrents,
as well as the handlers to match against those torrents. Here is an
example configuration.

    {
        "http": ":6060",
        "logPath": "&2",
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

When handler 'match' properties are unspecified, they will match any
torrent. Torrents are matched against handlers in the order they are
specified in the config file. So, in the example above, the 'other'
handler acts as a catch-all and will match all torrents not matched by
any other handler.

Usage
-----

Run gutterd with the command

    gutterd [OPTIONS]

For help with command line options.

    gutterd -h

Prerequisites
-------------

[Install Go][].

Installation
-------------

Use goinstall to install gutterd

    go get github.com/bmatsuo/gutterd

General Documentation
---------------------

Use godoc to vew the documentation for gutterd

    godoc github.com/bmatsuo/gutterd

Or alternatively, use a godoc http server

    godoc -http=:6060

and visit [the Godoc URL][]


Author
======

Bryan Matsuo &lt;bmatsuo@soe.ucsc.edu&gt;

Copyright & License
===================

Copyright (c) 2012, Bryan Matsuo.
All rights reserved.
Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.
