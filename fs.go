package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

func removeFiles(directories []fs.DirEntry, final_date time.Time) {
	for _, dir := range directories {
		info, _ := dir.Info()
		if info.ModTime().After(final_date) {
			continue
		}
		if info.IsDir() {
			nextDirectories, err := os.ReadDir(info.Name())
			if err != nil {
				fmt.Println("Error removing file:", err)
			}
			removeFiles(nextDirectories, final_date)
		}

		if !info.IsDir() {
			fullPath, err := filepath.Abs(info.Name())
			if err != nil {
				fmt.Println("Error getting full path:", err)
				return
			}

			// Remove the file
			err = os.Remove(fullPath)
			if err != nil {
				fmt.Println("Error removing file:", err)
			}
		}

	}
}

func FindDir(path string, final_date time.Time) {
	directories, err := os.ReadDir(path)
	file, errFile := os.ReadFile(path)
	if err != nil && errFile != nil {
		log.Fatal("Error reading directory:", err)
		return
	}
	if errFile == nil {
		print(file)
	}
	removeFiles(directories, final_date)
}
