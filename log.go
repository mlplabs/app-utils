package utils

import (
	"io"
	"log"
	"os"
)

type logObj struct{
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

var Log logObj 


func (l *logObj) Init(filePathStdLog, filePathErrLog string) {

	aWritersStd := make([]io.Writer, 0, 2)
	aWritersStd = append(aWritersStd, os.Stdout)

	aWritersErr := make([]io.Writer, 0, 2)
	aWritersErr = append(aWritersErr, os.Stderr)

	if filePathStdLog != ""{
		file, err := os.OpenFile(filePathStdLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file ", filePathStdLog, ":", err)
		}
		aWritersStd = append(aWritersStd, file)
	}
	if filePathErrLog != ""{
		file, err := os.OpenFile(filePathErrLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file ", filePathErrLog, ":", err)
		}
		aWritersStd = append(aWritersStd, file)
	}

	writerStd  := io.MultiWriter(aWritersStd...)
	writerErr := io.MultiWriter(aWritersErr...)


	l.Trace = log.New(writerStd, "TRACE: ", log.Ldate|log.Ltime)
	l.Info = log.New(writerStd, "INFO: ",  log.Ldate|log.Ltime)
	l.Warning = log.New(writerErr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.Error = log.New(writerErr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}