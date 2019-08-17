package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/TomOnTime/utfutil"
)

// Получение названия последней директории по указанному пути
func getLastDirName(path string) (string, error) {
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

// Функция цикличного наблюдения за директорией
func watchBuildDir(path string) {
	fmt.Printf("Start watch folder %s\n", path)
	watcherInterval := time.Second * 5
	for true {
		watchCycle(path)

		fmt.Printf("Sleep %s\n", watcherInterval)
		time.Sleep(watcherInterval)
	}
}

// Цикл обработки
func watchCycle(path string) {
	lastDir, err := getLastDirName(path)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("lastDir=%s\n", lastDir)
	fullBuildPath := path + string(os.PathSeparator) + lastDir
	if isBuildDirProcessed(fullBuildPath) {
		fmt.Printf("Folder %s already processed.\n", fullBuildPath)
		return
	}

	buildFinished, err := isBuildDirContainFinishedBuild(fullBuildPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !buildFinished {
		fmt.Printf("Folder %s contains build in the process.\n", fullBuildPath)
		return
	}

	err = processFinishedBuild(fullBuildPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	appendBuildToProcessed(fullBuildPath)
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

func processFinishedBuild(buildDir string) error {
	fmt.Printf("Process build dir %s.\n", buildDir)

	containMessage, message, err := getBuildMessage(buildDir)
	if err != nil {
		return err
	}

	fmt.Printf("containMessage = %t\n", containMessage)
	fmt.Printf("message=%s.\n", message)

	return nil
}

func getBuildMessage(buildDir string) (bool, string, error) {
	buildLogFileName := buildDir + string(os.PathSeparator) + "logs\\Release\\DOCs.log"

	scanner, err := utfutil.NewScanner(buildLogFileName, utfutil.HTML5)
	if err != nil {
		return false, "", err
	}

	defer scanner.Close()

	scannerBuffer := make([]byte, 0, 64*1024)
	scanner.Buffer(scannerBuffer, 1024*1024)

	for scanner.Scan() {
		logLine := scanner.Text()

		if strings.Contains(strings.ToUpper(logLine), "ПРЕДУПРЕЖДЕНИЙ") {
			return true, logLine, nil
		}

	}

	return false, "", scanner.Err()
}

func appendBuildToProcessed(buildDir string) {
	f, err := os.OpenFile(processedBuildFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	if _, err = f.WriteString(buildDir); err != nil {
		log.Fatal(err)
	}
}

var buildDir = "\\\\s6\\BuildArchive\\T-FLEX DOCs 17\\DOCsDev"
var processedBuildFile = "processed.txt"

func main() {
	watchBuildDir(buildDir)
}
