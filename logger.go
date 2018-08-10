package sewan_go_sdk

import (
	"io"
	"log"
	"os"
	"strings"
)

// This wrapper is used only for plugin developpment debug purpose,
// it must be removed at the end of the develpment cycle, before delivery to prod.
// It creates a logger that write logs to files in sdk-logs/ folder, stored in
// current folder.
func LoggerCreate(logFile string) *log.Logger {
	return loggerCreate(logFile)
}

func loggerCreate(logFile string) *log.Logger {
	logFolder := "sdk-logs/"
	var logFilePath strings.Builder
	logFilePath.WriteString(logFolder)
	logFilePath.WriteString(logFile)
	_, folder_existsError := os.Stat(logFolder)
	if folder_existsError != nil {
		os.Mkdir(logFolder, 0777)
	}
	var _, file_existsError = os.Stat(logFilePath.String())
	if file_existsError == nil {
		os.Remove(logFilePath.String())
	}
	os.Create(logFilePath.String())
	logFileObject, logFileErr := os.OpenFile(logFilePath.String(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if logFileErr != nil {
		log.Fatalln("Failed to open log file :", logFileErr)
	}
	logWriter := io.MultiWriter(logFileObject)
	return log.New(logWriter, "Sewan Provider : ", log.Ldate|log.Ltime)
}
