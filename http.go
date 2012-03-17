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
	"code.google.com/p/gorilla/mux"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type HTTPFormat uint

const (
	Hhtml HTTPFormat = iota
	Hjson
	Hinvalid
)

var httpLogger Logger

func _initHTTP() {
	if loggerMux == nil {
		panic("nil mux")
	}
	httpLogger = loggerMux.NewSource("http")
}

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
<html>
	<head>
		<title>config | gutterd</title>
		<meta http-equiv="content-type" content="text/html;charset=UTF-8" />
		</head>
	<body>
		<style type="text/css">
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
				background:#FFF;
			}
			</style>

		<h1>Configuration</h1>
		<div class="config">
			<div>
				<h2>Web server:</h2>
				{{.HTTP}}
				</div>
			<div>
				<h2>Logs:</h2>
				{{.LogPath}}
				</div>
			<div>
				<h2>Poll frequency:</h2>
				<table class="config">
					<tr>
						<td>{{.PollFrequency}}s</td>
						<td>
							<form name="upPoll" action="/config/pollFrequency" method="post">
								<input type="submit" value="⬆" />
								</form>
							<form name="downPoll" action="/config/pollFrequency" method="post">
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
							<form name="delWatch{{$i}}" action="/config/watch" method="post">
								<input type="submit" value="×" />
								</form>
							</td>
						</tr>
					{{end}}
					<tr>
						<form name="newWatch" action="/config/watch" method="post">
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
							<form name="delHandler{{$i}}" action="/config/handler/{{$h.Name}}/delete" method="post">
								<input type="submit" value="×" />
								</form>
							<form name="upHandler{{$i}}" action="/config/handler/{{$h.Name}}/up" method="post">
								<input type="submit" value="⬆" />
								</form>
							<form name="downHandler{{$i}}" action="/config/handler/{{$h.Name}}/down" method="post">
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
					</table>
				</div>
			</div>
		</body>
	</html>
`
var configHTMLTemplate = template.Must(template.New("config").Parse(configHTMLTemplateString))

func ConfigControllerShow(w http.ResponseWriter, r *http.Request) {
	format := HTTPRequestFormat(r)
	switch format {
	case Hhtml:
		err := configHTMLTemplate.Execute(w, config)
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
	http.Redirect(w, r, "/config", http.StatusFound)
}

func ConfigControllerWatchUpdate(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/config", http.StatusFound)
}

func ConfigControllerSave(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/config", http.StatusFound)
}

func HandlerControllerNew(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "New handler form.")
}

func HandlerControllerPrepend(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/config", http.StatusFound)
}

func HandlerControllerAppend(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/config", http.StatusFound)
}

func HandlerControllerInsert(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/config", http.StatusFound)
}

func HandlerControllerDelete(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/config", http.StatusFound)
}

func HandlerControllerUp(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/config", http.StatusFound)
}

func HandlerControllerDown(w http.ResponseWriter, r *http.Request) {
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
	router.HandleFunc("/config/handlers/prepend", HandlerControllerPrepend).
		Methods("POST")
	router.HandleFunc("/config/handlers/append", HandlerControllerAppend).
		Methods("POST")
	router.HandleFunc("/config/handlers/insert", HandlerControllerInsert).
		Methods("POST")
	router.HandleFunc("/config/handlers/new", HandlerControllerNew).
		Methods("GET")
	router.HandleFunc("/config/handler/{handler}/delete", HandlerControllerDelete).
		Methods("POST")
	router.HandleFunc("/config/handler/{handler}/up", HandlerControllerUp).
		Methods("POST")
	router.HandleFunc("/config/handler/{handler}/down", HandlerControllerDown).
		Methods("POST")
	router.HandleFunc("/config/handler/{handler}/name", HandlerControllerNameUpdate).
		Methods("POST")
	router.HandleFunc("/config/handler/{handler}/watch", HandlerControllerWatchUpdate).
		Methods("POST")
	router.HandleFunc("/config/pollFrequency", ConfigControllerPollUpdate).
		Methods("POST")
	router.HandleFunc("/config/watch", ConfigControllerWatchUpdate).
		Methods("POST")
	router.HandleFunc("/config/save", ConfigControllerSave).
		Methods("POST")
	router.HandleFunc(`/config{format:(\.(json|html))?}`, ConfigControllerShow).
		Methods("GET")
	http.ListenAndServe(config.HTTP, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpLogger.Infof("Request\t%v\t%v", r.Method, r.URL.Path)
		router.ServeHTTP(w, r)
	}))
}
