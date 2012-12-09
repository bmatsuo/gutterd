// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

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

var DefaultLoggerMux *LoggerMux
var DefaultLogger Logger

func Output(calldepth int, s string) error      { return DefaultLogger.Output(calldepth+1, s) }
func Debug(v ...interface{})                    { DefaultLogger.Output(4, fmt.Sprint("DEBUG\t", fmt.Sprint(v...))) }
func Debugf(format string, v ...interface{})    { DefaultLogger.Output(4, fmt.Sprintf(format, v...)) }
func Debugln(v ...interface{})                  { DefaultLogger.Output(4, fmt.Sprintln(v...)) }
func Info(v ...interface{})                     { DefaultLogger.Output(4, fmt.Sprint("INFO\t", fmt.Sprint(v...))) }
func Infof(format string, v ...interface{})     { DefaultLogger.Output(4, fmt.Sprintf(format, v...)) }
func Infoln(v ...interface{})                   { DefaultLogger.Output(4, fmt.Sprintln(v...)) }
func Notice(v ...interface{})                   { DefaultLogger.Output(4, fmt.Sprint("NOTICE\t", fmt.Sprint(v...))) }
func Noticef(format string, v ...interface{})   { DefaultLogger.Output(4, fmt.Sprintf(format, v...)) }
func Noticeln(v ...interface{})                 { DefaultLogger.Output(4, fmt.Sprintln(v...)) }
func Warn(v ...interface{})                     { DefaultLogger.Output(4, fmt.Sprint("WARN\t", fmt.Sprint(v...))) }
func Warnf(format string, v ...interface{})     { DefaultLogger.Output(4, fmt.Sprintf(format, v...)) }
func Warnln(v ...interface{})                   { DefaultLogger.Output(4, fmt.Sprintln(v...)) }
func Error(v ...interface{})                    { DefaultLogger.Output(4, fmt.Sprint("ERROR\t", fmt.Sprint(v...))) }
func Errorf(format string, v ...interface{})    { DefaultLogger.Output(4, fmt.Sprintf(format, v...)) }
func Errorln(v ...interface{})                  { DefaultLogger.Output(4, fmt.Sprintln(v...)) }
func Critical(v ...interface{})                 { DefaultLogger.Output(4, fmt.Sprint("CRITICAL\t", fmt.Sprint(v...))) }
func Criticalf(format string, v ...interface{}) { DefaultLogger.Output(4, fmt.Sprintf(format, v...)) }
func Criticalln(v ...interface{})               { DefaultLogger.Output(4, fmt.Sprintln(v...)) }
func Print(v ...interface{})                    { DefaultLogger.Output(4, fmt.Sprint(v...)) }
func Printf(format string, v ...interface{})    { DefaultLogger.Output(4, fmt.Sprintf(format, v...)) }
func Println(v ...interface{})                  { DefaultLogger.Output(4, fmt.Sprintln(v...)) }
func Fatal(v ...interface{}) {
	DefaultLogger.Output(4, fmt.Sprint("FATAL\t", fmt.Sprint(v...)))
	os.Exit(1)
}
func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Output(4, fmt.Sprint("FATAL\t", fmt.Sprintf(format, v...)))
	os.Exit(1)
}
func Fatalln(v ...interface{}) {
	DefaultLogger.Output(4, fmt.Sprint("FATAL\t", fmt.Sprintln(v...)))
	os.Exit(1)
}
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	DefaultLogger.Output(4, fmt.Sprint("PANIC\t", s))
	panic(s)
}
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	DefaultLogger.Output(4, fmt.Sprint("PANIC\t", s))
	panic(s)
}
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	DefaultLogger.Output(4, fmt.Sprint("PANIC\t", s))
	panic(s)
}

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})
	Notice(v ...interface{})
	Noticef(format string, v ...interface{})
	Noticeln(v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnln(v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})
	Critical(v ...interface{})
	Criticalf(format string, v ...interface{})
	Criticalln(v ...interface{})
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
func (l *gLogger) Debug(v ...interface{})                    { l.Output(4, fmt.Sprint("DEBUG\t", fmt.Sprint(v...))) }
func (l *gLogger) Debugf(format string, v ...interface{})    { l.Output(4, fmt.Sprint("DEBUG\t",fmt.Sprintf(format, v...))) }
func (l *gLogger) Debugln(v ...interface{})                  { l.Output(4, fmt.Sprint("DEBUG\t",fmt.Sprintln(v...))) }
func (l *gLogger) Info(v ...interface{})                     { l.Output(4, fmt.Sprint("INFO\t", fmt.Sprint(v...))) }
func (l *gLogger) Infof(format string, v ...interface{})     { l.Output(4, fmt.Sprint("INFO\t",fmt.Sprintf(format, v...))) }
func (l *gLogger) Infoln(v ...interface{})                   { l.Output(4, fmt.Sprint("INFO\t",fmt.Sprintln(v...))) }
func (l *gLogger) Notice(v ...interface{})                   { l.Output(4, fmt.Sprint("NOTICE\t", fmt.Sprint(v...))) }
func (l *gLogger) Noticef(format string, v ...interface{})   { l.Output(4, fmt.Sprint("NOTICE\t",fmt.Sprintf(format, v...))) }
func (l *gLogger) Noticeln(v ...interface{})                 { l.Output(4, fmt.Sprint("NOTICE\t",fmt.Sprintln(v...))) }
func (l *gLogger) Warn(v ...interface{})                     { l.Output(4, fmt.Sprint("WARN\t", fmt.Sprint(v...))) }
func (l *gLogger) Warnf(format string, v ...interface{})     { l.Output(4, fmt.Sprint("WARN\t",fmt.Sprintf(format, v...))) }
func (l *gLogger) Warnln(v ...interface{})                   { l.Output(4, fmt.Sprint("WARN\t",fmt.Sprintln(v...))) }
func (l *gLogger) Error(v ...interface{})                    { l.Output(4, fmt.Sprint("ERROR\t", fmt.Sprint(v...))) }
func (l *gLogger) Errorf(format string, v ...interface{})    { l.Output(4, fmt.Sprint("ERROR\t",fmt.Sprintf(format, v...))) }
func (l *gLogger) Errorln(v ...interface{})                  { l.Output(4, fmt.Sprint("ERROR\t",fmt.Sprintln(v...))) }
func (l *gLogger) Critical(v ...interface{})                 { l.Output(4, fmt.Sprint("CRITICAL\t", fmt.Sprint(v...))) }
func (l *gLogger) Criticalf(format string, v ...interface{}) { l.Output(4, fmt.Sprint("CRITICAL\t",fmt.Sprintf(format, v...))) }
func (l *gLogger) Criticalln(v ...interface{})               { l.Output(4, fmt.Sprint("CRITICAL\t",fmt.Sprintln(v...))) }
func (l *gLogger) Print(v ...interface{})                    { l.Output(4, fmt.Sprint(v...)) }
func (l *gLogger) Printf(format string, v ...interface{})    { l.Output(4, fmt.Sprintf(format, v...)) }
func (l *gLogger) Println(v ...interface{})                  { l.Output(4, fmt.Sprintln(v...)) }
func (l *gLogger) Fatal(v ...interface{}) {
	l.Output(4, fmt.Sprint("FATAL\t", fmt.Sprint(v...)))
	os.Exit(1)
}
func (l *gLogger) Fatalf(format string, v ...interface{}) {
	l.Output(4, fmt.Sprint("FATAL\t", fmt.Sprintf(format, v...)))
	os.Exit(1)
}
func (l *gLogger) Fatalln(v ...interface{}) {
	l.Output(4, fmt.Sprint("FATAL\t", fmt.Sprintln(v...)))
	os.Exit(1)
}
func (l *gLogger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Output(4, fmt.Sprint("PANIC\t", s))
	panic(s)
}
func (l *gLogger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(4, fmt.Sprint("PANIC\t", s))
	panic(s)
}
func (l *gLogger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Output(4, fmt.Sprint("PANIC\t", s))
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
	if len(names) == 0 {
		panic("sink accepts no sources")
	}
	mux.sinks = append(mux.sinks, &sinkLogger{l, names})
}

func (mux *LoggerMux) RemoveSink(l *log.Logger) {
	for i, s := range mux.sinks {
		if s.Logger == l {
			mux.sinks = append(append(make([]*sinkLogger, 0, len(mux.sinks)-1),
				mux.sinks[:i]...),
				mux.sinks[i+1:]...)
		}
	}
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

func (mux *LoggerMux) RemoveSource(name string) {
	for i, s := range mux.sources {
		if s.name == name {
			mux.sources = append(append(make([]*gLogger, 0, len(mux.sources)-1),
				mux.sources[:i]...),
				mux.sources[i+1:]...)
		}
	}
	// Don't remove the source from the sinks' "accepts" lists.
}
