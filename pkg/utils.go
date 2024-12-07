package pkg

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	LOG_FILENAME        = ".goclean.log"
	SCHEDULE_FILENAME   = ".goclean.cron.json"
	DEFAULT_DATE_FORMAT = "Jan 02 2006 - 15:04:05"
)

type LogFile struct {
	Name string
	File *os.File
}

var (
	LOG_FILE = LogFile{
		Name: LOG_FILENAME,
	}
)

func Filter[T any](arr []T, filterFn func(T) bool) []T {
	var filtered []T
	for _, item := range arr {
		if filterFn(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func FindOne[T any](arr []T, filterFn func(T) bool) T {
	subArr := Filter(arr, filterFn)
	return subArr[0]
}

func WriteLog(log_content string) {
	nowFormatted := time.Now().Format(DEFAULT_DATE_FORMAT)
	var fileContent []byte
	_, err := LOG_FILE.File.Read(fileContent)
	if err != nil {
		log.Fatalln("Could not open file", err)
	}
	line := fmt.Sprintf("[%s] %s\n", nowFormatted, log_content)
	LOG_FILE.File.Write([]byte(append(fileContent, []byte(line)...)))
}

func LogDeletions(targets []string, space_restored_in_bytes int) {
	var log_content string

	kb_restored := space_restored_in_bytes / 1024
	if len(targets) > 0 {
		log_content = fmt.Sprintf("Content removed (%d KB): \n", kb_restored)
		for _, t := range targets {
			log_content += fmt.Sprintf(" - %s \n", t)
		}
	} else {
		log_content = "No content was removed. \n"
	}

	WriteLog(log_content)
}
