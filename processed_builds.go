package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func isBuildDirProcessed(buildDir string) bool {
	file, err := os.OpenFile(processedBuildFile, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		scanLine := scanner.Text()

		if strings.EqualFold(scanLine, buildDir) || strings.EqualFold(buildDir, scanLine) {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return false
}

func appendBuildToProcessed(buildDir string) {
	f, err := os.OpenFile(processedBuildFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	if _, err = f.WriteString(buildDir + "\n"); err != nil {
		log.Fatal(err)
	}
}
