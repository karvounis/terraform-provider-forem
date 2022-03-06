package internal

import "log"

const (
	LOG_LEVEL_DEBUG = "DEBUG"
	LOG_LEVEL_INFO  = "INFO"
)

func LogDebug(msg string) {
	logLevel(LOG_LEVEL_DEBUG, msg)
}

func LogInfo(msg string) {
	logLevel(LOG_LEVEL_INFO, msg)
}

func logLevel(lvl, msg string) {
	log.Printf("[%s] %s", lvl, msg)
}
