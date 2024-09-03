package utils

import (
	"bufio"
	"os"
	"strings"
)

// ReadLinesFromFile reads all lines from a file and returns them as a slice of strings.
func ReadLinesFromFile(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		return []string{}
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines
}

func CreateDirIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.Mkdir(dir, os.ModePerm)
	}
	return nil
}
