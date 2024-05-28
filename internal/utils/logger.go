package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func Logger(logLevel, field1, field2, msg string) (bool, string) {
	f, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return false, "error opening file"
	}
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(f)

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	log.WithFields(log.Fields{
		field1: field2,
	}).Info(msg)
	f.Close()

	return true, ""
}
