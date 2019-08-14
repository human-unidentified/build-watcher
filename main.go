package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"time"
)

// Получение названия последней директории по указанному пути
func getLastDirName(path string) (lastDir string) {
	allFiles, _ := ioutil.ReadDir(path)

	sort.SliceStable(allFiles, func(i, j int) bool {
		diff := allFiles[i].ModTime().Sub(allFiles[j].ModTime())
		return diff > 0
	})

	for _, f := range allFiles {
		if f.IsDir() {
			return f.Name()
		}
	}

	return ""
}

func watchDir(path string) {
	fmt.Println("Start watch folder ", path)
	watcherInterval := time.Second * 5
	for true {
		lastDir := getLastDirName(path)
		if lastDir != "" {
			fmt.Println("lastDir=", lastDir)
		}

		fmt.Printf("Sleep %s\n", watcherInterval)
		time.Sleep(watcherInterval)
	}

}

func main() {
	watchDir("\\\\s6\\BuildArchive\\T-FLEX DOCs 17\\DOCsDev\\")
}
