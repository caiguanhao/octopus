package main

import (
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("octopus")

func initLogger(verbosity string, color bool) {
	var lvl logging.Level
	switch verbosity {
	case "debug":
		lvl = logging.DEBUG
	case "notice":
		lvl = logging.NOTICE
	case "warning":
		lvl = logging.WARNING
	case "error":
		lvl = logging.ERROR
	case "critical":
		lvl = logging.CRITICAL
	default:
		lvl = logging.INFO
	}

	var format logging.Formatter
	if color {
		format = logging.MustStringFormatter(`%{color}%{level: 6s} %{time:02/01/2006 15:04:05}%{color:reset} %{message}`)
	} else {
		format = logging.MustStringFormatter(`%{level: 6s} %{time:02/01/2006 15:04:05} %{message}`)
	}
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	formatter := logging.NewBackendFormatter(backend, format)
	level := logging.AddModuleLevel(formatter)
	level.SetLevel(lvl, "")
	logging.SetBackend(level)
}
