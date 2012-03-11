
[install go]: http://golang.org/install.html "Install Go"
[the godoc url]: http://localhost:6060/pkg/github.com/bmatsuo/gutterd/ "the Godoc URL"

About gutterd
=============

Gutterd is a deamon process that inspects downloaded .torrent files and
organizes them into specific directories based on their content. 

This is meant to be used in conjuction with multiple sessions of `rtorrent`.

Documentation
=============

Usage
-----

Run gutterd with the command

    gutterd [options]

Prerequisites
-------------

[Install Go][].

Installation
-------------

Use goinstall to install gutterd

    goinstall github.com/bmatsuo/gutterd

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
