package main

import (
	"strconv"
	"os"
	"time"
	"log"
)

var logFile *os.File
const logName = "dsysb_log"

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

func print(level int, v ...any) error {
	switch level {
	case 0:
		log.SetPrefix("[DEBUG]")
	case 1:
		log.SetPrefix("[INFO]")
	case 2:
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
