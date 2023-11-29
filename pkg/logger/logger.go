package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	Log *log.Logger
)

func init() {
	logFile, err := os.Create("info.log")
	fmt.Println("Logfile: info.log")
	if err != nil {
		panic(err)
	}
	Log = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	mw := io.MultiWriter(os.Stdout, logFile)
	Log.SetOutput(mw)
}
