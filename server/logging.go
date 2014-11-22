package server

import (
	log "github.com/funkygao/log4go"
	"os"
	"syscall"
)

func SetupLogging(logFile, logLevel, crashLogFile string) {
	level := log.DEBUG
	switch logLevel {
	case "info":
		level = log.INFO

	case "warn":
		level = log.WARNING

	case "error":
		level = log.ERROR

	case "debug":
		level = log.DEBUG

	case "trace":
		level = log.TRACE
	}

	for _, filter := range log.Global {
		filter.Level = level
	}

	if logFile == "stdout" {
		log.AddFilter("stdout", level, log.NewConsoleLogWriter())
	} else {
		writer := log.NewFileLogWriter(logFile, false)
		log.DeleteFilter("stdout")
		log.AddFilter("file", level, writer)
		writer.SetFormat("[%d %T] [%L] (%S) %M")
		writer.SetRotate(true)
		writer.SetRotateSize(0)
		writer.SetRotateLines(0)
		writer.SetRotateDaily(true)
	}

	if crashLogFile != "" {
		f, err := os.OpenFile(crashLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}

		syscall.Dup2(int(f.Fd()), 2)
	}

}
