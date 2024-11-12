package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const LOG_FILENAME = ".goclean.log"

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
	path string
	age  string
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

			chunks := strings.Split(age, "")
			if len(chunks) > 2 {
				log.Fatalln("invalid age")
			}

			age, err := strconv.Atoi(age)
			if err != nil {
				log.Fatalln("age is NAN")
			}

			now := time.Now()
			final := now.AddDate(0, 0, age)

			FindDir(path, final)
		},
	}

	rootCmd.Flags().StringVarP(&path, "path", "p", ".", "path to directory")
	rootCmd.Flags().StringVarP(&age, "age", "a", "1", "Target age")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	file, err := os.OpenFile(LOG_FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	LOG_FILE.File = file
}
