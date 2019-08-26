package main

import (
	"log"
	"strconv"

	"github.com/TomOnTime/utfutil"
)

type smtpSettings struct {
	from       string
	to         string
	smtpServer string
	port       int
	user       string
	password   string
}

func getSmtpSettings(credentialsFileName string) smtpSettings {
	scanner, err := utfutil.NewScanner(credentialsFileName, utfutil.UTF8)
	if err != nil {
		log.Fatal(err)
	}
	defer scanner.Close()

	usedSettings := smtpSettings{}

	scanner.Scan()
	usedSettings.from = scanner.Text()
	scanner.Scan()
	usedSettings.to = scanner.Text()
	scanner.Scan()
	usedSettings.smtpServer = scanner.Text()
	scanner.Scan()
	usedSettings.port, _ = strconv.Atoi(scanner.Text())
	scanner.Scan()
	usedSettings.user = scanner.Text()
	scanner.Scan()
	usedSettings.password = scanner.Text()

	return usedSettings

}
