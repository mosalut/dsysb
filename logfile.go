// dsysb

package main

import (
	"strings"
	"runtime"
	"strconv"
	"os"
	"time"
	"log"
)

var logFile *os.File
const logName = "dsysb_log"

const (
	log_debug = iota
	log_info
	log_warning
	log_error
)

func openLogFile(host string) error {
	var err error
	filename := logName + host
	logFile, err = os.OpenFile(filename, os.O_WRONLY | os.O_APPEND, 0644)
	if err != nil {
		logFile, err = os.Create(filename)
		if err != nil {
			return err
		}
	}

	log.SetOutput(logFile)

	return nil
}

func setLogFile() error {
	logFile.Close()
	now := time.Now().Unix()

	err := os.Rename(logName, logName + "_" + strconv.FormatInt(now, 16))
	if err != nil {
		return err
	}

	logFile, err = os.Create(logName)
	if err != nil {
		return err
	}

	log.SetOutput(logFile)

	return nil
}

func print(level int, v ...any) {
	switch level {
	case log_debug:
		log.SetPrefix("[DEBUG]")
	case log_info:
		log.SetPrefix("[INFO]")
	case log_warning:
		log.SetPrefix("[WARNING]")
	case log_error:
		log.SetPrefix("[ERROR]")
	}

	if cmdFlag.logFile {
		info, err := logFile.Stat()
		if err != nil {
			log.Println(err)
		}

		if info.Size() >= 4194304 { // 4 * 1024 * 1024
			err := setLogFile()
			if err != nil {
				log.Println(err)
			}
		}
	}

	_, path, line, _ := runtime.Caller(1)
	names := strings.Split(path, `/`)
	filename := names[len(names) - 1]
	log.Println(filename, "line:", line, v)
}
