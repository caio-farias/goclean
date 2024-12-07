package cmd

import (
	"fmt"
	"goclean/pkg"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/spf13/cobra"
)

const (
	LOG_FILENAME      = ".goclean.log"
	SCHEDULE_FILENAME = ".goclean.cron.json"
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
				schdl := pkg.NewScheduler(SCHEDULE_FILENAME)
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

			pkg.Clean(path, ageVerified, filename)
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

func init() {
	file, err := os.OpenFile(LOG_FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	pkg.LOG_FILE.File = file
}
