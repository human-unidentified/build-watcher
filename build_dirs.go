package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

// Получение названия последней директории по указанному пути
func getLastBuildDirName(path string) (string, error) {
	allFiles, err := ioutil.ReadDir(path)

	if err != nil {
		return "", err
	}

	sort.SliceStable(allFiles, func(i, j int) bool {
		diff := allFiles[i].ModTime().Sub(allFiles[j].ModTime())
		return diff > 0
	})

	for _, f := range allFiles {
		if f.IsDir() {
			return f.Name(), nil
		}
	}

	return "", fmt.Errorf("No directories found in %s", path)
}

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

func isBuildDirContainFinishedBuild(buildDir string) (bool, error) {
	// Достаточно знать что есть дистрибутив - значит лог есть
	rusDistribFoldername := buildDir + string(os.PathSeparator) + "Rus"

	_, err := os.Stat(rusDistribFoldername)

	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func appendBuildDirToProcessed(buildDir string) {
	f, err := os.OpenFile(processedBuildFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	if _, err = f.WriteString(buildDir + "\n"); err != nil {
		log.Fatal(err)
	}
}
