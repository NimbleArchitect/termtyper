package main

import (
	"log"
)

func logError(msg ...interface{}) {
	if loglevel >= 1 {
		log.Print("[ERROR] ", msg)
	}
}

func logWarn(msg ...interface{}) {
	if loglevel >= 2 {
		log.Print("[WARN]", msg)
	}
}

func logInfo(msg ...interface{}) {
	if loglevel >= 3 {
		log.Print("[INFO] ", msg)
	}
}

func logDebug(msg ...interface{}) {
	if loglevel >= 4 {
		log.Print("[DEBUG] ", msg)
	}
}
