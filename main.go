package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/TomOnTime/utfutil"
)

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
	lastDir, err := getLastBuildDirName(path)

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

	appendBuildDirToProcessed(fullBuildPath)
}

func processFinishedBuild(buildDir string) error {
	fmt.Printf("Process build dir %s.\n", buildDir)

	containMessage, message, err := getBuildMessage(buildDir)
	if err != nil {
		return err
	}

	fmt.Printf("containMessage = %t\n", containMessage)
	fmt.Printf("message=%s\n", message)

	if containMessage {
		err = sendEmail(message)
		if err != nil {
			return err
		}
	}

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
		logLine := strings.Trim(scanner.Text(), " ")

		if strings.Contains(strings.ToUpper(logLine), "ПРЕДУПРЕЖДЕНИЙ") {
			if logLine == "Предупреждений: 0" {
				return false, logLine, nil
			} else {
				return true, logLine, nil
			}
		}

	}

	return false, "", scanner.Err()
}

func sendEmail(message string) error {
	sendParameters := getSmtpSettings(credentialsFileName)

	m := gomail.NewMessage()
	m.SetHeader("From", sendParameters.from)
	m.SetHeader("To", sendParameters.to)
	m.SetHeader("Subject", "Build monitor")
	m.SetBody("text/html", message)

	d := gomail.NewDialer(sendParameters.smtpServer, sendParameters.port, sendParameters.user, sendParameters.password)
	d.SSL = true

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil

}

var buildDir = "\\\\s6\\BuildArchive\\T-FLEX DOCs 17\\DOCsDev"
var processedBuildFile = "data/processed-list.txt"
var credentialsFileName = "data/credentials"

func main() {
	startupDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	credentialsFileName = startupDir + string(os.PathSeparator) + credentialsFileName
	processedBuildFile = startupDir + string(os.PathSeparator) + processedBuildFile

	os.OpenFile(processedBuildFile, os.O_RDONLY|os.O_CREATE, 0666)

	watchBuildDir(buildDir)
}
