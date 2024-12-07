package pkg

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// Q ---  Q --- S
// A1 15:00 --- A2 14 --- 15:00

func wildcardMatch(text, pattern string) bool {
	textLen := len(text)
	patternLen := len(pattern)
	textIdx, patternIdx := 0, 0
	starIdx, matchIdx := -1, 0

	for textIdx < textLen {
		// If characters match or pattern has '*', continue
		if patternIdx < patternLen && (pattern[patternIdx] == text[textIdx]) {
			textIdx++
			patternIdx++
		} else if patternIdx < patternLen && pattern[patternIdx] == '*' {
			// '*' found, record the position
			starIdx = patternIdx
			matchIdx = textIdx
			patternIdx++
		} else if starIdx != -1 {
			// Backtrack to the last '*' position
			patternIdx = starIdx + 1
			matchIdx++
			textIdx = matchIdx
		} else {
			// No match found
			return false
		}
	}

	// Check if remaining pattern characters are all '*'
	for patternIdx < patternLen && pattern[patternIdx] == '*' {
		patternIdx++
	}

	// If we reached the end of the pattern, it's a match
	return patternIdx == patternLen
}

func FindTargets(path string, final_date time.Time, filename string) ([]string, error) {
	targets := []string{}

	err := filepath.Walk(path, func(file_path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if anyMatch(
			true,
			path == file_path,                 // exclude root
			!final_date.After(info.ModTime()), // exclude filter files with early age
			// exclude non filename match
			filename != "" && !wildcardMatch(info.Name(), filename)) {
			return nil
		}

		targets = append(targets, file_path)
		return nil
	})

	return targets, err
}

func anyMatch(match bool, inputs ...bool) bool {
	for _, input := range inputs {
		if input == match {
			return match
		}
	}

	return !match
}

func DeleteTargets(targets []string) int {
	deletionClearedSpace := 0
	for _, t := range targets {
		bytes, errRead := os.ReadFile(t)
		err := os.Remove(t)
		if err != nil || errRead != nil {
			log.Fatalln("error trying to delete target", t)
		}
		deletionClearedSpace += len(bytes)
	}
	return deletionClearedSpace
}

func Clean(path string, age int, filename string) {
	now := time.Now()
	final_date := now.AddDate(0, 0, -age)

	targets, _ := FindTargets(path, final_date, filename)

	if len(targets) > 0 {
		space_restored_in_bytes := DeleteTargets(targets)
		log.Println(space_restored_in_bytes)
		LogDeletions(targets, space_restored_in_bytes)
	}

	log.Printf("No targets like [path: %s] [filename %s] [age: %d] to clean.", path, filename, age)
}
