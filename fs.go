package main

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

func DeleteTargets(targets []string) {
	for _, t := range targets {
		err := os.Remove(t)
		if err != nil {
			log.Fatalln("error trying to delete target", t)
		}
	}
}
