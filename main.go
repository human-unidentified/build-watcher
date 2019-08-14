package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"time"
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

	if lastDir != "" {
		fmt.Printf("lastDir=%s\n", lastDir)
	}
}

func main() {
	//buildDir := "\\\\s6\\BuildArchive\\T-FLEX DOCs 17\\DOCsDev\\DOCsDev 17.0.0.0 24.05.2019 15.06\\logs\\CompactRelease"
	buildDir := "\\\\s6\\BuildArchive\\T-FLEX DOCs 17\\DOCsDev\\"
	watchBuildDir(buildDir)
}
