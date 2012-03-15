// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    log.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     2012-03-13 00:45:19.076882 -0700 PDT
 *  Description: 
 */

import (
	"fmt"
	"log"
	"os"
)

var loggerMux = new(LoggerMux)

var DefaultLogger = loggerMux.NewSource("gutterd")

func Output(calldepth int, s string) error      { return DefaultLogger.Output(calldepth+1, s) }
func Debug(v ...interface{})                    { DefaultLogger.Output(4, fmt.Sprintf("DEBUG\t", fmt.Sprint(v...))) }
func Debugf(format string, v ...interface{})    { DefaultLogger.Output(4, fmt.Sprintf(format, v...)) }
func Debugln(v ...interface{})                  { DefaultLogger.Output(4, fmt.Sprintln(v...)) }
func Info(v ...interface{})                     { DefaultLogger.Output(4, fmt.Sprintf("INFO\t", fmt.Sprint(v...))) }
func Infof(format string, v ...interface{})     { DefaultLogger.Output(4, fmt.Sprintf(format, v...)) }
func Infoln(v ...interface{})                   { DefaultLogger.Output(4, fmt.Sprintln(v...)) }
func Notice(v ...interface{})                   { DefaultLogger.Output(4, fmt.Sprintf("NOTICE\t", fmt.Sprint(v...))) }
func Noticef(format string, v ...interface{})   { DefaultLogger.Output(4, fmt.Sprintf(format, v...)) }
func Noticeln(v ...interface{})                 { DefaultLogger.Output(4, fmt.Sprintln(v...)) }
func Warn(v ...interface{})                     { DefaultLogger.Output(4, fmt.Sprintf("WARN\t", fmt.Sprint(v...))) }
func Warnf(format string, v ...interface{})     { DefaultLogger.Output(4, fmt.Sprintf(format, v...)) }
func Warnln(v ...interface{})                   { DefaultLogger.Output(4, fmt.Sprintln(v...)) }
func Error(v ...interface{})                    { DefaultLogger.Output(4, fmt.Sprintf("ERROR\t", fmt.Sprint(v...))) }
func Errorf(format string, v ...interface{})    { DefaultLogger.Output(4, fmt.Sprintf(format, v...)) }
func Errorln(v ...interface{})                  { DefaultLogger.Output(4, fmt.Sprintln(v...)) }
func Critical(v ...interface{})                 { DefaultLogger.Output(4, fmt.Sprintf("CRITICAL\t", fmt.Sprint(v...))) }
func Criticalf(format string, v ...interface{}) { DefaultLogger.Output(4, fmt.Sprintf(format, v...)) }
func Criticalln(v ...interface{})               { DefaultLogger.Output(4, fmt.Sprintln(v...)) }
func Print(v ...interface{})                    { DefaultLogger.Output(4, fmt.Sprint(v...)) }
func Printf(format string, v ...interface{})    { DefaultLogger.Output(4, fmt.Sprintf(format, v...)) }
func Println(v ...interface{})                  { DefaultLogger.Output(4, fmt.Sprintln(v...)) }
func Fatal(v ...interface{}) {
	DefaultLogger.Output(4, fmt.Sprint(v...))
	os.Exit(1)
}
func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Output(4, fmt.Sprintf(format, v...))
	os.Exit(1)
}
func Fatalln(v ...interface{}) {
	DefaultLogger.Output(4, fmt.Sprintln(v...))
	os.Exit(1)
}
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	DefaultLogger.Output(4, s)
	panic(s)
}
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	DefaultLogger.Output(4, s)
	panic(s)
}
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	DefaultLogger.Output(4, s)
	panic(s)
}

type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Output(calldepth int, s string) error
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// A gutterd internal logger object. Implements the Logger interface.
type gLogger struct {
	name   string
	output func(calldepth int, s string) error
}

func newGgLogger(name string, output func(calldepth int, s string) error) *gLogger {
	if name == "" {
		panic("no name")
	}
	if output == nil {
		panic("nil func")
	}
	return &gLogger{name, output}
}

func (l *gLogger) Output(calldepth int, s string) error {
	return l.output(calldepth, fmt.Sprintf("%s\t%s", l.name, s))
}
func (l *gLogger) Print(v ...interface{})                 { l.Output(4, fmt.Sprint(v...)) }
func (l *gLogger) Printf(format string, v ...interface{}) { l.Output(4, fmt.Sprintf(format, v...)) }
func (l *gLogger) Println(v ...interface{})               { l.Output(4, fmt.Sprintln(v...)) }
func (l *gLogger) Fatal(v ...interface{}) {
	l.Output(4, fmt.Sprint(v...))
	os.Exit(1)
}
func (l *gLogger) Fatalf(format string, v ...interface{}) {
	l.Output(4, fmt.Sprintf(format, v...))
	os.Exit(1)
}
func (l *gLogger) Fatalln(v ...interface{}) {
	l.Output(4, fmt.Sprintln(v...))
	os.Exit(1)
}
func (l *gLogger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Output(4, s)
	panic(s)
}
func (l *gLogger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(4, s)
	panic(s)
}
func (l *gLogger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Output(4, s)
	panic(s)
}

// An 'output' logger object that that can collect multiple log outputs.
type sinkLogger struct {
	*log.Logger
	accepts []string
}

// Implements the Logger interface
type LoggerMux struct {
	sources []*gLogger
	sinks   []*sinkLogger
}

func (mux *LoggerMux) NewSink(l *log.Logger, names ...string) {
	mux.sinks = append(mux.sinks, &sinkLogger{l, names})
}

func (mux *LoggerMux) NewSource(name string) Logger {
	for i := range mux.sources {
		if mux.sources[i].name == name {
			panic("duplicate source name")
		}
	}
	logger := &gLogger{name, func(calldepth int, s string) error {
		_mux := mux
		for _, lg := range _mux.sinks {
			for _, aname := range lg.accepts {
				if aname == name {
					if err := lg.Output(calldepth, s); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}}
	mux.sources = append(mux.sources, logger)
	return logger
}
