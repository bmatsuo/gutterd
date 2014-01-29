// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    http.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     2012-03-11 16:24:28.171009 -0700 PDT
 *  Description:
 */

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"github.com/bmatsuo/gutterd/handler"
	"github.com/bmatsuo/gutterd/watcher"
)

type HTTPFormat uint

const (
	Hhtml HTTPFormat = iota
	Hjson
	Hinvalid
)

func HTTPRequestFormat(r *http.Request) HTTPFormat {
	if format := mux.Vars(r)["format"]; format != "" {
		switch strings.ToLower(format) {
		case ".html":
			return Hhtml
		case ".json":
			return Hjson
		}
		return Hinvalid
	}
	f := r.URL.Query().Get("format")
	if f == "" {
		return Hhtml
	}
	switch strings.ToLower(f) {
	case "html":
		return Hhtml
	case "json":
		return Hjson
	}
	return Hinvalid
}

var configHTMLTemplateString = `
{{define "css"}}
.config {
	margin-left:50;
}
.config ul {
	list-style-type:none;
}
.config table {
	/*border-width: 1px;*/
	/*border-color: gray;*/
	/*border-style: outset;*/
}
.config th {
	padding: 3px;
	padding-left: 5px;
	padding-right: 5px;
	/*border-width: 1px;*/
	/*border-style: solid;*/
	/*border-color: gray;*/
	/*background-color: #EEEEEE;*/
}
.config td {
	padding: 3px;
	padding-left: 5px;
	padding-right: 5px;
	/*border-width: 1px;*/
	/*border-style: solid;*/
	/*border-color: gray;*/
}
.config form {
	padding:0;
	margin:0;
}
.config input {
	padding:0;
	margin:0;
	width:100%;
	background:#FFF;
}
{{end}}

{{define "html"}}
<html>
	<head>
		<title>config | gutterd</title>
		<meta http-equiv="content-type" content="text/html;charset=UTF-8" />
		</head>
	<body>
		<style type="text/css">
			{{template "css"}}
			</style>

		<h1>Configuration</h1>
		<div class="config note">
			Changes will not overwrite the configuration file at {{.Path}}.
			</div>
		<div class="config">
			<div>
				<h2>Web server:</h2>
				<div class="config">{{.HTTP}}</div>
				</div>
			<div>
				<h2>Logs:</h2>
				<table class="config">
					<tr>
						<th>Name</th>
						<th>Accepts</th>
						</tr>
					{{range .Logs}}
					<tr>
						<td>{{.Path}}</td>
						<td>
							{{range .Accepts}}
							{{.}},
							{{end}}
							</td>
						</tr>
					{{end}}
					</table>
				</div>
			<div>
				<h2>Poll frequency:</h2>
				<table class="config">
					<tr>
						<td>{{.PollFrequency}}s</td>
						<td>
							<form name="upPoll" action="/config/pollFrequency" method="post">
								<input type="hidden" name="delta" value="5" />
								<input type="submit" value="⬆" />
								</form>
							<form name="downPoll" action="/config/pollFrequency" method="post">
								<input type="hidden" name="delta" value="-5" />
								<input type="submit" value="⬇" />
								</form>
							</td>
						</tr>
					</table>
				</div>
			<div>
				<h2>Watch directories:</h2>
				<table class="config">
					{{range $i, $w := .Watch}}
					<tr>
						<td>{{$w}}</td>
						<td>
							<form name="delWatch{{$i}}" action="/config/watch/{{$i}}/delete" method="post">
								<input type="submit" value="×" />
								</form>
							</td>
						</tr>
					{{end}}
					<tr>
						<form name="newWatch" action="/config/watch/add" method="post">
							<td><input type="text" name="watch" /></td>
							<td><input type="submit" value="+" /></td>
							</form>
						</tr>
					</table>
				</div>
			<div>
				<h2>Handlers:</h2>
				<table class="config">
					<tr>
						<th></th>
						<th>Name</th>
						<th>Watch</th>
						<th colspan=3><strong>Match</strong></th>
						</tr>
					<tr>
						<th></th>
						<th></th>
						<th></th>
						<th>Tracker</th>
						<th>Basename</th>
						<th>Extension</th>
						</tr>
					{{range $i, $h := .Handlers}}
					<tr>
						<th>
							<form name="delHandler{{$i}}" action="/config/handlers/{{$h.Name}}/delete" method="post">
								<input type="submit" value="×" />
								</form>
							<form name="upHandler{{$i}}" action="/config/handlers/{{$h.Name}}/up" method="post">
								<input type="submit" value="⬆" />
								</form>
							<form name="downHandler{{$i}}" action="/config/handlers/{{$h.Name}}/down" method="post">
								<input type="submit" value="⬇" />
								</form>
							</th>
						<td>{{$h.Name}}</td>
						<td>{{$h.Watch}}</td>
						<td>{{$h.Match.Tracker}}</td>
						<td>{{$h.Match.Basename}}</td>
						<td>{{$h.Match.Ext}}</td>
						</tr>
					{{end}}
					<tr>
						<form name="newHondler" id="newHandler" action="/config/handlers/create" method="post">
							<th>
								<input type="submit" value="✓" />
								<input type="button" value="×" />
								</th>
							<td><input type="text" name="name" /></td>
							<td><input type="text" name="watch" /></td>
							<td><input type="text" name="tracker" /></td>
							<td><input type="text" name="basename" /></td>
							<td><input type="text" name="ext" /></td>
							</form>
						</tr>
					</table>
				</div>
			</div>
		</body>
	</html>
{{end}}
`
var configHTMLTemplate = template.Must(template.New("config").Parse(configHTMLTemplateString))

func ConfigControllerShow(w http.ResponseWriter, r *http.Request) {
	format := HTTPRequestFormat(r)
	switch format {
	case Hhtml:
		err := configHTMLTemplate.ExecuteTemplate(w, "html", config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case Hjson:
		var p []byte
		var err error
		w.Header().Set("Content-type", "application/json")
		if p, err = json.Marshal(config); err != nil {
			http.Error(w, `{"error": "couldn't marshal configuration"}`,
				http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", p)
	default:
		http.Error(w, `{"error": "Invalid parameter: format"}`,
			http.StatusBadRequest)
	}
}

func ConfigControllerPollUpdate(w http.ResponseWriter, r *http.Request) {
	_delta := r.FormValue("delta")
	if _delta == "" {
		glog.Info("[http] no pollFrequency delta")
		http.Redirect(w, r, "/config", http.StatusFound)
		return
	}
	delta, err := strconv.ParseInt(_delta, 10, 64)
	if err != nil {
		glog.Infof("[http] error parsing delta: %v", err)
		http.Redirect(w, r, "/config", http.StatusFound)
		return
	}
	freq := config.PollFrequency + delta
	if freq < 0 {
		glog.Infof("[http] negative frequency: ", freq)
		http.Redirect(w, r, "/config", http.StatusFound)
		return
	}

	glog.Infof("updating poll frequency : %d", freq)
	config.PollFrequency = freq
	http.Redirect(w, r, "/config", http.StatusFound)
}

func ConfigControllerWatchAdd(w http.ResponseWriter, r *http.Request) {
	_path := strings.TrimFunc(r.FormValue("watch"), unicode.IsSpace)
	if _path == "" {
		glog.Warning("[http] no 'watch' specified.")
		http.Redirect(w, r, "/config", http.StatusFound)
		return
	}
	if !strings.HasPrefix(_path, "/") {
		glog.Warning("[http] non-absolute path specified.")
		http.Redirect(w, r, "/config", http.StatusFound)
		return
	}
	path, err := filepath.EvalSymlinks(_path)
	if err != nil {
		glog.Warning("[http] error evaluating symlinks: ", err)
		http.Redirect(w, r, "/config", http.StatusNotFound)
		return
	}
	stat, err := os.Stat(path)
	if err != nil {
		glog.Warningf("[http] cannot stat: %v", err)
		http.Redirect(w, r, "/config", http.StatusNotFound)
		return
	}
	if !stat.IsDir() {
		glog.Warningf("[http] not a directory: %v", _path)
		http.Redirect(w, r, "/config", http.StatusNotFound)
		return
	}
	for _, _watch := range config.Watch {
		watch, err := filepath.EvalSymlinks(string(_watch))
		if err != nil {
			glog.Warningf("[http] error evaluating symlinks: %v", err)
			http.Redirect(w, r, "/config", http.StatusNotFound)
			return
		}
		if watch == path {
			glog.Warningf("[http] already watching: %v", _path)
			http.Redirect(w, r, "/config", http.StatusNotFound)
			return
		}
	}

	glog.Warningf("watching directory: %v", _path)
	config.Watch = append(config.Watch, watcher.Config(_path))
	http.Redirect(w, r, "/config", http.StatusFound)
}

func ConfigControllerWatchDelete(w http.ResponseWriter, r *http.Request) {
	_i, err := strconv.ParseInt(mux.Vars(r)["index"], 10, 0)
	if err != nil {
		http.Redirect(w, r, "/config", http.StatusInternalServerError)
	}
	i := int(_i)
	if i >= len(config.Watch) {
		http.Redirect(w, r, "/config", http.StatusNotFound)
	}

	glog.Warningf("no longer watching directory: %v", config.Watch[i])
	config.Watch = append(config.Watch[:i], config.Watch[i+1:]...)
	http.Redirect(w, r, "/config", http.StatusFound)
}

func ConfigControllerSave(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/config", http.StatusFound)
}

func HandlerControllerNew(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "New handler form.")
}

func HandlerControllerCreate(w http.ResponseWriter, r *http.Request) {
	var hc handler.Config
	hc.Name = r.FormValue("name")
	hc.Watch = r.FormValue("watch")
	hc.Match.Tracker = r.FormValue("tracker")
	hc.Match.Basename = r.FormValue("basename")
	hc.Match.Ext = r.FormValue("ext")
	if err := hc.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, h := range config.Handlers {
		if h.Name == hc.Name {
			http.Error(w, fmt.Sprintf("handler exists with name %q", hc.Name), http.StatusBadRequest)
			return
		}
	}

	config.Handlers = append(config.Handlers, hc)
	handlers = append(handlers, hc.Handler())

	glog.Infof("created handler: %v", hc)

	http.Redirect(w, r, "/config", http.StatusFound)
}

func HandlerControllerDelete(w http.ResponseWriter, r *http.Request) {
	h := mux.Vars(r)["handler"]
	hIndex := -1
	for i := range config.Handlers {
		if h == config.Handlers[i].Name {
			hIndex = i
			break
		}
	}

	if hIndex < 0 {
		http.Redirect(w, r, "/config", http.StatusNotFound)
	}
	if hIndex >= len(config.Handlers)+1 {
		http.Redirect(w, r, "/config", http.StatusNotFound)
	}

	handlers = append(handlers[:hIndex], handlers[hIndex+1:]...)
	config.Handlers = append(config.Handlers[:hIndex], config.Handlers[hIndex+1:]...)

	glog.Warningf("deleted handler: %v", h)

	http.Redirect(w, r, "/config", http.StatusFound)
}

func HandlerControllerUp(w http.ResponseWriter, r *http.Request) {
	h := mux.Vars(r)["handler"]
	hIndex := -1
	for i := range config.Handlers {
		if h == config.Handlers[i].Name {
			hIndex = i
			break
		}
	}

	if hIndex < 1 {
		http.Redirect(w, r, "/config", http.StatusNotFound)
	}

	handlers[hIndex-1], handlers[hIndex] = handlers[hIndex], handlers[hIndex-1]
	config.Handlers[hIndex-1], config.Handlers[hIndex] = config.Handlers[hIndex], config.Handlers[hIndex-1]

	handlerNames := make([]string, len(config.Handlers))
	for i := range config.Handlers {
		handlerNames[i] = config.Handlers[i].Name
	}
	glog.Warningf("new handler order: %v", handlerNames)

	http.Redirect(w, r, "/config", http.StatusFound)
}

func HandlerControllerDown(w http.ResponseWriter, r *http.Request) {
	h := mux.Vars(r)["handler"]
	hIndex := -1
	for i := range config.Handlers {
		if h == config.Handlers[i].Name {
			hIndex = i
			break
		}
	}

	if hIndex >= len(config.Handlers)-1 {
		http.Redirect(w, r, "/config", http.StatusNotFound)
	}

	handlers[hIndex], handlers[hIndex+1] = handlers[hIndex+1], handlers[hIndex]
	config.Handlers[hIndex], config.Handlers[hIndex+1] = config.Handlers[hIndex+1], config.Handlers[hIndex]

	handlerNames := make([]string, len(config.Handlers))
	for i := range config.Handlers {
		handlerNames[i] = config.Handlers[i].Name
	}
	glog.Infof("new handler order: %v", handlerNames)
	http.Redirect(w, r, "/config", http.StatusFound)
}

func HandlerControllerNameUpdate(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/config", http.StatusFound)
}

func HandlerControllerWatchUpdate(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/config", http.StatusFound)
}

func ListenAndServe() {
	router := mux.NewRouter()
	router.HandleFunc("/config/handlers/new", HandlerControllerNew).
		Methods("GET")
	router.HandleFunc("/config/handlers/create", HandlerControllerCreate).
		Methods("POST")
	router.HandleFunc("/config/handlers/{handler}/delete", HandlerControllerDelete).
		Methods("POST")
	router.HandleFunc("/config/handlers/{handler}/up", HandlerControllerUp).
		Methods("POST")
	router.HandleFunc("/config/handlers/{handler}/down", HandlerControllerDown).
		Methods("POST")
	router.HandleFunc("/config/handlers/{handler}/name", HandlerControllerNameUpdate).
		Methods("POST")
	router.HandleFunc("/config/handlers/{handler}/watch", HandlerControllerWatchUpdate).
		Methods("POST")
	router.HandleFunc("/config/pollFrequency", ConfigControllerPollUpdate).
		Methods("POST")
	router.HandleFunc("/config/watch/{index:[0-9]+}/delete", ConfigControllerWatchDelete).
		Methods("POST")
	router.HandleFunc("/config/watch/add", ConfigControllerWatchAdd).
		Methods("POST")
	router.HandleFunc("/config/save", ConfigControllerSave).
		Methods("POST")
	router.HandleFunc(`/config{format:(\.(json|html))?}`, ConfigControllerShow).
		Methods("GET")
	http.ListenAndServe(config.HTTP, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		router.ServeHTTP(w, r)
		dur := time.Since(start)
		glog.Infof("[http] request %v %v %v %dms",
			r.Method, r.URL, r.Proto, dur/time.Millisecond)
	}))
}
