package logger

import (
	"log"
	"os"
)

var (
	infoLog  = log.New(os.Stdout, "\033[35m[info]\033[0m", log.LstdFlags|log.Lshortfile)
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m", log.LstdFlags|log.Lshortfile)
)

var (
	Info   = infoLog.Println
	Infof  = infoLog.Panicf
	Error  = errorLog.Println
	Errorf = errorLog.Printf
)
