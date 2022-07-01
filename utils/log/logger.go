package log

import (
	"io"
	"log"
	"os"

	"jarvis-trading-bot/consts"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func Init() {
	var mw io.Writer
	logFile := os.Getenv(consts.LogFile)
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv(consts.LogToFile) == "true" {
		mw = io.MultiWriter(os.Stdout, file)
	} else {
		mw = io.Writer(os.Stdout)
	}

	InfoLogger = log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(mw, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
