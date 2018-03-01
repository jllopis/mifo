// Copyright 2017 Joan Llopis. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package log

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/jllopis/mifo/version"
)

var (
	defaultLogger *Log
)

// StdLogger defines an interface that will accept a stdlib
// compatible logger.
type StdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

// Log permite la configuración del logger
type Log struct {
	logger  StdLogger
	Service string
	Version string
}

// LogConfig contiene la configuración para el logger
type LogConfig struct {
	Logger  *log.Logger
	Service string
	Version string
}

func init() {
	serviceName := version.Name
	defaultLogger = &Log{
		logger: log.New(os.Stdout,
			"",
			log.Ldate|log.Ltime),
		Service: serviceName,
		Version: version.Version,
	}
}

// New creates a new Log instance. If an error
// occurs, *Log should be nil.
func New(c *LogConfig) *Log {
	return &Log{
		logger:  c.Logger,
		Service: c.Service,
		Version: c.Version,
	}
}

// SetDefaultLogger ajusta las opciones por defecto del logger
func SetLogger(logger StdLogger) {
	if logger != nil {
		defaultLogger.logger = logger
	}
}

func (l *Log) SetService(s string) *Log {
	if s != "" {
		l.Service = s
	}
	return l
}

func (l *Log) SetVersion(v string) *Log {
	if v != "" {
		l.Version = v
	}
	return l
}

// Err publicará un log de error utilizando el formato para GCE
// (https://cloud.google.com/error-reporting/docs/formatting-error-messages?hl=es)
func Err(str ...interface{}) {
	if defaultLogger.logger == nil {
		return
	}

	pc, fl, ln, _ := runtime.Caller(1)

	pre := fmt.Sprintf("svc=%s tp=error src=%s:%d fn=%s", defaultLogger.Service, fl, ln, path.Base(runtime.FuncForPC(pc).Name()))
	data := prepareEntry(pre, str)

	defaultLogger.logger.Printf(data)
}

// Info publicará información de contexto, no referida a un error.
func Info(str ...interface{}) {
	if defaultLogger.logger == nil {
		return
	}

	pc, fl, ln, _ := runtime.Caller(1)

	pre := fmt.Sprintf("svc=%s tp=info src=%s:%d fn=%s", defaultLogger.Service, fl, ln, path.Base(runtime.FuncForPC(pc).Name()))
	data := prepareEntry(pre, str)

	defaultLogger.logger.Print(data)
}

// prepareEntry will parse the input params and build a hashmap from them.
// It will return the string with the message to log
func prepareEntry(pre string, m []interface{}) string {
	var msg string
	if len(m)%2 != 0 {
		msg = fmt.Sprintf("%s msg=%s", pre, m[0].(string))
		m = m[1:]
	}
	if len(m)%2 != 0 {
		return msg
	}
	for i := 0; i < len(m); i = i + 2 {
		msg = fmt.Sprintf("%s %s=%v", msg, m[i].(string), m[i+1])
	}
	return msg
}
