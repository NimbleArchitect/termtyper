package main

import (
	"fmt"
	"log"
	"os"
)

func logError(msg ...interface{}) {
	if settings.termtyper.LogLevel >= 1 {
		log.Print("[ERROR] ", msg)
		writeFile("[ERROR] ", msg)
	}
}

func logWarn(msg ...interface{}) {
	if settings.termtyper.LogLevel >= 2 {
		log.Print("[WARN]", msg)
		writeFile("[WARN] ", msg)
	}
}

func logInfo(msg ...interface{}) {
	if settings.termtyper.LogLevel >= 3 {
		log.Print("[INFO] ", msg)
		writeFile("[INFO] ", msg)
	}
}

func logDebug(msg ...interface{}) {
	if settings.termtyper.LogLevel >= 4 {
		log.Print("[DEBUG] ", msg)
		writeFile("[DEBUG] ", msg)
	}
}

func writeFile(msg ...interface{}) {
	f, err := os.OpenFile("termtyper.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	str := fmt.Sprintf("%v\n", msg)
	if _, err := f.WriteString(str); err != nil {
		log.Println(err)
	}
}
