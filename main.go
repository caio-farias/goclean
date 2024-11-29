package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
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

var (
	path          string
	age           string
	filename      string
	cron_exp      string
	schedule_mode bool
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "goclean",                                // Name of your command
		Short: "Schedule tasks for cleaning irectories", // Short description
		Long:  "TODO",                                   // Long description
		Run: func(cmd *cobra.Command, args []string) {

			if schedule_mode {
				schdl := NewScheduler(SCHEDULE_FILENAME)
				var wg sync.WaitGroup
				wg.Add(1)
				defer wg.Wait()
				go schdl.Init(&wg)
				return
			}

			if _, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Println("File/dir does not exist")
				return
			}

			ageVerified, err := strconv.Atoi(age)
			if err != nil {
				log.Fatalln("age is NAN")
			}

			Clean(path, ageVerified, filename)
		},
	}

	rootCmd.Flags().BoolVarP(&schedule_mode, "schedule_mode", "s", false, "path to directory")
	rootCmd.Flags().StringVarP(&path, "path", "p", "./aaa", "path to directory")
	rootCmd.Flags().StringVarP(&age, "age", "a", "0", "Target age")
	rootCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename")
	rootCmd.Flags().StringVarP(&cron_exp, "cron_exp", "c", "", "cron_exp")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func logDeletions(targets []string, space_restored_in_bytes int) {
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

func init() {
	file, err := os.OpenFile(LOG_FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	LOG_FILE.File = file
}
