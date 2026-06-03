package logger

import (
	"log"
	"os"
)

var Log *log.Logger

func Init() {
	Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}
