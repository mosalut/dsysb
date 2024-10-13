package main

import (
	"strconv"
	"os"
	"time"
	"log"
)

var logFile *os.File
const logName = "dsysb_log"

func openLogFile() error {
	var err error
	logFile, err = os.OpenFile(logName, os.O_WRONLY, 0644)
	if err != nil {
		logFile, err = os.Create(logName)
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

func print(level int, v ...any) error {
	if level == 0 {
		log.SetPrefix("[DEBUG]")
	} else {
		log.SetPrefix("[ERROR]")
	}

	info, err := logFile.Stat()
	if err != nil {
		return err
	}

	if info.Size() >= 4194304 { // 4 * 1024 * 1024
		err := setLogFile()
		if err != nil {
			return err
		}
	}

	log.Println(v)

	return nil
}
