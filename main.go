package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

const (
	LOG_FILENAME      = ".goclean.log"
	SCHEDULE_FILENAME = ".goclean.cron.json"
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
	path     string
	age      string
	filename string
	pattern  string
	cron_exp string
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "goclean",                                // Name of your command
		Short: "Schedule tasks for cleaning irectories", // Short description
		Long:  "TODO",                                   // Long description
		Run: func(cmd *cobra.Command, args []string) {

			if _, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Println("File/dir does not exist")
				return
			}

			age, err := strconv.Atoi(age)
			if err != nil {
				log.Fatalln("age is NAN")
			}

			if cron_exp == "" {
				execClean(path, age, filename)
			}

			schdl := NewScheduler(SCHEDULE_FILENAME)
			schdl.add(cron_exp, func() {
				execClean(path, age, filename)
			})
			schdl.cron.Start()
			time.Sleep(2 * time.Minute)

			// Stop the scheduler gracefully
			schdl.cron.Stop()
			fmt.Println("Scheduler stopped")
		},
	}

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

func execClean(path string, age int, filename string) {
	now := time.Now()
	final_date := now.AddDate(0, 0, -age)

	targets, _ := FindTargets(path, final_date, filename)
	defer LOG_FILE.File.Close()
	logDeletions(targets)
	if len(targets) > 0 {
		DeleteTargets(targets)
	}
	return
}

func logDeletions(targets []string) {
	var log_content string

	now := time.Now()
	now_as_str := now.Format("Jan 02 2006 - 15:04:05")

	if len(targets) > 0 {
		log_content = fmt.Sprintf("[%s] Content removed: \n", now_as_str)
		for _, t := range targets {
			log_content += fmt.Sprintf(" - %s \n", t)
		}
	} else {
		log_content = fmt.Sprintf("[%s] No content was removed. \n", now_as_str)
	}

	WriteLog(log_content)
}

func WriteLog(log_content string) {
	var fileContent []byte
	_, err := LOG_FILE.File.Read(fileContent)
	if err != nil {
		log.Fatalln("Could not open file", err)
	}
	LOG_FILE.File.Write([]byte(append(fileContent, []byte(log_content)...)))
}

func init() {
	file, err := os.OpenFile(LOG_FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	LOG_FILE.File = file
}
